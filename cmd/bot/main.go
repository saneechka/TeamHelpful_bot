package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgclient "HelpBot/client/telegram"
	"HelpBot/internal/config"
	tgdelivery "HelpBot/internal/delivery/telegram"
	"HelpBot/internal/repository"
	"HelpBot/internal/repository/sqlite"
	"HelpBot/internal/service"
)

func main() {
	// Инициализируем конфигурацию
	cfg := config.NewConfig()

	// Настраиваем логирование
	if cfg.Debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("Debug mode enabled")
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// Инициализируем базу данных
	db, err := sqlite.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Инициализируем репозитории
	userRepo := sqlite.NewUserRepository(db)
	paymentRepo := sqlite.NewPaymentRepository(db)

	// Создаем репозитории
	repos := repository.NewRepositories(userRepo, paymentRepo)

	// Инициализируем сервисы
	userService := service.NewUserService(repos.UserRepository)
	authService := service.NewAuthService(repos.UserRepository, cfg)
	sessionService := service.NewSessionService(userService, authService)
	paymentService := service.NewPaymentService(repos.PaymentRepository, userService)

	// Инициализируем клиент Telegram
	client, err := tgclient.NewClient(cfg.TelegramToken, cfg.PollTimeout, cfg.MessagesLimit)
	if err != nil {
		log.Fatalf("Error creating Telegram client: %v", err)
	}

	// Инициализируем обработчик
	handler := tgdelivery.NewHandler(client, userService, paymentService, sessionService)

	log.Printf("Bot started with poll timeout: %v", cfg.PollTimeout)

	// Настраиваем обработку сигналов для корректного завершения
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Запускаем обработку сообщений в отдельной горутине
	go client.StartPolling(handler.HandleUpdate)

	// Ожидаем сигнал завершения
	<-c
	log.Println("Shutting down bot...")
}
