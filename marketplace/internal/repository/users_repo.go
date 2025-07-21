package repository

import (
	"errors"
	"github.com/google/uuid"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (uuid.UUID, error)
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
func (r *userRepository) Create(user *models.User) (uuid.UUID, error) {
	if err := r.db.Create(user).Error; err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

// GetByID возвращает пользователя по ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByLogin возвращает пользователя по логину
func (r *userRepository) GetByLogin(login string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
