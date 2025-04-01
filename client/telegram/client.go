package telegram

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"HelpBot/internal/domain"
)

// Client представляет клиент для работы с Telegram API
type Client struct {
	bot           *tgbotapi.BotAPI
	updateTimeout time.Duration
	messagesLimit int
}

// NewClient создает новый экземпляр клиента Telegram
func NewClient(token string, updateTimeout time.Duration, messagesLimit int) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// Включаем режим отладки
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Client{
		bot:           bot,
		updateTimeout: updateTimeout,
		messagesLimit: messagesLimit,
	}, nil
}

// StartPolling начинает получение обновлений от Telegram API
func (c *Client) StartPolling(handler func(update *tgbotapi.Update)) {
	// Сначала удаляем вебхуки, если они были установлены
	log.Println("Deleting webhook...")
	_, err := c.bot.Request(tgbotapi.DeleteWebhookConfig{})
	if err != nil {
		log.Printf("Error deleting webhook: %v", err)
	}

	log.Println("Starting to receive updates...")
	// Начинаем получать обновления с ID 0 (получаем все доступные обновления)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(c.updateTimeout.Seconds())
	updates := c.bot.GetUpdatesChan(u)

	log.Printf("Bot started polling for updates with timeout %d seconds", int(c.updateTimeout.Seconds()))

	for update := range updates {
		log.Printf("Received update ID: %d", update.UpdateID)
		if update.Message != nil {
			log.Printf("Message from %s: %s", update.Message.From.UserName, update.Message.Text)
		}
		handler(&update)
	}
}

// SendMessage отправляет сообщение пользователю
func (c *Client) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	return err
}

// SendMessageWithKeyboard отправляет сообщение с клавиатурой
//переписать на дженерики
func (c *Client) SendMessageWithKeyboard(chatID int64, text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)

	// Проверяем тип клавиатуры
	switch k := keyboard.(type) {
	case ReplyKeyboardMarkup:
		msg.ReplyMarkup = k.ToTelegramKeyboard()
	case tgbotapi.ReplyKeyboardMarkup:
		msg.ReplyMarkup = k
	case tgbotapi.InlineKeyboardMarkup:
		msg.ReplyMarkup = k
	default:
		msg.ReplyMarkup = keyboard
	}

	_, err := c.bot.Send(msg)
	return err
}

// CreateReplyKeyboard создает клавиатуру с кнопками
func (c *Client) CreateReplyKeyboard(buttons [][]string) tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.KeyboardButton
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewKeyboardButton(button))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewReplyKeyboard(keyboard...)
}

// CreateInlineKeyboard создает инлайн-клавиатуру с кнопками
func (c *Client) CreateInlineKeyboard(buttons [][]InlineButton) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.InlineKeyboardButton
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// RemoveKeyboard удаляет клавиатуру
func (c *Client) RemoveKeyboard(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := c.bot.Send(msg)
	return err
}

// GetUserFromMessage извлекает информацию о пользователе из сообщения
func (c *Client) GetUserFromMessage(message *tgbotapi.Message) *domain.User {
	if message == nil || message.From == nil {
		return nil
	}

	// Генерируем уникальное имя пользователя, добавляя chat_id
	username := fmt.Sprintf("%s_%d", message.From.UserName, message.Chat.ID)

	return &domain.User{
		ChatID:    message.Chat.ID,
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetLoginKeyboard возвращает клавиатуру для авторизации
func (c *Client) GetLoginKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := [][]string{
		{"Войти"},
		{"Зарегистрироваться"},
	}
	return c.CreateReplyKeyboard(buttons)
}

// GetMainMenuKeyboard возвращает клавиатуру главного меню
func (c *Client) GetMainMenuKeyboard(isAdmin bool) tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]string
	if isAdmin {
		buttons = [][]string{
			{"Мой профиль", "Пополнить баланс"},
			{"Список пользователей", "Управление пользователями"},
			{"Выйти"},
		}
	} else {
		buttons = [][]string{
			{"Мой профиль", "Пополнить баланс"},
			{"Выйти"},
		}
	}
	return c.CreateReplyKeyboard(buttons)
}

// GetUserManagementKeyboard возвращает клавиатуру управления пользователями
func (c *Client) GetUserManagementKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttons := [][]string{
		{"Добавить пользователя", "Удалить пользователя"},
		{"Изменить роль пользователя"},
		{"Назад"},
	}
	return c.CreateReplyKeyboard(buttons)
}
