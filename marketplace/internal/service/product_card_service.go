package service

import (
	"errors"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/repository"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/transport_rest/middleware"

	"github.com/google/uuid"
)

type ProductCardService struct {
	productCard repository.ProductCardRepository
	userRepo    repository.UserRepository
}

func NewProductCardService(productCard repository.ProductCardRepository, userRepo repository.UserRepository) *ProductCardService {
	logger.Logger.Debug("created new ProductCardService")
	return &ProductCardService{
		productCard: productCard,
		userRepo:    userRepo,
	}
}

// CreateProductCard создает новое объявление с проверкой существования пользователя
func (s *ProductCardService) CreateProductCard(authStatus middleware.AuthStatus, title, description, imageURL string, price float64) (*models.ProductCard, error) {
	// Проверка существования пользователя( нужна была бы если бы не было мидлвеер для проверки авторизации)
	//if _, err := s.userRepo.GetByID(userID); err != nil {
	//	logger.Logger.Error("user not found")
	//	return nil, errors.New("user not found")
	//}

	// Валидация данных
	if len(title) < 3 || len(title) > 99 {
		logger.Logger.Error("title must be between 3-100 characters")
		return nil, errors.New("title must be between 3-100 characters")
	}
	if len(description) > 999 {
		logger.Logger.Error("description is too long")
		return nil, errors.New("description too long")
	}
	if price <= 0 {
		logger.Logger.Error("price must be greater than zero")
		return nil, errors.New("price must be positive")
	}

	productCard := &models.ProductCard{
		ID:          uuid.New(),
		UserID:      authStatus.UserID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
	}
	id
	if err := s.productCard.Create(productCard); err != nil {
		return nil, err
	}

	return productCard, nil
}
