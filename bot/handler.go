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
	WelcomeMessage            = "Выберите действие:"
	BalanceMessage            = "Ваш текущий баланс 10 бабл-ти."
	PaymentMessage            = "Выберите способ оплаты:"
	AccountMessage            = "Информация о вашем  персональном аккаунте"
	PaymentOption1            = "Нажмите 'Произвести оплату' для продолжения"
	PaymentOption2            = "Нажмите 'Произвести оплату' для продолжения"
	TeamMessage               = "Сборище тех кто в пустые не забивает"
	ProcessingPayment         = "Обработка оплаты..."
	AlreadyProcessing         = "Оплата уже обрабатывается. Пожалуйста, подождите."
	LoginMessage              = "Пожалуйста, войдите в систему"
	LoginSuccessMessage       = "Вы успешно вошли в систему!"
	LoginFailMessage          = "Неверный пароль. Попробуйте снова."
	AuthStartMessage          = "Добро пожаловать! Нажмите кнопку для входа:"
	AskUsernameMessage        = "Введите ваше имя пользователя:"
	AskPasswordMessage        = "Введите ваш пароль:"
	CancelAuthMessage         = "Авторизация отменена"
	MASTER_PASSWORD           = "ks" // Задайте нужный пароль здесь
	LogoutMessage             = "Вы успешно вышли из системы!"
	TeamRosterMessage         = "Состав команды:"
	AskPositionMessage        = "Введите вашу позицию в команде:"
	AskBirthdayMessage        = "Введите вашу дату рождения (например, 01.01.1990):"
	ProfileSetupComplete      = "Информация сохранена!"
	AskNumberMessage          = "Введите ваш игровой номер:"
	PositionForward           = "Нападающий"
	PositionDefender          = "Защитник"
	PositionGoalie            = "Вратарь"
	PasswordSetMessage        = "Придумайте пароль:"
	PasswordSetSuccessMessage = "Пароль успешно установлен!"
	UserNotFoundMessage       = "Пользователь не найден. Пожалуйста, зарегистрируйтесь."
	SetPasswordButtonText     = "Установить пароль"
	PasswordNotSetMessage     = "Сначала установите пароль с помощью кнопки 'Установить пароль'"
	OnlyPasswordSetMessage    = "Сначала необходимо установить пароль"
	AskNewPasswordMessage     = "Введите пароль, который будет использоваться для входа:"
	SetPasswordSuccessMessage = "Пароль сохранен! Теперь вы можете войти, используя любой логин и установленный пароль."
	EnterAmountMessage        = "Введите сумму для пополнения (в рублях):"
	PaymentDetailsMessage     = "Для пополнения баланса переведите указанную сумму на счет:\n" +
		"Номер карты: 1234 5678 9012 3456\n" +
		"После перевода нажмите 'Подтвердить оплату'"
	PaymentConfirmationMessage = "Оплата будет проверена администратором в течение 24 часов"
	InvalidAmountMessage       = "Пожалуйста, введите корректную сумму (целое число больше 0)"
)

type AuthState int

const (
	StateNone AuthState = iota
	StateAwaitingUsername
	StateAwaitingPassword
	StateAwaitingPosition
	StateAwaitingBirthday
	StateAwaitingNumber
	StateAuthenticated
	StateAwaitingNewPassword
	StateAwaitingPaymentAmount
	StateAwaitingPaymentConfirmation
)

type UserAuthInfo struct {
	State         AuthState
	Username      string
	Position      string
	Birthday      string
	Number        string
	IsRegistering bool
	PaymentAmount float64
}

type Handler struct {
	client             *telegram.Client
	keyboard           telegram.ReplyKeyboardMarkup
	paymentKeyboard    telegram.ReplyKeyboardMarkup
	processKeyboard    telegram.ReplyKeyboardMarkup
	processingUsers    map[int64]bool
	mu                 sync.RWMutex
	authenticatedUsers map[int64]string // maps chatID to username
	authMu             sync.RWMutex
	authStates         map[int64]UserAuthInfo
	db                 *database.Database
	teamKeyboard       telegram.ReplyKeyboardMarkup
}

func NewHandler(client *telegram.Client, dbPath string) (*Handler, error) {
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		client:             client,
		keyboard:           CreateMainKeyboard(),
		paymentKeyboard:    CreatePaymentKeyboard(),
		processKeyboard:    CreatePaymentProcessKeyboard(),
		processingUsers:    make(map[int64]bool),
		authenticatedUsers: make(map[int64]string),
		authStates:         make(map[int64]UserAuthInfo),
		db:                 db,
		teamKeyboard:       CreateTeamKeyboard(),
	}
	return h, nil
}

func CreateAuthKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Войти"}},
			{{Text: SetPasswordButtonText}},
			{{Text: "Отмена"}},
		},
		ResizeKeyboard: true,
	}
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
		delete(h.processingUsers, userID) // Fixed: delete takes map and key
	}
}

