package bot

import (
	"HelpBot/bot/database"
	"HelpBot/client/telegram"
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	WelcomeMessage        = "Выберите действие:"
	BalanceMessage        = "Ваш текущий баланс 10 бабл-ти."
	PaymentMessage        = "Выберите способ оплаты:"
	AccountMessage        = "Информация о вашем  персональном аккаунте"
	PaymentOption1        = "Нажмите 'Произвести оплату' для продолжения"
	PaymentOption2        = "Нажмите 'Произвести оплату' для продолжения"
	TeamMessage           = "Сборище тех кто в пустые не забивает"
	ProcessingPayment     = "Обработка оплаты..."
	AlreadyProcessing     = "Оплата уже обрабатывается. Пожалуйста, подождите."
	AskPositionMessage    = "Введите вашу позицию в команде:"
	AskBirthdayMessage    = "Введите вашу дату рождения (например, 01.01.1990):"
	ProfileSetupComplete  = "Информация сохранена!"
	AskNumberMessage      = "Введите ваш игровой номер:"
	PositionForward       = "Нападающий"
	PositionDefender      = "Защитник"
	PositionGoalie        = "Вратарь"
	TeamRosterMessage     = "Состав команды:"
	EnterAmountMessage    = "Введите сумму для пополнения (в рублях):"
	PaymentDetailsMessage = "Для пополнения баланса переведите указанную сумму на счет:\n" +
		"Номер карты: 1234 5678 9012 3456\n" +
		"После перевода нажмите 'Подтвердить оплату'"
	PaymentConfirmationMessage = "Оплата будет проверена администратором в течение 24 часов"
	InvalidAmountMessage       = "Пожалуйста, введите корректную сумму (целое число больше 0)"
)

type UserState int

const (
	StateNone UserState = iota
	StateAwaitingPosition
	StateAwaitingBirthday
	StateAwaitingNumber
	StateAwaitingPaymentAmount
	StateAwaitingPaymentConfirmation
)

type UserInfo struct {
	State         UserState
	Username      string
	Position      string
	Birthday      string
	Number        string
	PaymentAmount float64
}

type Handler struct {
	client          *telegram.Client
	keyboard        telegram.ReplyKeyboardMarkup
	paymentKeyboard telegram.ReplyKeyboardMarkup
	processKeyboard telegram.ReplyKeyboardMarkup
	processingUsers map[int64]bool
	mu              sync.RWMutex
	userStates      map[int64]UserInfo
	db              *database.Database
	teamKeyboard    telegram.ReplyKeyboardMarkup
}

func NewHandler(client *telegram.Client, dbPath string) (*Handler, error) {
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		client:          client,
		keyboard:        CreateMainKeyboard(),
		paymentKeyboard: CreatePaymentKeyboard(),
		processKeyboard: CreatePaymentProcessKeyboard(),
		processingUsers: make(map[int64]bool),
		userStates:      make(map[int64]UserInfo),
		db:              db,
		teamKeyboard:    CreateTeamKeyboard(),
	}
	return h, nil
}

func CreateTeamKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Состав"}},
			{{Text: "Назад"}},
		},
		ResizeKeyboard: true,
	}
}

func CreatePositionKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: PositionForward}},
			{{Text: PositionDefender}},
			{{Text: PositionGoalie}},
			{{Text: "Отмена"}},
		},
		ResizeKeyboard: true,
	}
}

func (h *Handler) isProcessingPayment(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.processingUsers[userID]
}

func (h *Handler) setProcessingPayment(userID int64, status bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if status {
		h.processingUsers[userID] = true
	} else {
		delete(h.processingUsers, userID)
	}
}

func (h *Handler) getTeamRosterMessage() (string, error) {
	users, err := h.db.GetAllActiveUsers()
	if err != nil {
		return "", err
	}

	message := "Состав команды:\n\n"
	for i, user := range users {
		message += fmt.Sprintf("%d. %s\n", i+1, user.Username)
	}

	if len(users) == 0 {
		message += "Пока никого нет в команде"
	}

	return message, nil
}

