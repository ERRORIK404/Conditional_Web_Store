package service

import (
	"errors"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/hashing"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/localJWT"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/repository"

	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register создает нового пользователя с хешированием пароля
func (s *AuthService) Register(login, password string) (*models.User, error) {
	// Проверка существования пользователя
	existing, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Хеширование пароля
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Создание пользователя
	newUser := &models.User{
		ID:             &uuid.UUID{},
		Login:          login,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}
	id, err := s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = id
	return newUser, nil
}

// Login проверяет учетные данные и генерирует JWT токен
func (s *AuthService) Login(login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user == nil {
		logger.Logger.Error("user not found")
		return "", errors.New("user not found")
	}

	// Проверка пароля
	if isCorrect := hashing.IsCorrectPassword(user.HashedPassword, password); !isCorrect {
		logger.Logger.Error("invalid password")
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT токена
	token, err := localjwt.GenerateJWT(user.ID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