func (h *Handler) isAuthenticated(chatID int64) bool {
	h.authMu.RLock()
	defer h.authMu.RUnlock()
	_, ok := h.authenticatedUsers[chatID]
	return ok
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

	if session.Position == "" || session.Birthday == "" {
		h.authStates[chatID] = UserAuthInfo{
			State:    StateAwaitingPosition,
			Username: session.Username,
		}
		return AskPositionMessage, nil
	}

	return fmt.Sprintf("Ваш аккаунт:\nИмя: %s\nПозиция: %s\nДата рождения: %s",
		session.Username, session.Position, session.Birthday), nil
}

func (h *Handler) createInitialKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: SetPasswordButtonText}},
			{{Text: "Войти"}},
			{{Text: "Отмена"}},
		},
		ResizeKeyboard: true,
	}
}

func (h *Handler) HandleUpdate(update telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Проверяем состояние авторизации
	authInfo, exists := h.authStates[chatID]
	if !exists {
		authInfo = UserAuthInfo{State: StateNone}
		h.authStates[chatID] = authInfo
	}

	// Обработка команды отмены в любом состоянии
	if text == "Отмена" {
		h.authStates[chatID] = UserAuthInfo{State: StateNone}
		return h.client.SendMessageWithKeyboard(chatID, AuthStartMessage, CreateAuthKeyboard())
	}

	// Обработка процесса авторизации
	if !h.isAuthenticated(chatID) {
		switch authInfo.State {
		case StateNone:
			switch text {
			case "/start":
				return h.client.SendMessageWithKeyboard(chatID, AuthStartMessage, h.createInitialKeyboard())
			case SetPasswordButtonText:
				authInfo.State = StateAwaitingNewPassword
				h.authStates[chatID] = authInfo
				return h.client.SendMessage(chatID, AskNewPasswordMessage)
			case "Войти":
				authInfo.State = StateAwaitingUsername
				h.authStates[chatID] = authInfo
				return h.client.SendMessage(chatID, AskUsernameMessage)
			default:
				return h.client.SendMessageWithKeyboard(chatID, AuthStartMessage, h.createInitialKeyboard())
			}

		case StateAwaitingNewPassword:
			// Сохраняем общий пароль для всех пользователей
			if err := h.db.SetMasterPassword(text); err != nil {
				return err
			}
			h.authStates[chatID] = UserAuthInfo{State: StateNone}
			return h.client.SendMessageWithKeyboard(chatID, SetPasswordSuccessMessage, h.createInitialKeyboard())

		case StateAwaitingUsername:
			authInfo.Username = text
			authInfo.State = StateAwaitingPassword
			h.authStates[chatID] = authInfo
			return h.client.SendMessage(chatID, AskPasswordMessage)

		case StateAwaitingPassword:
			// Проверяем пароль напрямую с мастер-паролем
			masterPass, err := h.db.GetMasterPassword()
			if err != nil {
				return err
			}
			if text != masterPass {
				h.authStates[chatID] = UserAuthInfo{State: StateNone}
				return h.client.SendMessageWithKeyboard(chatID, LoginFailMessage, h.createInitialKeyboard())
			}
			return h.handleLoginComplete(chatID, authInfo.Username, text)
		}
		return nil
	}

	// Обработка ввода данных профиля для авторизованных пользователей
	switch authInfo.State {
	case StateAwaitingPosition:
		if text != PositionForward && text != PositionDefender && text != PositionGoalie {
			return h.client.SendMessageWithKeyboard(chatID, "Пожалуйста, выберите позицию из предложенных", CreatePositionKeyboard())
		}
		authInfo.Position = text
		authInfo.State = StateAwaitingBirthday
		h.authStates[chatID] = authInfo
		return h.client.SendMessage(chatID, AskBirthdayMessage)

	case StateAwaitingBirthday:
		authInfo.Birthday = text
		authInfo.State = StateAwaitingNumber
		h.authStates[chatID] = authInfo
		return h.client.SendMessage(chatID, AskNumberMessage)

	case StateAwaitingNumber:
		authInfo.Number = text
		if err := h.db.SaveUserSession(chatID, authInfo.Username, authInfo.Position, authInfo.Birthday, authInfo.Number, 0.0); err != nil {
			return err
		}
		authInfo.State = StateAuthenticated
		h.authStates[chatID] = authInfo
		return h.client.SendMessageWithKeyboard(chatID, ProfileSetupComplete, h.keyboard)
	case StateAwaitingNewPassword:
		if err := h.db.SetUserPassword(authInfo.Username, text); err != nil {
			return err
		}
		authInfo.State = StateAuthenticated
		h.authStates[chatID] = authInfo
		return h.client.SendMessage(chatID, PasswordSetSuccessMessage)
	case StateAwaitingPaymentAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil || amount <= 0 {
			return h.client.SendMessage(chatID, InvalidAmountMessage)
		}
		authInfo.PaymentAmount = amount
		authInfo.State = StateAwaitingPaymentConfirmation
		h.authStates[chatID] = authInfo
		return h.client.SendMessageWithKeyboard(
			chatID,
			fmt.Sprintf(PaymentDetailsMessage+"\nСумма: %.2f руб.", amount),
			h.processKeyboard,
		)
	}

	// Основная логика обработки команд для авторизованных пользователей
	switch text {
	case "/start":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, WelcomeMessage, h.keyboard)
	case "Мой баланс":
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}
		return h.client.SendMessage(update.Message.Chat.ID,
			fmt.Sprintf("Ваш текущий баланс: %.2f ", session.Balance))
	case "Способ оплаты":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentMessage, h.paymentKeyboard)
	case "Мой аккаунт":
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}

		if session.Position == "" || session.Birthday == "" || session.Number == "" {
			authInfo.State = StateAwaitingPosition
			authInfo.Username = session.Username
			h.authStates[chatID] = authInfo
			return h.client.SendMessageWithKeyboard(chatID, AskPositionMessage, CreatePositionKeyboard())
		}

		return h.client.SendMessage(chatID, fmt.Sprintf(
			"Ваш аккаунт:\nИмя: %s\nПозиция: %s\nДата рождения: %s\nИгровой номер: %s",
			session.Username, session.Position, session.Birthday, session.Number))
	case "Информация о команде":
		// Теперь просто показываем клавиатуру без сообщения о команде
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, "Выберите раздел:", h.teamKeyboard)
	case "Состав":
		message, err := h.getTeamRosterMessage()
		if err != nil {
			return err
		}
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, message, h.teamKeyboard)
	case "1":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentOption1, h.processKeyboard)
	case "2":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentOption2, h.processKeyboard)
	case "Произвести оплату":
		if h.isProcessingPayment(update.Message.Chat.ID) {
			return h.client.SendMessage(update.Message.Chat.ID, AlreadyProcessing)
		}
		h.setProcessingPayment(update.Message.Chat.ID, true)
		return h.client.SendMessage(update.Message.Chat.ID, ProcessingPayment)
	case "Назад к способам оплаты":
		h.setProcessingPayment(update.Message.Chat.ID, false)
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentMessage, h.paymentKeyboard)
	case "Назад":
		h.setProcessingPayment(update.Message.Chat.ID, false)
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, WelcomeMessage, h.keyboard)
	case "Выйти":
		return h.logout(update.Message.Chat.ID)
	case "/setpassword":
		if !h.isAuthenticated(chatID) {
			return h.client.SendMessage(chatID, "Сначала войдите в систему")
		}
		authInfo.State = StateAwaitingNewPassword
		h.authStates[chatID] = authInfo
		return h.client.SendMessage(chatID, PasswordSetMessage)
	case "Пополнить баланс":
		authInfo.State = StateAwaitingPaymentAmount
		h.authStates[chatID] = authInfo
		return h.client.SendMessage(chatID, EnterAmountMessage)

	case "Подтвердить оплату":
		if authInfo.State != StateAwaitingPaymentConfirmation {
			return h.client.SendMessage(chatID, "Сначала введите сумму для оплаты")
		}
		// Добавляем сумму к текущему балансу пользователя
		session, err := h.db.GetUserSession(chatID)
		if err != nil {
			return err
		}
		newBalance := session.Balance + authInfo.PaymentAmount
		if err := h.db.UpdateBalance(chatID, newBalance); err != nil {
			return err
		}
		authInfo.State = StateAuthenticated
		h.authStates[chatID] = authInfo
		return h.client.SendMessageWithKeyboard(chatID,
			fmt.Sprintf("Баланс успешно пополнен на %.2f руб. Новый баланс: %.2f руб.",
				authInfo.PaymentAmount, newBalance),
			h.keyboard)

	case "Отменить":
		authInfo.State = StateAuthenticated
		h.authStates[chatID] = authInfo
		return h.client.SendMessageWithKeyboard(chatID, "Оплата отменена", h.keyboard)
	default:
		return nil
	}
}

func (h *Handler) handleLoginComplete(chatID int64, username, password string) error {
	h.authMu.Lock()
	h.authenticatedUsers[chatID] = username
	h.authMu.Unlock()

	// Сохраняем сессию
	if err := h.db.SaveUserSession(chatID, username, "", "", "", 0.0); err != nil {
		log.Printf("Error saving user session: %v", err)
	}

	h.authStates[chatID] = UserAuthInfo{
		State:    StateAuthenticated,
		Username: username,
	}

	return h.client.SendMessageWithKeyboard(chatID, LoginSuccessMessage, h.keyboard)
}

func (h *Handler) logout(chatID int64) error {
	h.authMu.Lock()
	delete(h.authenticatedUsers, chatID)
	h.authMu.Unlock()

	// Удаляем сессию из БД
	if err := h.db.RemoveUserSession(chatID); err != nil {
		log.Printf("Error removing user session: %v", err)
	}

	h.mu.Lock()
	delete(h.processingUsers, chatID)
	h.mu.Unlock()

	delete(h.authStates, chatID)
	return h.client.SendMessageWithKeyboard(chatID, LogoutMessage, CreateAuthKeyboard())
}
