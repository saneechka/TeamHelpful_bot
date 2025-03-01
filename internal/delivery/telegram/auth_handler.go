package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"HelpBot/client/telegram"
	"HelpBot/internal/domain"
)

// AuthHandler обрабатывает команды авторизации
type AuthHandler struct {
	client         *telegram.Client
	sessionService domain.SessionService
	userService    domain.UserService
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(client *telegram.Client, sessionService domain.SessionService, userService domain.UserService) *AuthHandler {
	return &AuthHandler{
		client:         client,
		sessionService: sessionService,
		userService:    userService,
	}
}

// HandleStart обрабатывает команду /start
func (h *AuthHandler) HandleStart(message *tgbotapi.Message) error {
	// Получаем сессию пользователя
	session, err := h.sessionService.GetSession(message.Chat.ID)
	if err != nil {
		return err
	}

	// Если пользователь не существует, создаем его
	if session == nil {
		user := h.client.GetUserFromMessage(message)
		if err := h.userService.SaveUser(user); err != nil {
			return err
		}

		// Создаем новую сессию
		session = &domain.UserSession{
			User:        user,
			State:       domain.StateNone,
			LastCommand: "", // Сбрасываем последнюю команду
		}
	} else {
		// Сбрасываем состояние существующей сессии
		session.State = domain.StateNone
		// Сбрасываем имя пользователя в сессии, чтобы начать процесс заново
		session.User.Username = ""
		// Сбрасываем последнюю команду
		session.LastCommand = ""
	}

	// Обновляем сессию
	if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
		return err
	}

	// Если пользователь уже авторизован, показываем главное меню
	if session.IsAuthorized {
		isAdmin, err := h.sessionService.IsAdmin(message.Chat.ID)
		if err != nil {
			return err
		}

		keyboard := h.client.GetMainMenuKeyboard(isAdmin)
		return h.client.SendMessageWithKeyboard(message.Chat.ID, "Добро пожаловать! Выберите действие:", keyboard)
	}

	// Если пользователь не авторизован, показываем меню авторизации
	keyboard := h.client.GetLoginKeyboard()
	return h.client.SendMessageWithKeyboard(message.Chat.ID, "Добро пожаловать! Для начала работы необходимо авторизоваться:", keyboard)
}

// HandleLogin обрабатывает процесс входа в систему
func (h *AuthHandler) HandleLogin(message *tgbotapi.Message) error {
	// Получаем сессию пользователя
	session, err := h.sessionService.GetSession(message.Chat.ID)
	if err != nil {
		return err
	}
	if session == nil {
		// Если сессия не существует, создаем новую
		user := h.client.GetUserFromMessage(message)
		if err := h.userService.SaveUser(user); err != nil {
			return err
		}

		session = &domain.UserSession{
			User:        user,
			State:       domain.StateNone,
			LastCommand: "login", // Устанавливаем последнюю команду как "login"
		}
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}
	}

	// В зависимости от текущего состояния сессии
	switch session.State {
	case domain.StateNone:
		// Если текст сообщения "Войти", начинаем процесс входа
		if message.Text == "Войти" {
			// Запрашиваем имя пользователя
			session.State = domain.StateAwaitingUsername
			session.LastCommand = "login" // Устанавливаем последнюю команду как "login"
			if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
				return err
			}
			return h.client.SendMessage(message.Chat.ID, "Введите имя пользователя:")
		}
		return nil

	case domain.StateAwaitingUsername:
		// Сохраняем имя пользователя и запрашиваем пароль
		username := message.Text
		if username == "" {
			return h.client.SendMessage(message.Chat.ID, "Имя пользователя не может быть пустым. Попробуйте еще раз:")
		}

		// Проверяем, существует ли пользователь с таким именем
		users, err := h.userService.GetAllUsers()
		if err != nil {
			return err
		}

		userExists := false
		for _, u := range users {
			if u.Username == username {
				userExists = true
				break
			}
		}

		if !userExists {
			// Если пользователь не существует, предлагаем зарегистрироваться
			keyboard := h.client.GetLoginKeyboard()
			return h.client.SendMessageWithKeyboard(message.Chat.ID, "Пользователь с таким именем не найден. Выберите действие:", keyboard)
		}

		session.User.Username = username
		session.State = domain.StateAwaitingPassword
		session.LastCommand = "login" // Устанавливаем последнюю команду как "login"
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}
		return h.client.SendMessage(message.Chat.ID, "Введите пароль:")

	case domain.StateAwaitingPassword:
		// Пытаемся авторизовать пользователя
		password := message.Text
		if password == "" {
			return h.client.SendMessage(message.Chat.ID, "Пароль не может быть пустым. Попробуйте еще раз:")
		}

		// Авторизуем пользователя
		err := h.sessionService.Login(message.Chat.ID, session.User.Username, password)
		if err != nil {
			// Если произошла ошибка авторизации, сбрасываем состояние и предлагаем выбрать действие
			session.State = domain.StateNone
			session.User.Username = "" // Сбрасываем имя пользователя
			if updateErr := h.sessionService.UpdateSession(message.Chat.ID, session); updateErr != nil {
				return updateErr
			}
			keyboard := h.client.GetLoginKeyboard()
			return h.client.SendMessageWithKeyboard(message.Chat.ID, fmt.Sprintf("Ошибка авторизации: %s. Выберите действие:", err.Error()), keyboard)
		}

		// Получаем обновленную сессию
		session, err = h.sessionService.GetSession(message.Chat.ID)
		if err != nil {
			return err
		}

		// Сбрасываем состояние сессии
		session.State = domain.StateNone
		session.LastCommand = "" // Сбрасываем последнюю команду
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}

		// Показываем главное меню
		isAdmin, err := h.sessionService.IsAdmin(message.Chat.ID)
		if err != nil {
			return err
		}

		keyboard := h.client.GetMainMenuKeyboard(isAdmin)
		return h.client.SendMessageWithKeyboard(message.Chat.ID, "Вы успешно авторизованы! Выберите действие:", keyboard)

	default:
		return fmt.Errorf("неизвестное состояние сессии")
	}
}

