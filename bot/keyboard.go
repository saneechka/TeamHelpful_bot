package bot

import "HelpBot/client/telegram"

func CreateMainKeyboard() telegram.ReplyKeyboardMarkup {
    return telegram.ReplyKeyboardMarkup{
        Keyboard: [][]telegram.KeyboardButton{
            {{Text: "Мой баланс"}, {Text: "Способ оплаты"}, {Text: "Мой аккаунт"}},
            {{Text: "Информация о команде"}},
        },
        ResizeKeyboard: true,
        OneTimeKeyboard: false,
    }
}

func CreatePaymentKeyboard() telegram.ReplyKeyboardMarkup {
    return telegram.ReplyKeyboardMarkup{
        Keyboard: [][]telegram.KeyboardButton{
            {{Text: "1"}, {Text: "2"}},
            {{Text: "Назад"}},
        },
        ResizeKeyboard: true,
        OneTimeKeyboard: false,
    }
}

func CreatePaymentProcessKeyboard() telegram.ReplyKeyboardMarkup {
    return telegram.ReplyKeyboardMarkup{
        Keyboard: [][]telegram.KeyboardButton{
            {{Text: "Произвести оплату"}},
            {{Text: "Назад к способам оплаты"}},
        },
        ResizeKeyboard: true,
        OneTimeKeyboard: false,
    }
}
