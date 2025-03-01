package service

import (
	"errors"
	"fmt"
	"time"

	"HelpBot/internal/config"
	"HelpBot/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// AuthService реализует интерфейс domain.AuthService
type AuthService struct {
	userRepo domain.UserRepository
	config   *config.Config
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo domain.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// hashPassword хеширует пароль с использованием bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash проверяет, соответствует ли пароль хешу
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(user *domain.User) error {
	// Проверяем, существует ли пользователь с таким именем
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("пользователь с таким именем уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}
	user.Password = hashedPassword

	// Устанавливаем роль по умолчанию
	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	// Сохраняем пользователя
	return s.userRepo.Save(user)
}

// Login авторизует пользователя
func (s *AuthService) Login(username, password string) (*domain.User, error) {
	// Получаем пользователя по имени
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем пароль
	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("неверный пароль")
	}

	// Обновляем время последнего входа
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword изменяет пароль пользователя
func (s *AuthService) ChangePassword(chatID int64, oldPassword, newPassword string) error {
	// Получаем пользователя
	user, err := s.userRepo.GetByID(chatID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем старый пароль
	if !checkPasswordHash(oldPassword, user.Password) {
		return errors.New("неверный пароль")
	}

	// Хешируем новый пароль
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	// Обновляем пароль
	return s.userRepo.UpdatePassword(chatID, hashedPassword)
}

// ChangeRole изменяет роль пользователя (только для администраторов)
func (s *AuthService) ChangeRole(adminChatID int64, targetUsername string, newRole string) error {
	// Проверяем, является ли пользователь администратором
	admin, err := s.userRepo.GetByID(adminChatID)
	if err != nil {
		return err
	}
	if admin == nil || admin.Role != domain.RoleAdmin {
		return errors.New("недостаточно прав для выполнения операции")
	}

	// Получаем целевого пользователя
	targetUser, err := s.userRepo.GetByUsername(targetUsername)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем валидность роли
	if newRole != domain.RoleAdmin && newRole != domain.RoleUser {
		return errors.New("недопустимая роль")
	}

	// Обновляем роль
	return s.userRepo.UpdateRole(targetUser.ChatID, newRole)
}
