package telegram

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CreateMainKeyboard создает основную клавиатуру
func CreateMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Войти"}},
			{{Text: "Зарегистрироваться"}},
		},
		ResizeKeyboard: true,
	}
}

// CreateAuthenticatedKeyboard создает клавиатуру для авторизованного пользователя
func CreateAuthenticatedKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Выйти"}},
		},
		ResizeKeyboard: true,
	}
}

// CreateCancelKeyboard создает клавиатуру с кнопкой отмены
func CreateCancelKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Отмена"}},
		},
		ResizeKeyboard: true,
	}
}

// CreatePaymentKeyboard создает клавиатуру для выбора способа оплаты
func CreatePaymentKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "1"}, {Text: "2"}},
			{{Text: "Пополнить баланс"}},
			{{Text: "Назад"}},
		},
		ResizeKeyboard: true,
	}
}

// CreatePaymentProcessKeyboard создает клавиатуру для процесса оплаты
func CreatePaymentProcessKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Произвести оплату"}},
			{{Text: "Подтвердить оплату"}},
			{{Text: "Отменить"}},
			{{Text: "Назад к способам оплаты"}},
		},
		ResizeKeyboard: true,
	}
}

// CreateTeamKeyboard создает клавиатуру для информации о команде
func CreateTeamKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Состав"}},
			{{Text: "Назад"}},
		},
		ResizeKeyboard: true,
	}
}

// CreateLoginKeyboard создает клавиатуру для авторизации
func CreateLoginKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Войти"}},
			{{Text: "Зарегистрироваться"}},
		},
		ResizeKeyboard: true,
	}
}
