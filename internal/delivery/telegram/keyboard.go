package telegram

import (
	"HelpBot/client/telegram"
	"HelpBot/internal/domain"
)

// CreateMainKeyboard создает основную клавиатуру
func CreateMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: "Мой баланс"}, {Text: "Способ оплаты"}, {Text: "Мой аккаунт"}},
			{{Text: "Информация о команде"}},
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

// CreatePositionKeyboard создает клавиатуру для выбора позиции
func CreatePositionKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: domain.PositionForward}},
			{{Text: domain.PositionDefender}},
			{{Text: domain.PositionGoalie}},
			{{Text: "Отмена"}},
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
