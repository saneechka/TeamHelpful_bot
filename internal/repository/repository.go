package repository

import (
	"HelpBot/internal/domain"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	domain.UserRepository
}

// PaymentRepository определяет методы для работы с платежами
type PaymentRepository interface {
	domain.PaymentRepository
}

// Repositories содержит все репозитории
type Repositories struct {
	UserRepository    UserRepository
	PaymentRepository PaymentRepository
}

// NewRepositories создает новый экземпляр Repositories
func NewRepositories(userRepo UserRepository, paymentRepo PaymentRepository) *Repositories {
	return &Repositories{
		UserRepository:    userRepo,
		PaymentRepository: paymentRepo,
	}
}