// HandleRegister обрабатывает процесс регистрации
func (h *AuthHandler) HandleRegister(message *tgbotapi.Message) error {
	// Получаем сессию пользователя
	session, err := h.sessionService.GetSession(message.Chat.ID)
	if err != nil {
		return err
	}
	if session == nil {
		// Если сессия не существует, создаем новую
		user := h.client.GetUserFromMessage(message)
		if err := h.userService.SaveUser(user); err != nil {
			return err
		}

		session = &domain.UserSession{
			User:        user,
			State:       domain.StateNone,
			LastCommand: "register", // Устанавливаем последнюю команду как "register"
		}
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}
	}

	// В зависимости от текущего состояния сессии
	switch session.State {
	case domain.StateNone:
		// Если текст сообщения "Зарегистрироваться", начинаем процесс регистрации
		if message.Text == "Зарегистрироваться" {
			// Запрашиваем имя пользователя
			session.State = domain.StateAwaitingUsername
			session.LastCommand = "register" // Устанавливаем последнюю команду как "register"
			if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
				return err
			}
			return h.client.SendMessage(message.Chat.ID, "Введите имя пользователя для регистрации:")
		}
		return nil

	case domain.StateAwaitingUsername:
		// Сохраняем имя пользователя и запрашиваем пароль
		username := message.Text
		if username == "" {
			return h.client.SendMessage(message.Chat.ID, "Имя пользователя не может быть пустым. Попробуйте еще раз:")
		}

		// Проверяем, существует ли пользователь с таким именем
		users, err := h.userService.GetAllUsers()
		if err != nil {
			return err
		}

		for _, u := range users {
			if u.Username == username {
				return h.client.SendMessage(message.Chat.ID, "Пользователь с таким именем уже существует. Введите другое имя:")
			}
		}

		session.User.Username = username
		session.State = domain.StateAwaitingPassword
		session.LastCommand = "register" // Устанавливаем последнюю команду как "register"
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}
		return h.client.SendMessage(message.Chat.ID, "Введите пароль для регистрации:")

	case domain.StateAwaitingPassword:
		// Регистрируем пользователя
		password := message.Text
		if password == "" {
			return h.client.SendMessage(message.Chat.ID, "Пароль не может быть пустым. Попробуйте еще раз:")
		}

		// Создаем нового пользователя
		user := &domain.User{
			ChatID:    message.Chat.ID,
			Username:  session.User.Username,
			Password:  password,
			Role:      domain.RoleUser,
			CreatedAt: session.User.CreatedAt,
			UpdatedAt: session.User.UpdatedAt,
		}

		// Регистрируем пользователя
		if err := h.sessionService.Register(user); err != nil {
			// Если произошла ошибка регистрации, сбрасываем состояние и предлагаем выбрать действие
			session.State = domain.StateNone
			session.User.Username = "" // Сбрасываем имя пользователя
			session.LastCommand = ""   // Сбрасываем последнюю команду
			if updateErr := h.sessionService.UpdateSession(message.Chat.ID, session); updateErr != nil {
				return updateErr
			}
			keyboard := h.client.GetLoginKeyboard()
			return h.client.SendMessageWithKeyboard(message.Chat.ID, fmt.Sprintf("Ошибка регистрации: %s. Выберите действие:", err.Error()), keyboard)
		}

		// Авторизуем пользователя
		if err := h.sessionService.Login(message.Chat.ID, user.Username, password); err != nil {
			// Если произошла ошибка авторизации, сбрасываем состояние и предлагаем выбрать действие
			session.State = domain.StateNone
			session.User.Username = "" // Сбрасываем имя пользователя
			session.LastCommand = ""   // Сбрасываем последнюю команду
			if updateErr := h.sessionService.UpdateSession(message.Chat.ID, session); updateErr != nil {
				return updateErr
			}
			keyboard := h.client.GetLoginKeyboard()
			return h.client.SendMessageWithKeyboard(message.Chat.ID, fmt.Sprintf("Регистрация успешна, но произошла ошибка авторизации: %s. Выберите действие:", err.Error()), keyboard)
		}

		// Получаем обновленную сессию
		session, err = h.sessionService.GetSession(message.Chat.ID)
		if err != nil {
			return err
		}

		// Сбрасываем состояние сессии
		session.State = domain.StateNone
		session.LastCommand = "" // Сбрасываем последнюю команду
		if err := h.sessionService.UpdateSession(message.Chat.ID, session); err != nil {
			return err
		}

		// Показываем главное меню
		isAdmin, err := h.sessionService.IsAdmin(message.Chat.ID)
		if err != nil {
			return err
		}

		keyboard := h.client.GetMainMenuKeyboard(isAdmin)
		return h.client.SendMessageWithKeyboard(message.Chat.ID, "Вы успешно зарегистрированы и авторизованы! Выберите действие:", keyboard)

	default:
		return fmt.Errorf("неизвестное состояние сессии")
	}
}

