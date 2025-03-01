package domain

// UserRepository определяет методы для работы с пользователями в хранилище
type UserRepository interface {
	// GetByID возвращает пользователя по его ChatID
	GetByID(chatID int64) (*User, error)

	// GetByUsername возвращает пользователя по его имени пользователя
	GetByUsername(username string) (*User, error)

	// Save сохраняет или обновляет пользователя
	Save(user *User) error

	// Delete удаляет пользователя
	Delete(chatID int64) error

	// GetAll возвращает всех пользователей
	GetAll() ([]*User, error)

	// UpdateBalance обновляет баланс пользователя
	UpdateBalance(chatID int64, newBalance float64) error

	// UpdatePassword обновляет пароль пользователя
	UpdatePassword(chatID int64, newPassword string) error

	// UpdateRole обновляет роль пользователя
	UpdateRole(chatID int64, newRole string) error
}

// PaymentRepository определяет методы для работы с платежами
type PaymentRepository interface {
	// SaveRequest сохраняет запрос на пополнение баланса
	SaveRequest(chatID int64, amount float64) error

	// GetRequests возвращает все запросы на пополнение баланса
	GetRequests() ([]*PaymentRequest, error)

	// UpdateRequestStatus обновляет статус запроса
	UpdateRequestStatus(id int64, status string) error
}
