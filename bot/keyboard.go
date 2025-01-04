package bot

import "HelpBot/client/telegram"

func CreateMainKeyboard() telegram.ReplyKeyboardMarkup {
    return telegram.ReplyKeyboardMarkup{
        Keyboard: [][]telegram.KeyboardButton{
            {{Text: "Мой баланс"}},
        },
        ResizeKeyboard: true,
        OneTimeKeyboard: false,
    }
}
