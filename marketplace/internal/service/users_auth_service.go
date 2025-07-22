package service

import (
	"errors"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/hashing"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/localJWT"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/repository"
	"go.uber.org/zap"

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
		logger.Logger.Error("users_auth_service.Register", zap.Error(err))
		return nil, err
	}
	if existing != nil {
		logger.Logger.Error("users_auth_service.Register user already exists", zap.Any("id", existing.ID), zap.String("login", existing.Login))
		return nil, errors.New("user already exists")
	}

	// Хеширование пароля
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Создание пользователя
	newUser := models.User{
		ID:             uuid.Nil,
		Login:          login,
		HashedPassword: string(hashedPassword),
	}
	id, err := s.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = id

	logger.Logger.Debug("new user on level service", zap.Any("id", id), zap.String("login", login))
	return &newUser, nil
}

// Login проверяет учетные данные и генерирует JWT токен
func (s *AuthService) Login(login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return "", err
	}
	if user == nil {
		logger.Logger.Error("users_auth_service.Login user not found")
		return "", errors.New("user not found")
	}

	// Проверка пароля
	if isCorrect := hashing.IsCorrectPassword(user.HashedPassword, password); !isCorrect {
		logger.Logger.Error("users_auth_service.Login invalid password")
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT токена
	token, err := localjwt.GenerateJWT(user.ID.String())
	if err != nil {
		return "", err
	}

	logger.Logger.Debug("users_auth_service.Login user login with jwt", zap.String("token", token), zap.String("login", login))
	return token, nil
}
