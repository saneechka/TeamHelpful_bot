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
	// Проверяем, существует ли пользователь
	existingUser, err := s.userRepo.GetByID(user.ChatID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		// Если пользователь не существует, создаем нового
		return s.userRepo.Save(user)
	} else {
		// Если пользователь существует, обновляем его
		return s.userRepo.Update(user)
	}
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(chatID int64) error {
	return s.userRepo.Delete(chatID)
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
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

// GetUserByUsername возвращает пользователя по его имени пользователя
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(username)
}
