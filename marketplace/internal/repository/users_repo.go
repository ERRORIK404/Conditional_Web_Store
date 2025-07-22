package repository

import (
	"errors"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user models.User) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByLogin(login string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(user models.User) (uuid.UUID, error) {
	if err := r.db.Create(user).Error; err != nil {
		logger.Logger.Error("UserRepository.Create", zap.Error(err))
		return uuid.Nil, err
	}
	
	logger.Logger.Debug("New Users info on level repo ",
		zap.String("id", user.ID.String()),
		zap.String("login", user.Login))
	return user.ID, nil
}

// GetByID возвращает пользователя по ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("UserRepository.GetByID(NOT FOUND)", zap.Error(err))
			return nil, nil
		}
		logger.Logger.Error("UserRepository.GetByID", zap.Error(err))
		return nil, err
	}

	logger.Logger.Debug("GetedByID User info on level repo",
		zap.String("id", user.ID.String()),
		zap.String("login", user.Login))
	return &user, nil
}

// GetByLogin возвращает пользователя по логину
func (r *userRepository) GetByLogin(login string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("UserRepository.GetByLogin(NOT FOUND)", zap.Error(err))
			return nil, nil
		}
		logger.Logger.Error("UserRepository.GetByLogin", zap.Error(err))
		return nil, err
	}

	logger.Logger.Debug("GetedByLogin User info on level repo",
		zap.String("id", user.ID.String()),
		zap.String("login", user.Login))

	return &user, nil
}
