package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

type Database struct {
    db *sql.DB
}

type UserSession struct {
    ChatID    int64
    Username  string
    LoginTime time.Time
    Position  string
    Birthday  string
    Number    string    // Добавляем поле для игрового номера
    Balance   float64
}

func NewDatabase(dbPath string) (*Database, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if (err != nil) {
        return nil, err
    }

    // Создаем таблицы, включая таблицу для мастер-пароля
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_sessions (
            chat_id INTEGER PRIMARY KEY,
            username TEXT NOT NULL,
            login_time DATETIME NOT NULL,
            position TEXT DEFAULT '',
            birthday TEXT DEFAULT '',
            number TEXT DEFAULT '',
            balance REAL DEFAULT 0.0
        );

        CREATE TABLE IF NOT EXISTS user_credentials (
            username TEXT PRIMARY KEY,
            password TEXT NOT NULL,
            created_at DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS master_password (
            password TEXT NOT NULL,
            created_at DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS payment_requests (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            chat_id INTEGER NOT NULL,
            amount REAL NOT NULL,
            status TEXT DEFAULT 'pending',
            created_at DATETIME NOT NULL
        );
    `)
    if err != nil {
        return nil, err
    }

    return &Database{db: db}, nil
}

func (d *Database) SaveUserSession(chatID int64, username string, position string, birthday string, number string, balance float64) error {
    _, err := d.db.Exec(`
        INSERT OR REPLACE INTO user_sessions (chat_id, username, login_time, position, birthday, number, balance)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, chatID, username, time.Now(), position, birthday, number, balance)
    return err
}

// Add method to update balance
func (d *Database) UpdateUserBalance(chatID int64, newBalance float64) error {
    _, err := d.db.Exec(`UPDATE user_sessions SET balance = ? WHERE chat_id = ?`, newBalance, chatID)
    return err
}

func (d *Database) RemoveUserSession(chatID int64) error {
    _, err := d.db.Exec(`DELETE FROM user_sessions WHERE chat_id = ?`, chatID)
    return err
}

func (d *Database) GetUserSession(chatID int64) (*UserSession, error) {
    var session UserSession
    err := d.db.QueryRow(`
        SELECT chat_id, username, login_time, position, birthday, number, balance
        FROM user_sessions 
        WHERE chat_id = ?
    `, chatID).Scan(&session.ChatID, &session.Username, &session.LoginTime, &session.Position, &session.Birthday, &session.Number, &session.Balance)
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (d *Database) GetAllActiveUsers() ([]UserSession, error) {
    rows, err := d.db.Query(`
        SELECT chat_id, username, login_time, position, birthday, number, balance
        FROM user_sessions 
        ORDER BY login_time DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []UserSession
    for rows.Next() {
        var user UserSession
        if err := rows.Scan(&user.ChatID, &user.Username, &user.LoginTime, &user.Position, &user.Birthday, &user.Number, &user.Balance); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func (d *Database) UpdateBalance(chatID int64, newBalance float64) error {
    _, err := d.db.Exec(`
        UPDATE user_sessions SET balance = ? WHERE chat_id = ?
    `, newBalance, chatID)
    return err
}

// Добавляем методы для работы с паролями
func (d *Database) SetUserPassword(username, password string) error {
    _, err := d.db.Exec(`
        INSERT OR REPLACE INTO user_credentials (username, password, created_at)
        VALUES (?, ?, ?)
    `, username, password, time.Now())
    return err
}

func (d *Database) CheckUserPassword(username, password string) (bool, error) {
    var storedPassword string
    err := d.db.QueryRow(`
        SELECT password FROM user_credentials WHERE username = ?
    `, username).Scan(&storedPassword)

    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    return storedPassword == password, nil
}

func (d *Database) CheckUserExists(username string) (bool, error) {
    var exists bool
    err := d.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM user_credentials WHERE username = ?)
    `, username).Scan(&exists)
    
    if err != nil {
        return false, err
    }
    return exists, nil
}

func (d *Database) SetMasterPassword(password string) error {
    _, err := d.db.Exec(`
        DELETE FROM master_password;
        INSERT INTO master_password (password, created_at)
        VALUES (?, ?)
    `, password, time.Now())
    return err
}

func (d *Database) GetMasterPassword() (string, error) {
    var password string
    err := d.db.QueryRow(`
        SELECT password FROM master_password LIMIT 1
    `).Scan(&password)

    if err == sql.ErrNoRows {
        return "", nil
    }
    if err != nil {
        return "", err
    }

    return password, nil
}

func (d *Database) SavePaymentRequest(chatID int64, amount float64) error {
    _, err := d.db.Exec(`
        INSERT INTO payment_requests (chat_id, amount, created_at)
        VALUES (?, ?, ?)
    `, chatID, amount, time.Now())
    return err
}

func (d *Database) Close() error {
    return d.db.Close()
}
