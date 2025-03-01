package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// NewDB создает новое подключение к базе данных SQLite
func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables создает необходимые таблицы в базе данных
func createTables(db *sql.DB) error {
	// Создаем таблицу пользователей с полями для авторизации
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_sessions (
			chat_id INTEGER PRIMARY KEY,
			username TEXT,
			password TEXT,
			role TEXT DEFAULT 'user',
			position TEXT,
			birthday TEXT,
			number TEXT,
			balance REAL DEFAULT 0,
			login_time TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Error creating user_sessions table: %v", err)
		return err
	}

	// Создаем таблицу платежей
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER,
			amount REAL,
			status TEXT,
			created_at TIMESTAMP,
			FOREIGN KEY (chat_id) REFERENCES user_sessions (chat_id)
		)
	`)
	if err != nil {
		log.Printf("Error creating payments table: %v", err)
		return err
	}

	return nil
}
