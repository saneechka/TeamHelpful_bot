package bot

import (
	"HelpBot/client/telegram"
	"sync"
)

const (
	WelcomeMessage = "Выберите действие:"
	BalanceMessage = "Ваш текущий баланс примерно ничего."
	PaymentMessage = "Выберите способ оплаты:" 
	AccountMessage = "Информация о вашем  персональном аккаунте"
	PaymentOption1 = "Нажмите 'Произвести оплату' для продолжения"
	PaymentOption2 = "Нажмите 'Произвести оплату' для продолжения"
	TeamMessage    =  "Сборище тех кто в пустые не забивает"
	ProcessingPayment = "Обработка оплаты..."
	AlreadyProcessing = "Оплата уже обрабатывается. Пожалуйста, подождите."
)

type Handler struct {
	client          *telegram.Client
	keyboard        telegram.ReplyKeyboardMarkup
	paymentKeyboard telegram.ReplyKeyboardMarkup
	processKeyboard telegram.ReplyKeyboardMarkup
	processingUsers map[int64]bool
	mu              sync.RWMutex
}

func NewHandler(client *telegram.Client) *Handler {
	return &Handler{
		client:          client,
		keyboard:        CreateMainKeyboard(),
		paymentKeyboard: CreatePaymentKeyboard(),
		processKeyboard: CreatePaymentProcessKeyboard(),
		processingUsers: make(map[int64]bool),
	}
}

func (h *Handler) isProcessingPayment(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.processingUsers[userID]
}

func (h *Handler) setProcessingPayment(userID int64, status bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if status {
		h.processingUsers[userID] = true
	} else {
		delete(h.processingUsers, userID)
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
	case "Способ оплаты":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentMessage, h.paymentKeyboard)
	case "Мой аккаунт":
		return h.client.SendMessage(update.Message.Chat.ID, AccountMessage)
	case "Информация о команде":
		return h.client.SendMessage(update.Message.Chat.ID,TeamMessage) 
	case "1":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentOption1, h.processKeyboard)
	case "2":
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentOption2, h.processKeyboard)
	case "Произвести оплату":
		if h.isProcessingPayment(update.Message.Chat.ID) {
			return h.client.SendMessage(update.Message.Chat.ID, AlreadyProcessing)
		}
		h.setProcessingPayment(update.Message.Chat.ID, true)
		return h.client.SendMessage(update.Message.Chat.ID, ProcessingPayment)
	case "Назад к способам оплаты":
		h.setProcessingPayment(update.Message.Chat.ID, false)
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, PaymentMessage, h.paymentKeyboard)
	case "Назад":
		h.setProcessingPayment(update.Message.Chat.ID, false)
		return h.client.SendMessageWithKeyboard(update.Message.Chat.ID, WelcomeMessage, h.keyboard)
	default:
		return nil 
	}
}