// HandleLogout обрабатывает выход из системы
func (h *AuthHandler) HandleLogout(message *tgbotapi.Message) error {
	// Удаляем сессию пользователя
	if err := h.sessionService.Logout(message.Chat.ID); err != nil {
		return err
	}

	// Показываем меню авторизации
	keyboard := h.client.GetLoginKeyboard()
	return h.client.SendMessageWithKeyboard(message.Chat.ID, "Вы вышли из системы. Для продолжения работы необходимо авторизоваться:", keyboard)
}

// HandleToken обрабатывает авторизацию по JWT токену
func (h *AuthHandler) HandleToken(message *tgbotapi.Message) error {
	// Получаем токен из сообщения
	token := message.Text
	if token == "" {
		return h.client.SendMessage(message.Chat.ID, "Токен не может быть пустым. Попробуйте еще раз:")
	}

	// Проверяем токен и обновляем сессию
	err := h.sessionService.ValidateToken(message.Chat.ID, token)
	if err != nil {
		// Если произошла ошибка валидации токена, сбрасываем состояние и предлагаем выбрать действие
		session, _ := h.sessionService.GetSession(message.Chat.ID)
		if session != nil {
			session.State = domain.StateNone
			session.LastCommand = ""
			h.sessionService.UpdateSession(message.Chat.ID, session)
		}
		keyboard := h.client.GetLoginKeyboard()
		return h.client.SendMessageWithKeyboard(message.Chat.ID, fmt.Sprintf("Ошибка валидации токена: %s. Выберите действие:", err.Error()), keyboard)
	}

	// Показываем главное меню
	isAdmin, err := h.sessionService.IsAdmin(message.Chat.ID)
	if err != nil {
		return err
	}

	keyboard := h.client.GetMainMenuKeyboard(isAdmin)
	return h.client.SendMessageWithKeyboard(message.Chat.ID, "Вы успешно авторизованы по токену! Выберите действие:", keyboard)
}
