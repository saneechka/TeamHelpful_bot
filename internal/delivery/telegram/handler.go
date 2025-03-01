package telegram

import (
	"fmt"
	"log"
	"sync"

	"HelpBot/client/telegram"
	"HelpBot/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает сообщения от Telegram
type Handler struct {
	client          *telegram.Client
	userService     domain.UserService
	paymentService  domain.PaymentService
	sessionService  domain.SessionService
	keyboard        telegram.ReplyKeyboardMarkup
	paymentKeyboard telegram.ReplyKeyboardMarkup
	processKeyboard telegram.ReplyKeyboardMarkup
	teamKeyboard    telegram.ReplyKeyboardMarkup
	processingUsers map[int64]bool
	mu              sync.RWMutex
	authHandler     *AuthHandler
}

// NewHandler создает новый экземпляр Handler
func NewHandler(
	client *telegram.Client,
	userService domain.UserService,
	paymentService domain.PaymentService,
	sessionService domain.SessionService,
) *Handler {
	authHandler := NewAuthHandler(client, sessionService, userService)

	return &Handler{
		client:          client,
		userService:     userService,
		paymentService:  paymentService,
		sessionService:  sessionService,
		keyboard:        CreateMainKeyboard(),
		paymentKeyboard: CreatePaymentKeyboard(),
		processKeyboard: CreatePaymentProcessKeyboard(),
		teamKeyboard:    CreateTeamKeyboard(),
		processingUsers: make(map[int64]bool),
		authHandler:     authHandler,
	}
}

// isProcessingPayment проверяет, обрабатывается ли платеж для пользователя
func (h *Handler) isProcessingPayment(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.processingUsers[userID]
}

// setProcessingPayment устанавливает статус обработки платежа для пользователя
func (h *Handler) setProcessingPayment(userID int64, status bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if status {
		h.processingUsers[userID] = true
	} else {
		delete(h.processingUsers, userID)
	}
}

// getTeamRosterMessage формирует сообщение со списком команды
func (h *Handler) getTeamRosterMessage() (string, error) {
	users, err := h.userService.GetAllUsers()
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

// HandleUpdate обрабатывает обновление от Telegram
func (h *Handler) HandleUpdate(update *tgbotapi.Update) {
	// Обрабатываем только сообщения
	if update.Message == nil {
		log.Println("Received update without message, skipping")
		return
	}

	log.Printf("Received message from %s (%d): %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)

	// Получаем сессию пользователя
	session, err := h.sessionService.GetSession(update.Message.Chat.ID)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return
	}

	if session == nil {
		log.Printf("Session is nil for user %d, creating new session", update.Message.Chat.ID)
	} else {
		log.Printf("Session state for user %d: %v, authorized: %v", update.Message.Chat.ID, session.State, session.IsAuthorized)
	}

	// Обрабатываем команды
	if update.Message.IsCommand() {
		log.Printf("Processing command: %s", update.Message.Command())
		h.handleCommand(update.Message, session)
		return
	}

	// Обрабатываем текстовые сообщения
	log.Printf("Processing text message: %s", update.Message.Text)
	h.handleMessage(update.Message, session)
}

// handleCommand обрабатывает команды
func (h *Handler) handleCommand(message *tgbotapi.Message, session *domain.UserSession) {
	var err error

	switch message.Command() {
	case "start":
		err = h.authHandler.HandleStart(message)
	case "help":
		err = h.client.SendMessage(message.Chat.ID, "Доступные команды:\n/start - начать работу с ботом\n/help - показать справку")
	default:
		err = h.client.SendMessage(message.Chat.ID, "Неизвестная команда. Используйте /help для получения списка доступных команд.")
	}

	if err != nil {
		log.Printf("Error handling command: %v", err)
	}
}

// handleMessage обрабатывает текстовые сообщения
func (h *Handler) handleMessage(message *tgbotapi.Message, session *domain.UserSession) {
	var err error

	// Если сессия не существует, создаем ее
	if session == nil {
		err = h.authHandler.HandleStart(message)
		if err != nil {
			log.Printf("Error handling start: %v", err)
		}
		return
	}

	// Обрабатываем сообщения в зависимости от состояния сессии
	if session.State == domain.StateAwaitingUsername || session.State == domain.StateAwaitingPassword {
		// Пользователь в процессе авторизации или регистрации
		if message.Text == "Войти" {
			// Сбрасываем состояние и начинаем процесс входа
			session.State = domain.StateNone
			session.LastCommand = ""
			if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
				log.Printf("Error updating session: %v", err)
			}
			err = h.authHandler.HandleLogin(message)
		} else if message.Text == "Зарегистрироваться" {
			// Сбрасываем состояние и начинаем процесс регистрации
			session.State = domain.StateNone
			session.LastCommand = ""
			if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
				log.Printf("Error updating session: %v", err)
			}
			err = h.authHandler.HandleRegister(message)
		} else {
			// Продолжаем текущий процесс
			log.Printf("Processing input for state %v, last command: %s", session.State, session.LastCommand)

			// Определяем, какой процесс был начат по LastCommand
			if session.LastCommand == "register" {
				err = h.authHandler.HandleRegister(message)
			} else if session.LastCommand == "login" {
				err = h.authHandler.HandleLogin(message)
			} else {
				// Если LastCommand не установлен, пробуем определить по состоянию
				log.Printf("LastCommand not set, trying to determine by state")
				if session.State == domain.StateAwaitingUsername {
					// Если мы ожидаем имя пользователя, предполагаем, что это вход
					err = h.authHandler.HandleLogin(message)
				} else {
					// Если мы ожидаем пароль, проверяем, существует ли пользователь
					users, err := h.userService.GetAllUsers()
					if err == nil {
						userExists := false
						for _, u := range users {
							if u.Username == session.User.Username {
								userExists = true
								break
							}
						}

						if userExists {
							// Если пользователь существует, это вход
							err = h.authHandler.HandleLogin(message)
						} else {
							// Если пользователь не существует, это регистрация
							err = h.authHandler.HandleRegister(message)
						}
					} else {
						// В случае ошибки, пробуем продолжить процесс входа
						err = h.authHandler.HandleLogin(message)
					}
				}
			}
		}
	} else {
		// Обрабатываем сообщения авторизованного пользователя
		switch message.Text {
		case "Войти":
			err = h.authHandler.HandleLogin(message)
		case "Зарегистрироваться":
			err = h.authHandler.HandleRegister(message)
		case "Выйти":
			err = h.authHandler.HandleLogout(message)
		case "Мой профиль":
			// Здесь будет обработка профиля пользователя
			err = h.client.SendMessage(message.Chat.ID, "Функция просмотра профиля в разработке.")
		case "Пополнить баланс":
			// Здесь будет обработка пополнения баланса
			err = h.client.SendMessage(message.Chat.ID, "Функция пополнения баланса в разработке.")
		case "Список пользователей":
			// Здесь будет обработка списка пользователей
			err = h.client.SendMessage(message.Chat.ID, "Функция просмотра списка пользователей в разработке.")
		case "Управление пользователями":
			// Здесь будет обработка управления пользователями
			err = h.client.SendMessage(message.Chat.ID, "Функция управления пользователями в разработке.")
		default:
			// Здесь будет обработка других команд
			err = h.client.SendMessage(message.Chat.ID, "Неизвестная команда. Используйте кнопки для навигации.")
		}
	}

	if err != nil {
		log.Printf("Error handling message: %v", err)
	}
}
