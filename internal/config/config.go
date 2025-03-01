package config

import (
	"os"
	"strconv"
	"time"
)

// Config содержит конфигурацию приложения
type Config struct {
	TelegramToken string
	DBPath        string
	PollTimeout   time.Duration
	MessagesLimit int
	Debug         bool
	JWTSecret     string        // Секретный ключ для JWT токенов
	JWTExpiration time.Duration // Время жизни JWT токена
}

// NewConfig создает новый экземпляр Config
func NewConfig() *Config {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		token = "7828331860:AAG_XkEaE2vY4EKdGZaOJ9xD74D1fVV0U_k" // Значение по умолчанию
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "users.db" // Значение по умолчанию
	}

	// Проверяем режим отладки
	debug := false
	debugEnv := os.Getenv("DEBUG")
	if debugEnv != "" {
		debug, _ = strconv.ParseBool(debugEnv)
	}

	// Устанавливаем таймаут
	timeout := 60 * time.Second // Увеличиваем таймаут до 60 секунд
	timeoutEnv := os.Getenv("POLL_TIMEOUT")
	if timeoutEnv != "" {
		if t, err := strconv.Atoi(timeoutEnv); err == nil {
			timeout = time.Duration(t) * time.Second
		}
	}

	// Получаем секретный ключ для JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-256-bit-secret" // Значение по умолчанию, в продакшене следует использовать сложный ключ
	}

	// Устанавливаем время жизни JWT токена (по умолчанию 24 часа)
	jwtExpiration := 24 * time.Hour
	jwtExpirationEnv := os.Getenv("JWT_EXPIRATION")
	if jwtExpirationEnv != "" {
		if exp, err := strconv.Atoi(jwtExpirationEnv); err == nil {
			jwtExpiration = time.Duration(exp) * time.Hour
		}
	}

	return &Config{
		TelegramToken: token,
		DBPath:        dbPath,
		PollTimeout:   timeout,
		MessagesLimit: 100,
		Debug:         debug,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExpiration,
	}
}
