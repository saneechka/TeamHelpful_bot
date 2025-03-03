package repository

import (
	"HelpBot/internal/domain"
)

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	domain.UserRepository
}

// Repositories содержит все репозитории
type Repositories struct {
	UserRepository domain.UserRepository
}

// NewRepositories создает новый экземпляр Repositories
func NewRepositories(userRepo domain.UserRepository) *Repositories {
	return &Repositories{
		UserRepository: userRepo,
	}
}
