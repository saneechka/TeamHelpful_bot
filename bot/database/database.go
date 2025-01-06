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

    // Создаем таблицу с новой структурой
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_sessions_new (
            chat_id INTEGER PRIMARY KEY,
            username TEXT NOT NULL,
            login_time DATETIME NOT NULL,
            position TEXT DEFAULT '',
            birthday TEXT DEFAULT '',
            number TEXT DEFAULT '',
            balance REAL DEFAULT 0.0
        );
        
        INSERT OR IGNORE INTO user_sessions_new (chat_id, username, login_time, position, birthday, number)
        SELECT chat_id, username, login_time, position, birthday, number FROM user_sessions;
        
        DROP TABLE IF EXISTS user_sessions;
        
        ALTER TABLE user_sessions_new RENAME TO user_sessions;
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

func (d *Database) Close() error {
    return d.db.Close()
}
