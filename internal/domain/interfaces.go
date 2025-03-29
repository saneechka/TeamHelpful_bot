package domain

// UserRepository определяет методы для работы с пользователями в БД
type UserRepository interface {
	// GetByID возвращает пользователя по его ChatID
	GetByID(chatID int64) (*User, error)

	// GetByUsername возвращает пользователя по его имени пользователя
	GetByUsername(username string) (*User, error)

	// Save сохраняет нового пользователя
	Save(user *User) error

	// Update обновляет существующего пользователя
	Update(user *User) error

	// Delete удаляет пользователя
	Delete(chatID int64) error

	// GetAll возвращает всех пользователей
	GetAll() ([]*User, error)

	// UpdatePassword обновляет пароль пользователя
	UpdatePassword(chatID int64, newPassword string) error

	// UpdateRole обновляет роль пользователя
	UpdateRole(chatID int64, newRole string) error
}

// UserService определяет методы для работы с пользователями
type UserService interface {
	// GetUser возвращает пользователя по его ChatID
	GetUser(chatID int64) (*User, error)

	// GetUserByUsername возвращает пользователя по его имени пользователя
	GetUserByUsername(username string) (*User, error)

	// SaveUser сохраняет или обновляет пользователя
	SaveUser(user *User) error

	// DeleteUser удаляет пользователя
	DeleteUser(chatID int64) error

	// GetAllUsers возвращает всех пользователей
	GetAllUsers() ([]*User, error)

	// UpdateUserProfile обновляет профиль пользователя
	UpdateUserProfile(chatID int64, position, birthday, number string) error
}

// SessionService определяет методы для работы с сессиями пользователей
type SessionService interface {
	// GetSession возвращает сессию пользователя по его ChatID
	GetSession(chatID int64) (*UserSession, error)

	// UpdateSession обновляет сессию пользователя
	UpdateSession(chatID int64, session *UserSession) error

	// Login аутентифицирует пользователя
	Login(chatID int64, username, password string) error

	// Logout выходит из системы пользователя
	Logout(chatID int64) error

	// Register регистрирует нового пользователя
	Register(user *User) error

	// ValidateToken проверяет токен пользователя
	ValidateToken(chatID int64, token string) error

	// IsAdmin проверяет, является ли пользователь администратором
	IsAdmin(chatID int64) (bool, error)
}
