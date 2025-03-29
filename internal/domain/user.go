package domain

import "time"

// User представляет пользователя в системе
type User struct {
	ChatID    int64     `json:"chat_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Position  string    `json:"position"`
	Birthday  string    `json:"birthday"`
	Number    string    `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserState представляет состояние пользователя в диалоге с ботом
type UserState int

const (
	StateNone UserState = iota
	StateAwaitingUsername
	StateAwaitingPassword
)

// UserSession представляет текущую сессию пользователя
type UserSession struct {
	User         *User
	State        UserState
	IsAuthorized bool
	LastCommand  string // Последняя команда пользователя (login/register)
	Token        string // JWT токен для авторизации
}

// Константы для ролей пользователей
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
