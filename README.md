# TeamHelpful Bot

Telegram-бот для управления командой с возможностью пополнения баланса и просмотра информации о команде.

## Архитектура проекта

Проект построен с использованием принципов чистой архитектуры и разделения ответственности:

- **Domain** - модели и бизнес-логика
- **Repository** - работа с данными
- **Service** - бизнес-логика
- **Delivery** - обработка запросов от Telegram

## Структура проекта

```
├── cmd
│   └── bot
│       └── main.go           # Точка входа в приложение
├── client
│   └── telegram              # Клиент для работы с Telegram API
│       ├── client.go
│       └── types.go
├── internal
│   ├── config                # Конфигурация приложения
│   │   └── config.go
│   ├── domain                # Модели и интерфейсы
│   │   ├── messages.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── user.go
│   ├── repository            # Реализация репозиториев
│   │   └── sqlite
│   │       ├── db.go
│   │       ├── payment_repository.go
│   │       └── user_repository.go
│   ├── service               # Реализация сервисов
│   │   ├── payment_service.go
│   │   ├── session_service.go
│   │   └── user_service.go
│   └── delivery              # Обработка запросов
│       └── telegram
│           ├── handler.go
│           └── keyboard.go
├── Makefile                  # Команды для сборки и запуска
├── go.mod                    # Зависимости
└── README.md                 # Документация
```

## Функциональность

- Просмотр баланса
- Пополнение баланса
- Управление профилем пользователя
- Просмотр информации о команде

## Запуск

### Требования

- Go 1.21 или выше
- SQLite

### Переменные окружения

- `BOT_TOKEN` - токен Telegram бота
- `DB_PATH` - путь к файлу базы данных SQLite (по умолчанию: `users.db`)

### Команды

```bash
# Сборка
make build

# Запуск
make run

# Запуск без сборки
make run-dev

# Очистка
make clean

# Тестирование
make test
```

