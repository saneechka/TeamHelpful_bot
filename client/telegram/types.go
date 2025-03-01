package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InlineButton представляет кнопку инлайн-клавиатуры
type InlineButton struct {
	Text string
	Data string
}

// CallbackData представляет данные обратного вызова
type CallbackData struct {
	Action string
	Params map[string]string
}

// AuthStatus представляет статус авторизации
type AuthStatus int

const (
	// AuthStatusNone - пользователь не авторизован
	AuthStatusNone AuthStatus = iota
	// AuthStatusPending - пользователь в процессе авторизации
	AuthStatusPending
	// AuthStatusAuthorized - пользователь авторизован
	AuthStatusAuthorized
)

// KeyboardButton представляет кнопку клавиатуры
type KeyboardButton struct {
	Text string
}

// ReplyKeyboardMarkup представляет клавиатуру с кнопками
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton
	ResizeKeyboard  bool
	OneTimeKeyboard bool
	Selective       bool
}

// ToTelegramKeyboard конвертирует ReplyKeyboardMarkup в tgbotapi.ReplyKeyboardMarkup
func (k ReplyKeyboardMarkup) ToTelegramKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var keyboard [][]tgbotapi.KeyboardButton

	for _, row := range k.Keyboard {
		var keyboardRow []tgbotapi.KeyboardButton
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewKeyboardButton(button.Text))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        keyboard,
		ResizeKeyboard:  k.ResizeKeyboard,
		OneTimeKeyboard: k.OneTimeKeyboard,
		Selective:       k.Selective,
	}
}
