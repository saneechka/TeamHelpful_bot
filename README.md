# TeamHelpful Bot

Telegram-бот с системой регистрации и авторизации пользователей.

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
│   │       └── user_repository.go
│   ├── service               # Реализация сервисов
│   │   ├── auth_service.go
│   │   ├── session_service.go
│   │   └── user_service.go
│   └── delivery              # Обработка запросов
│       └── telegram
│           ├── handler.go
│           └── keyboard.go
├── Dockerfile               # Сборка Docker образа
├── docker-compose.yml      # Конфигурация Docker Compose
├── go.mod                  # Зависимости
└── README.md               # Документация
```

## Функциональность

- Регистрация новых пользователей
- Авторизация существующих пользователей
- JWT-авторизация
- Роли пользователей (admin/user)

## Запуск

### Требования

- Go 1.21 или выше
- SQLite
- Docker и Docker Compose (для контейнеризации)

### Переменные окружения

Создайте файл `.env` со следующими переменными:

```env
# Обязательные переменные
BOT_TOKEN=your_telegram_bot_token    # Токен вашего Telegram бота
DB_PATH=users.db                     # Путь к файлу базы данных

# Опциональные переменные
JWT_SECRET=your_secret_key           # Секретный ключ для JWT (если не указан, генерируется автоматически)
JWT_EXPIRATION=24                    # Время жизни JWT токена в часах (по умолчанию: 24)
DEBUG=false                          # Режим отладки (по умолчанию: false)
POLL_TIMEOUT=60                      # Таймаут опроса в секундах (по умолчанию: 60)
```

### Локальный запуск

```bash
# Сборка
make build

# Запуск
make run

# Запуск без сборки
make run-dev
```

### Запуск с Docker Compose

1. Создайте файл `.env` с необходимыми переменными окружения (как описано выше)
2. Запустите контейнер:
```bash
docker-compose up -d
```

## Команды бота

- `/start` - Начать работу с ботом
- `/help` - Получить справку
- `/login` - Войти в систему
- `/register` - Зарегистрироваться
- `/logout` - Выйти из системы

## Безопасность

<<<<<<< HEAD
### Предварительные требования
- Go 1.19+
- SQLite3
- Telegram Bot Token
- Docker и Docker Compose (для контейнеризации)

### Пошаговая установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/saneechka/teamHelpful.git
cd teamHelpfulBot
```

2. Создайте файл .env с необходимыми переменными окружения:
```
BOT_TOKEN=your_telegram_bot_token
DB_PATH=users.db
```

3. Установите зависимости (для локальной разработки):
```bash
go mod download
```

4. Создайте базу данных:
```bash
make init-db
```

5. Запустите бота локально:
```bash
go run main.go
```

### Запуск с использованием Docker Compose

1. Убедитесь, что у вас установлены Docker и Docker Compose:
```bash
docker --version
docker-compose --version
```

2. Создайте файл .env с необходимыми переменными окружения (как описано выше).

3. Запустите приложение с помощью Docker Compose:
```bash
./deploy.sh
```

или вручную:

```bash
docker-compose up -d --build
```

4. Проверьте статус контейнера:
```bash
docker-compose ps
```

5. Просмотр логов:
```bash
docker-compose logs -f
```

6. Остановка контейнера:
```bash
docker-compose down
```

## Деплой на продакшн

### Настройка GitHub Actions

Для автоматического деплоя на продакшн-сервер при пуше в ветку main используется GitHub Actions. Для настройки:

1. Добавьте следующие секреты в настройках вашего GitHub репозитория (Settings > Secrets and variables > Actions):

   - `BOT_TOKEN`: Токен вашего Telegram бота
   - `SSH_PRIVATE_KEY`: Приватный SSH-ключ для доступа к серверу
   - `SSH_HOST`: IP-адрес или домен вашего сервера
   - `SSH_USERNAME`: Имя пользователя для SSH-подключения
   - `DOCKER_USERNAME`: Имя пользователя Docker Hub
   - `DOCKER_TOKEN`: Токен доступа Docker Hub
   - `SLACK_WEBHOOK` (опционально): URL для отправки уведомлений в Slack

2. Убедитесь, что на вашем сервере установлены Docker и Docker Compose.

3. При пуше в ветку main GitHub Actions автоматически:
   - Соберет Docker-образ
   - Загрузит его в Docker Hub
   - Подключится к вашему серверу по SSH
   - Развернет приложение с использованием Docker Compose
   - Проверит успешность деплоя
   - Отправит уведомление о результате (если настроен Slack Webhook)

### Ручной деплой на продакшн

Вы также можете запустить workflow деплоя вручную через интерфейс GitHub Actions (вкладка Actions > Production Deployment > Run workflow).

## API Методы


- `GetBalance(userID)` - Получение баланса
- `MakePayment(userID, amount)` - Внесение оплаты
- `UpdateProfile(userID, data)` - Обновление профиля

### Административные
- `ConfirmPayment(paymentID)` - Подтверждение платежа
- `ManageTeam(action, data)` - Управление командой
- `SendNotification(userID, message)` - Отправка уведомлений
=======
### Хранение паролей
- Пароли хешируются с использованием bcrypt с оптимальной стоимостью
- Исходные пароли никогда не сохраняются в базе данных
- При каждом хешировании автоматически генерируется уникальная соль

### JWT авторизация
- Используются токены с ограниченным временем жизни
- Секретный ключ генерируется автоматически при запуске (если не указан)
- Поддерживается ротация ключей через переменные окружения
>>>>>>> Test

### База данных
- Защита от SQL-инъекций через подготовленные выражения
- Безопасное хранение чувствительных данных
- Автоматическое обновление временных меток

### Общие меры безопасности
- Валидация всех входных данных
- Ограничение попыток входа
- Безопасное хранение токенов и секретов
- Логирование важных событий безопасности

## Поддержка

- GitHub Issues: [создать issue](https://github.com/saneechka/teamHelpful_bot/issues)
- Telegram: [@bell1matita](https://t.me/bell1matita)
- Email: alexandro.dev.work@gmail.com



