package bot

import (
	"HelpBot/client/telegram"
	_"log"
)

const (
	WelcomeMessage = "Выберите действие:"
	BalanceMessage = "Ваш текущий баланс: 1000 руб."
)

type Handler struct {
	client   *telegram.Client
	keyboard telegram.ReplyKeyboardMarkup
}

func NewHandler(client *telegram.Client) *Handler {
	return &Handler{
		client:   client,
		keyboard: CreateMainKeyboard(),
	}
}

func (h *Handler) HandleUpdate(update telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	switch update.Message.Text {
	case "/start":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, WelcomeMessage, h.keyboard)
	case "Мой баланс":
		return h.client.SendMessage(update.Message.Chat.ID, BalanceMessage)
	default:
		return nil // Ignore all other messages
	}
}
