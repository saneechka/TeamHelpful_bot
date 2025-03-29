#!/bin/bash

# Проверяем наличие sqlite3
if ! command -v sqlite3 &> /dev/null; then
    echo "Error: sqlite3 is not installed"
    exit 1
fi

# Путь к базе данных (можно переопределить через переменную окружения)
DB_PATH=${DB_PATH:-"users.db"}

# Удаляем старую базу данных, если она существует
if [ -f "$DB_PATH" ]; then
    echo "Removing existing database..."
    rm "$DB_PATH"
fi

echo "Creating new database..."

# Создаем таблицу пользователей
sqlite3 "$DB_PATH" <<EOF
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    position TEXT,
    birthday TEXT,
    number TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_users_chat_id ON users(chat_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
EOF

# Устанавливаем права доступа
chmod 600 "$DB_PATH"

echo "Database initialized successfully at $DB_PATH" 