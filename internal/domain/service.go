package domain

// UserService определяет методы для работы с пользователями
type UserService interface {
	// GetUser возвращает пользователя по его ChatID
	GetUser(chatID int64) (*User, error)

	// SaveUser сохраняет или обновляет пользователя
	SaveUser(user *User) error

	// DeleteUser удаляет пользователя
	DeleteUser(chatID int64) error

	// GetAllUsers возвращает всех пользователей
	GetAllUsers() ([]*User, error)

	// UpdateUserBalance обновляет баланс пользователя
	UpdateUserBalance(chatID int64, amount float64) error

	// UpdateUserProfile обновляет профиль пользователя
	UpdateUserProfile(chatID int64, position, birthday, number string) error
}

// PaymentService определяет методы для работы с платежами
type PaymentService interface {
	// CreatePaymentRequest создает запрос на пополнение баланса
	CreatePaymentRequest(chatID int64, amount float64) error

	// ConfirmPayment подтверждает платеж и обновляет баланс пользователя
	ConfirmPayment(chatID int64, amount float64) error
}

// SessionService определяет методы для работы с сессиями пользователей
type SessionService interface {
	// GetSession возвращает текущую сессию пользователя
	GetSession(chatID int64) (*UserSession, error)

	// UpdateSession обновляет сессию пользователя
	UpdateSession(chatID int64, session *UserSession) error

	// DeleteSession удаляет сессию пользователя
	DeleteSession(chatID int64) error

	// Login авторизует пользователя и обновляет его сессию
	Login(chatID int64, username, password string) error

	// Logout выходит из системы и удаляет сессию пользователя
	Logout(chatID int64) error

	// Register регистрирует нового пользователя
	Register(user *User) error

	// IsAdmin проверяет, является ли пользователь администратором
	IsAdmin(chatID int64) (bool, error)
}

// AuthService определяет методы для авторизации пользователей
type AuthService interface {
	// Register регистрирует нового пользователя
	Register(user *User) error

	// Login авторизует пользователя
	Login(username, password string) (*User, error)

	// ChangePassword изменяет пароль пользователя
	ChangePassword(chatID int64, oldPassword, newPassword string) error

	// ChangeRole изменяет роль пользователя (только для администраторов)
	ChangeRole(adminChatID int64, targetUsername string, newRole string) error
}
