package service

import (
	"errors"
	"fmt"
	"time"

	"HelpBot/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims представляет собой данные, которые будут храниться в JWT токене
type JWTClaims struct {
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken генерирует JWT токен для пользователя
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	// Получаем текущее время
	now := time.Now()

	// Создаем claims для токена
	claims := JWTClaims{
		ChatID:   user.ChatID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "helpbot",
			Subject:   fmt.Sprintf("%d", user.ChatID),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает пользователя
func (s *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Получаем claims из токена
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Получаем пользователя из базы данных
	user, err := s.userRepo.GetByID(claims.ChatID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Проверяем, что имя пользователя и роль совпадают
	if user.Username != claims.Username || user.Role != claims.Role {
		return nil, errors.New("token data mismatch")
	}

	return user, nil
}
