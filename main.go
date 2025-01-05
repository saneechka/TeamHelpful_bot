package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"HelpBot/bot"
	"HelpBot/client/telegram"
)

const (
	tgHost         = "api.telegram.org"
	pollTimeout    = 100 * time.Millisecond
	messagesLimit  = 100
	maxRetries     = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	token := "7828331860:AAG_XkEaE2vY4EKdGZaOJ9xD74D1fVV0U_k"
	client := telegram.NewClient(tgHost, token)
	handler := bot.NewHandler(&client)

	// Start bot in separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		runBot(ctx, client, handler)
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down bot...")
	cancel()
	wg.Wait()
	log.Println("Bot shutdown complete")
}

func runBot(ctx context.Context, client telegram.Client, handler *bot.Handler) {
	offset := 0
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			updates, err := client.Updates(offset, messagesLimit)
			if err != nil {
				log.Printf("Error getting updates: %v", err)
				if retryCount < maxRetries {
					retryCount++
					time.Sleep(time.Second * time.Duration(retryCount))
					continue
				}
				log.Printf("Max retries reached, restarting bot")
				retryCount = 0
				offset = 0
				continue
			}
			retryCount = 0

			for _, update := range updates {
				if err := handler.HandleUpdate(update); err != nil {
					log.Printf("Error handling update: %v", err)
				}
				offset = update.ID + 1
			}

			time.Sleep(pollTimeout)
		}
	}
}