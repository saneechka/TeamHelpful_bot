package service

import (
	"HelpBot/internal/domain"
)

// PaymentService реализует интерфейс domain.PaymentService
type PaymentService struct {
	paymentRepo domain.PaymentRepository
	userService domain.UserService
}

// NewPaymentService создает новый экземпляр PaymentService
func NewPaymentService(paymentRepo domain.PaymentRepository, userService domain.UserService) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		userService: userService,
	}
}

// CreatePaymentRequest создает запрос на пополнение баланса
func (s *PaymentService) CreatePaymentRequest(chatID int64, amount float64) error {
	return s.paymentRepo.SaveRequest(chatID, amount)
}

// ConfirmPayment подтверждает платеж и обновляет баланс пользователя
func (s *PaymentService) ConfirmPayment(chatID int64, amount float64) error {
	// Обновляем баланс пользователя
	if err := s.userService.UpdateUserBalance(chatID, amount); err != nil {
		return err
	}

	return nil
}
