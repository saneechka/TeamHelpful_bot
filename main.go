package main

import (
	"log"
	"os"
	"time"

	"HelpBot/bot"
	"HelpBot/client/telegram"
)

const (
	pollTimeout    = 100 * time.Millisecond
	messagesLimit  = 100
)

func main() {
	// Initialize token directly, remove tgHost
	token := "bot_token"
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