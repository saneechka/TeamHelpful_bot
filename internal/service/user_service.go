package service

import (
	"fmt"
	"time"

	"HelpBot/internal/domain"
)

// UserService реализует интерфейс domain.UserService
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUser возвращает пользователя по его ChatID
func (s *UserService) GetUser(chatID int64) (*domain.User, error) {
	return s.userRepo.GetByID(chatID)
}

// SaveUser сохраняет или обновляет пользователя
func (s *UserService) SaveUser(user *domain.User) error {
	return s.userRepo.Save(user)
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(chatID int64) error {
	return s.userRepo.Delete(chatID)
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}

// UpdateUserBalance обновляет баланс пользователя
func (s *UserService) UpdateUserBalance(chatID int64, amount float64) error {
	user, err := s.userRepo.GetByID(chatID)
	if err != nil {
		return err
	}

	if user == nil {
		// Создаем нового пользователя, если он не существует
		username := fmt.Sprintf("user_%d", chatID)
		user = &domain.User{
			ChatID:    chatID,
			Username:  username,
			Balance:   amount,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return s.userRepo.Save(user)
	}

	// Обновляем баланс существующего пользователя
	newBalance := user.Balance + amount
	return s.userRepo.UpdateBalance(chatID, newBalance)
}

// UpdateUserProfile обновляет профиль пользователя
func (s *UserService) UpdateUserProfile(chatID int64, position, birthday, number string) error {
	user, err := s.userRepo.GetByID(chatID)
	if err != nil {
		return err
	}

	if user == nil {
		// Создаем нового пользователя, если он не существует
		username := fmt.Sprintf("user_%d", chatID)
		user = &domain.User{
			ChatID:    chatID,
			Username:  username,
			Position:  position,
			Birthday:  birthday,
			Number:    number,
			Balance:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else {
		// Обновляем существующего пользователя
		user.Position = position
		user.Birthday = birthday
		user.Number = number
		user.UpdatedAt = time.Now()
	}

	return s.userRepo.Save(user)
}
