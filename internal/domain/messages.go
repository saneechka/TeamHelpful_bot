package domain

// Константы для сообщений бота
const (
	WelcomeMessage        = "Выберите действие:"
	BalanceMessage        = "Ваш текущий баланс:"
	PaymentMessage        = "Выберите способ оплаты:"
	AccountMessage        = "Информация о вашем персональном аккаунте"
	PaymentOption1        = "Нажмите 'Произвести оплату' для продолжения"
	PaymentOption2        = "Нажмите 'Произвести оплату' для продолжения"
	TeamMessage           = "Сборище тех кто в пустые не забивает"
	ProcessingPayment     = "Обработка оплаты..."
	AlreadyProcessing     = "Оплата уже обрабатывается. Пожалуйста, подождите."
	AskPositionMessage    = "Введите вашу позицию в команде:"
	AskBirthdayMessage    = "Введите вашу дату рождения (например, 01.01.1990):"
	ProfileSetupComplete  = "Информация сохранена!"
	AskNumberMessage      = "Введите ваш игровой номер:"
	PositionForward       = "Нападающий"
	PositionDefender      = "Защитник"
	PositionGoalie        = "Вратарь"
	TeamRosterMessage     = "Состав команды:"
	EnterAmountMessage    = "Введите сумму для пополнения (в рублях):"
	PaymentDetailsMessage = "Для пополнения баланса переведите указанную сумму на счет:\n" +
		"Номер карты: 1234 5678 9012 3456\n" +
		"После перевода нажмите 'Подтвердить оплату'"
	PaymentConfirmationMessage = "Оплата будет проверена администратором в течение 24 часов"
	InvalidAmountMessage       = "Пожалуйста, введите корректную сумму (целое число больше 0)"
)
