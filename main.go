package main

import (
	"log"
	"time"

	"HelpBot/client/telegram"
)

const (
	tgHost         = "api.telegram.org"
	pollTimeout    = 100 * time.Millisecond
	messagesLimit  = 100
)

var predefinedResponse = "В 19:10"

func main() {
	token := "7828331860:AAG_XkEaE2vY4EKdGZaOJ9xD74D1fVV0U_k" // Replace with your bot token

	client := telegram.NewClient(tgHost, token)

	offset := 0
	for {
		updates, err := client.Updates(offset, messagesLimit)
		if err != nil {
			log.Printf("Error getting updates: %v", err)
			continue
		}

		for _, update := range updates {
			if update.Message == nil {
				continue
			}

			 // Send predefined message instead of echo
			err = client.SendMessage(update.Message.Chat.ID, predefinedResponse)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

			offset = update.ID + 1
		}

		time.Sleep(pollTimeout)
	}
}