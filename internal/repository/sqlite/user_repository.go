package sqlite

import (
	"database/sql"
	"fmt"
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
		SELECT chat_id, username, password, role, position, birthday, number, created_at, updated_at 
		FROM users 
		WHERE chat_id = ?`, chatID).Scan(
		&user.ChatID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Position,
		&user.Birthday,
		&user.Number,
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
		SELECT chat_id, username, password, role, position, birthday, number, created_at, updated_at 
		FROM users 
		WHERE username = ?`, username).Scan(
		&user.ChatID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Position,
		&user.Birthday,
		&user.Number,
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

// Save сохраняет нового пользователя
func (r *UserRepository) Save(user *domain.User) error {
	// Создаем нового пользователя
	_, err := r.db.Exec(`
		INSERT INTO users (chat_id, username, password, role, position, birthday, number, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ChatID,
		user.Username,
		user.Password,
		user.Role,
		user.Position,
		user.Birthday,
		user.Number,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(chatID int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE chat_id = ?", chatID)
	return err
}

// UpdatePassword обновляет пароль пользователя
func (r *UserRepository) UpdatePassword(chatID int64, newPassword string) error {
	_, err := r.db.Exec(`
		UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE chat_id = ?
	`, newPassword, chatID)
	return err
}

// UpdateRole обновляет роль пользователя
func (r *UserRepository) UpdateRole(chatID int64, newRole string) error {
	_, err := r.db.Exec(`
		UPDATE users SET role = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE chat_id = ?
	`, newRole, chatID)
	return err
}

// GetAll возвращает всех пользователей
func (r *UserRepository) GetAll() ([]*domain.User, error) {
	rows, err := r.db.Query(`
		SELECT chat_id, username, password, role, position, birthday, number, created_at, updated_at 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ChatID,
			&user.Username,
			&user.Password,
			&user.Role,
			&user.Position,
			&user.Birthday,
			&user.Number,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Update обновляет существующего пользователя
func (r *UserRepository) Update(user *domain.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET username = ?, password = ?, role = ?, position = ?, birthday = ?, number = ?, updated_at = ?
		WHERE chat_id = ?`,
		user.Username,
		user.Password,
		user.Role,
		user.Position,
		user.Birthday,
		user.Number,
		time.Now(),
		user.ChatID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
