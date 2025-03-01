package domain

import "time"

// User представляет пользователя системы
type User struct {
	ChatID    int64
	Username  string
	Password  string
	Role      string
	Position  string
	Birthday  string
	Number    string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserState представляет состояние пользователя в диалоге с ботом
type UserState int

const (
	StateNone UserState = iota
	StateAwaitingPosition
	StateAwaitingBirthday
	StateAwaitingNumber
	StateAwaitingPaymentAmount
	StateAwaitingPaymentConfirmation
	StateAwaitingUsername
	StateAwaitingPassword
)

// UserSession представляет текущую сессию пользователя
type UserSession struct {
	User          *User
	State         UserState
	PaymentAmount float64
	IsAuthorized  bool
	LastCommand   string // Последняя команда пользователя (login/register)
	Token         string // JWT токен для авторизации
}

// PaymentRequest представляет запрос на пополнение баланса
type PaymentRequest struct {
	ID        int64
	ChatID    int64
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// Роли пользователей
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