func (h *Handler) getAccountInfo(chatID int64) (string, error) {
	session, err := h.db.GetUserSession(chatID)
	if err != nil {
		return "", err
	}

	if session == nil {
		// Если сессии нет, создаем новую с именем пользователя по умолчанию
		username := fmt.Sprintf("user_%d", chatID)
		if err := h.db.SaveUserSession(chatID, username, "", "", "", 0.0); err != nil {
			return "", err
		}

		h.userStates[chatID] = UserInfo{
			State:    StateAwaitingPosition,
			Username: username,
		}
		return AskPositionMessage, nil
	}

	if session.Position == "" || session.Birthday == "" {
		h.userStates[chatID] = UserInfo{
			State:    StateAwaitingPosition,
			Username: session.Username,
		}
		return AskPositionMessage, nil
	}

	return fmt.Sprintf("Ваш аккаунт:\nИмя: %s\nПозиция: %s\nДата рождения: %s",
		session.Username, session.Position, session.Birthday), nil
}

func (h *Handler) HandleUpdate(update telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	userInfo, exists := h.userStates[chatID]
	if !exists {
		userInfo = UserInfo{State: StateNone}
		h.userStates[chatID] = userInfo

		// Автоматически создаем сессию для нового пользователя
		username := fmt.Sprintf("user_%d", chatID)
		if err := h.db.SaveUserSession(chatID, username, "", "", "", 0.0); err != nil {
			log.Printf("Error saving user session: %v", err)
		}
	}

	if text == "Отмена" {
		h.userStates[chatID] = UserInfo{State: StateNone}
		return h.client.SendMessageWithKeyboard(chatID, WelcomeMessage, h.keyboard)
	}

	// Обработка ввода данных профиля
	switch userInfo.State {
	case StateAwaitingPosition:
		if text != PositionForward && text != PositionDefender && text != PositionGoalie {
			return h.client.SendMessageWithKeyboard(chatID, "Пожалуйста, выберите позицию из предложенных", CreatePositionKeyboard())
		}
		userInfo.Position = text
		userInfo.State = StateAwaitingBirthday
		h.userStates[chatID] = userInfo
		return h.client.SendMessage(chatID, AskBirthdayMessage)

	case StateAwaitingBirthday:
		userInfo.Birthday = text
		userInfo.State = StateAwaitingNumber
		h.userStates[chatID] = userInfo
		return h.client.SendMessage(chatID, AskNumberMessage)

	case StateAwaitingNumber:
		userInfo.Number = text
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		username := session.Username
		if username == "" {
			username = fmt.Sprintf("user_%d", chatID)
		}

		if err := h.db.SaveUserSession(chatID, username, userInfo.Position, userInfo.Birthday, userInfo.Number, 0.0); err != nil {
			return err
		}
		userInfo.State = StateNone
		h.userStates[chatID] = userInfo
		return h.client.SendMessageWithKeyboard(chatID, ProfileSetupComplete, h.keyboard)

	case StateAwaitingPaymentAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil || amount <= 0 {
			return h.client.SendMessage(chatID, InvalidAmountMessage)
		}
		userInfo.PaymentAmount = amount
		userInfo.State = StateAwaitingPaymentConfirmation
		h.userStates[chatID] = userInfo
		return h.client.SendMessageWithKeyboard(
			chatID,
			fmt.Sprintf(PaymentDetailsMessage+"\nСумма: %.2f руб.", amount),
			h.processKeyboard,
		)
	}

	// Основная логика обработки команд
	switch text {
	case "/start":
		return h.client.SendMessageWithKeyboard(chatID, WelcomeMessage, h.keyboard)
	case "Мой баланс":
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		if session == nil {
			// Если сессии нет, создаем новую
			username := fmt.Sprintf("user_%d", chatID)
			if err := h.db.SaveUserSession(chatID, username, "", "", "", 0.0); err != nil {
				return err
			}
			return h.client.SendMessage(chatID, "Ваш текущий баланс: 0.00")
		}

		return h.client.SendMessage(chatID, fmt.Sprintf("Ваш текущий баланс: %.2f", session.Balance))
	case "Способ оплаты":
		return h.client.SendMessageWithKeyboard(chatID, PaymentMessage, h.paymentKeyboard)
	case "Мой аккаунт":
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		if session == nil {
			// Если сессии нет, создаем новую
			username := fmt.Sprintf("user_%d", chatID)
			if err := h.db.SaveUserSession(chatID, username, "", "", "", 0.0); err != nil {
				return err
			}

			userInfo.State = StateAwaitingPosition
			userInfo.Username = username
			h.userStates[chatID] = userInfo
			return h.client.SendMessageWithKeyboard(chatID, AskPositionMessage, CreatePositionKeyboard())
		}

		if session.Position == "" || session.Birthday == "" || session.Number == "" {
			userInfo.State = StateAwaitingPosition
			userInfo.Username = session.Username
			h.userStates[chatID] = userInfo
			return h.client.SendMessageWithKeyboard(chatID, AskPositionMessage, CreatePositionKeyboard())
		}

		return h.client.SendMessage(chatID, fmt.Sprintf(
			"Ваш аккаунт:\nИмя: %s\nПозиция: %s\nДата рождения: %s\nИгровой номер: %s",
			session.Username, session.Position, session.Birthday, session.Number))
	case "Информация о команде":
		return h.client.SendMessageWithKeyboard(chatID, "Выберите раздел:", h.teamKeyboard)
	case "Состав":
		message, err := h.getTeamRosterMessage()
		if err != nil {
			return err
		}
		return h.client.SendMessageWithKeyboard(chatID, message, h.teamKeyboard)
	case "1":
		return h.client.SendMessageWithKeyboard(chatID, PaymentOption1, h.processKeyboard)
	case "2":
		return h.client.SendMessageWithKeyboard(chatID, PaymentOption2, h.processKeyboard)
	case "Произвести оплату":
		if h.isProcessingPayment(chatID) {
			return h.client.SendMessage(chatID, AlreadyProcessing)
		}
		h.setProcessingPayment(chatID, true)
		return h.client.SendMessage(chatID, ProcessingPayment)
	case "Назад к способам оплаты":
		h.setProcessingPayment(chatID, false)
		return h.client.SendMessageWithKeyboard(chatID, PaymentMessage, h.paymentKeyboard)
	case "Назад":
		h.setProcessingPayment(chatID, false)
		return h.client.SendMessageWithKeyboard(chatID, WelcomeMessage, h.keyboard)
	case "Пополнить баланс":
		userInfo.State = StateAwaitingPaymentAmount
		h.userStates[chatID] = userInfo
		return h.client.SendMessage(chatID, EnterAmountMessage)
	case "Подтвердить оплату":
		if userInfo.State != StateAwaitingPaymentConfirmation {
			return h.client.SendMessage(chatID, "Сначала введите сумму для оплаты")
		}
		// Добавляем сумму к текущему балансу пользователя
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		if session == nil {

			username := fmt.Sprintf("user_%d", chatID)
			if err := h.db.SaveUserSession(chatID, username, "", "", "", userInfo.PaymentAmount); err != nil {
				return err
			}
		} else {
			newBalance := session.Balance + userInfo.PaymentAmount
			if err := h.db.UpdateBalance(chatID, newBalance); err != nil {
				return err
			}
		}

		userInfo.State = StateNone
		h.userStates[chatID] = userInfo

		// Получаем обновленный баланс
		updatedSession, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		return h.client.SendMessageWithKeyboard(chatID,
			fmt.Sprintf("Баланс успешно пополнен на %.2f руб. Новый баланс: %.2f руб.",
				userInfo.PaymentAmount, updatedSession.Balance),
			h.keyboard)
	case "Отменить":
		userInfo.State = StateNone
		h.userStates[chatID] = userInfo
		return h.client.SendMessageWithKeyboard(chatID, "Оплата отменена", h.keyboard)
	default:
		return nil
	}
}
