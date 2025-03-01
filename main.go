package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"HelpBot/bot"
	"HelpBot/client/telegram"

	"github.com/joho/godotenv"
)

const (
	pollTimeout   = 100 * time.Millisecond
	messagesLimit = 100
)

// startHealthCheckServer starts a simple HTTP server for health checks
func startHealthCheckServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go func() {
		log.Println("Starting health check server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()
}

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded: %v", err)
	}

	// Start health check server
	startHealthCheckServer()

	token := "bot_token" // Значение по умолчанию
	if envToken := os.Getenv("BOT_TOKEN"); envToken != "" {
		token = envToken
	}

	// Create client with just token
	client := telegram.NewClient(token)
	handler, err := bot.NewHandler(client, "users.db")
	if err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	offset := 0
	for {
		updates, err := client.Updates(offset, messagesLimit)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			continue
		}

		for _, update := range updates {
			if err := handler.HandleUpdate(update); err != nil {
				log.Printf("Error handling update: %v", err)
			}
			offset = update.ID + 1
		}

		time.Sleep(pollTimeout)
	}
}
