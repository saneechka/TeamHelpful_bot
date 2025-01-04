package telegram

type Update struct {
    ID      int      `json:"update_id"`
    Message *Message `json:"message"`
}

type Message struct {
    Text string `json:"text"`
    Chat Chat   `json:"chat"`
}

type Chat struct {
    ID int64 `json:"id"`
}

type KeyboardButton struct {
    Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
    Keyboard              [][]KeyboardButton `json:"keyboard"`
    ResizeKeyboard        bool               `json:"resize_keyboard"`
    OneTimeKeyboard       bool               `json:"one_time_keyboard"`
    InputFieldPlaceholder string             `json:"input_field_placeholder"`
    IsPersistent          bool               `json:"is_persistent"`
}