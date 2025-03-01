package sqlite

import (
	"database/sql"
	"time"

	"HelpBot/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

// UserRepository реализует интерфейс domain.UserRepository для SQLite
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID возвращает пользователя по его ChatID
func (r *UserRepository) GetByID(chatID int64) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT chat_id, username, password, role, position, birthday, number, balance, login_time, login_time
		FROM user_sessions 
		WHERE chat_id = ?
	`, chatID).Scan(
		&user.ChatID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Position,
		&user.Birthday,
		&user.Number,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername возвращает пользователя по его имени пользователя
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT chat_id, username, password, role, position, birthday, number, balance, login_time, login_time
		FROM user_sessions 
		WHERE username = ?
	`, username).Scan(
		&user.ChatID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Position,
		&user.Birthday,
		&user.Number,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save сохраняет или обновляет пользователя
func (r *UserRepository) Save(user *domain.User) error {
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	_, err := r.db.Exec(`
		INSERT OR REPLACE INTO user_sessions (
			chat_id, username, password, role, position, birthday, number, balance, login_time
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		user.ChatID,
		user.Username,
		user.Password,
		user.Role,
		user.Position,
		user.Birthday,
		user.Number,
		user.Balance,
		user.CreatedAt,
	)
	return err
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(chatID int64) error {
	_, err := r.db.Exec(`DELETE FROM user_sessions WHERE chat_id = ?`, chatID)
	return err
}

// GetAll возвращает всех пользователей
func (r *UserRepository) GetAll() ([]*domain.User, error) {
	rows, err := r.db.Query(`
		SELECT chat_id, username, password, role, position, birthday, number, balance, login_time, login_time
		FROM user_sessions 
		ORDER BY login_time DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ChatID,
			&user.Username,
			&user.Password,
			&user.Role,
			&user.Position,
			&user.Birthday,
			&user.Number,
			&user.Balance,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// UpdateBalance обновляет баланс пользователя
func (r *UserRepository) UpdateBalance(chatID int64, newBalance float64) error {
	_, err := r.db.Exec(`
		UPDATE user_sessions SET balance = ? WHERE chat_id = ?
	`, newBalance, chatID)
	return err
}

// UpdatePassword обновляет пароль пользователя
func (r *UserRepository) UpdatePassword(chatID int64, newPassword string) error {
	_, err := r.db.Exec(`
		UPDATE user_sessions SET password = ? WHERE chat_id = ?
	`, newPassword, chatID)
	return err
}

// UpdateRole обновляет роль пользователя
func (r *UserRepository) UpdateRole(chatID int64, newRole string) error {
	_, err := r.db.Exec(`
		UPDATE user_sessions SET role = ? WHERE chat_id = ?
	`, newRole, chatID)
	return err
}
