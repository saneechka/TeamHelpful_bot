package bot

import "HelpBot/client/telegram"

func CreateMainKeyboard() telegram.ReplyKeyboardMarkup {
    return telegram.ReplyKeyboardMarkup{
        Keyboard: [][]telegram.KeyboardButton{
            {{Text: "Мой баланс"}, {Text: "Способ оплаты"},
        {Text:"Мой аккаунт"},
        },
        },
        ResizeKeyboard: true,
        OneTimeKeyboard: true,
    }
}
