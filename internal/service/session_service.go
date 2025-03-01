package service

import (
	"sync"

	"HelpBot/internal/domain"
)

// SessionService реализует интерфейс domain.SessionService
type SessionService struct {
	userService domain.UserService
	authService *AuthService
	sessions    map[int64]*domain.UserSession
	mu          sync.RWMutex
}

// NewSessionService создает новый экземпляр SessionService
func NewSessionService(userService domain.UserService, authService *AuthService) *SessionService {
	return &SessionService{
		userService: userService,
		authService: authService,
		sessions:    make(map[int64]*domain.UserSession),
	}
}

// GetSession возвращает текущую сессию пользователя
func (s *SessionService) GetSession(chatID int64) (*domain.UserSession, error) {
	s.mu.RLock()
	session, exists := s.sessions[chatID]
	s.mu.RUnlock()

	if exists {
		return session, nil
	}

	// Если сессия не найдена, пытаемся получить пользователя из БД
	user, err := s.userService.GetUser(chatID)
	if err != nil {
		return nil, err
	}

	// Если пользователь не найден, возвращаем nil
	if user == nil {
		return nil, nil
	}

	// Создаем новую сессию
	session = &domain.UserSession{
		User:         user,
		State:        domain.StateNone,
		IsAuthorized: false, // По умолчанию пользователь не авторизован
	}

	// Сохраняем сессию
	s.mu.Lock()
	s.sessions[chatID] = session
	s.mu.Unlock()

	return session, nil
}

// UpdateSession обновляет сессию пользователя
func (s *SessionService) UpdateSession(chatID int64, session *domain.UserSession) error {
	s.mu.Lock()
	s.sessions[chatID] = session
	s.mu.Unlock()

	return nil
}

// DeleteSession удаляет сессию пользователя
func (s *SessionService) DeleteSession(chatID int64) error {
	s.mu.Lock()
	delete(s.sessions, chatID)
	s.mu.Unlock()

	return nil
}

// Login авторизует пользователя и обновляет его сессию
func (s *SessionService) Login(chatID int64, username, password string) error {
	// Авторизуем пользователя
	user, err := s.authService.Login(username, password)
	if err != nil {
		return err
	}

	// Генерируем JWT токен
	token, err := s.authService.GenerateToken(user)
	if err != nil {
		return err
	}

	// Получаем текущую сессию
	session, err := s.GetSession(chatID)
	if err != nil {
		return err
	}

	// Если сессия не существует, создаем новую
	if session == nil {
		session = &domain.UserSession{
			User:         user,
			State:        domain.StateNone,
			IsAuthorized: true,
			Token:        token,
		}
	} else {
		// Обновляем пользователя в сессии
		session.User = user
		session.State = domain.StateNone
		session.IsAuthorized = true
		session.Token = token
	}

	// Сохраняем сессию
	return s.UpdateSession(chatID, session)
}

// Logout выходит из системы и удаляет сессию пользователя
func (s *SessionService) Logout(chatID int64) error {
	return s.DeleteSession(chatID)
}

// Register регистрирует нового пользователя
func (s *SessionService) Register(user *domain.User) error {
	return s.authService.Register(user)
}

// IsAdmin проверяет, является ли пользователь администратором
func (s *SessionService) IsAdmin(chatID int64) (bool, error) {
	session, err := s.GetSession(chatID)
	if err != nil {
		return false, err
	}
	if session == nil || session.User == nil {
		return false, nil
	}
	return session.User.Role == domain.RoleAdmin, nil
}

// ValidateToken проверяет валидность JWT токена и обновляет сессию
func (s *SessionService) ValidateToken(chatID int64, tokenString string) error {
	// Проверяем токен
	user, err := s.authService.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// Получаем текущую сессию
	session, err := s.GetSession(chatID)
	if err != nil {
		return err
	}

	// Если сессия не существует, создаем новую
	if session == nil {
		session = &domain.UserSession{
			User:         user,
			State:        domain.StateNone,
			IsAuthorized: true,
			Token:        tokenString,
		}
	} else {
		// Обновляем пользователя в сессии
		session.User = user
		session.State = domain.StateNone
		session.IsAuthorized = true
		session.Token = tokenString
	}

	// Сохраняем сессию
	return s.UpdateSession(chatID, session)
}
