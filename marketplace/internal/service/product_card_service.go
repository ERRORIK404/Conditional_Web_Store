package service

import (
	"errors"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/repository"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/transport_rest/middleware"
	"go.uber.org/zap"

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
	if authStatus.IsAuthenticated == false {
		logger.Logger.Error("product_card_service.CreateProductCard unauthorized on level service")
		return nil, errors.New("unauthorized")
	}

	// Валидация данных
	if len(title) < 3 || len(title) > 99 {
		logger.Logger.Error("title must be between 3-100 characters on level service")
		return nil, errors.New("title must be between 3-100 characters")
	}
	if len(description) > 999 {
		logger.Logger.Error("description is too long on level service")
		return nil, errors.New("description too long")
	}
	if price <= 0 {
		logger.Logger.Error("price must be greater than zero on level service")
		return nil, errors.New("price must be positive")
	}

	productCard := models.ProductCard{
		ID:          uuid.Nil,
		UserID:      authStatus.UserID,
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
	}

	id, err := s.productCard.Create(productCard)
	if err != nil {
		return nil, err
	}
	productCard.ID = id

	logger.Logger.Debug("productCard created on lever service", zap.Any("productCard", productCard))
	return &productCard, nil
}
