package repository

import (
	"errors"
	"fmt"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type ProductCardRepository interface {
	Create(productCard models.ProductCard) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*models.ProductCard, error)
	ListProductCard(
		currentUserID uuid.UUID,
		minPrice, maxPrice *float64,
		sortField, sortOrder string,
		page, pageSize int,
	) ([]models.ProductCard, error)
}

type productCardRepository struct {
	db *gorm.DB
}

func NewProductCardRepository(db *gorm.DB) ProductCardRepository {
	return &productCardRepository{db: db}
}

func (r *productCardRepository) Create(productCard models.ProductCard) (uuid.UUID, error) {
	err := r.db.Create(productCard).Error
	if err != nil {
		logger.Logger.Error("productCardRepository.Create", zap.Error(err))
		return uuid.Nil, err
	}

	logger.Logger.Debug("Product Card Created", zap.Any("productCard", productCard))
	return productCard.ID, nil
}

func (r *productCardRepository) GetByID(id uuid.UUID) (*models.ProductCard, error) {
	var pc models.ProductCard

	if err := r.db.First(&pc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("productCardRepository.GetByID(NOT FOUND)", zap.Error(err))
			return nil, err
		}
		logger.Logger.Error("productCardRepository.GetByID", zap.Error(err))
		return nil, err
	}

	logger.Logger.Debug("Product Card Retrieved", zap.Any("pc", pc))
	return &pc, nil
}

func (r *productCardRepository) ListProductCard(
	currentUserID uuid.UUID,
	minPrice, maxPrice *float64,
	sortField, sortOrder string,
	page, pageSize int,
) ([]models.ProductCard, error) {
	// Валидация параметров сортировки
	sortField, sortOrder = validateSortParams(sortField, sortOrder)

	// Рассчитываем смещение для пагинации
	offset := (page - 1) * pageSize

	// Строим запрос
	query := r.db.Table("product_card").
		Select(
			"product_card.title",
			"product_card.description",
			"product_card.image_url",
			"product_card.price",
			"users.login AS author_login",
		).
		Joins("JOIN users ON product_card.user_id = users.id")

	// Добавляем проверку принадлежности для авторизованных пользователей
	if currentUserID != uuid.Nil {
		query = query.Select(
			query.Statement.Selects,
			fmt.Sprintf("(product_card.user_id = '%s') AS is_own", currentUserID.String()),
		)
	} else {
		query = query.Select(
			query.Statement.Selects,
			"false AS is_own",
		)
	}

	// Фильтрация по цене
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}

	// Применяем сортировку
	switch sortField {
	case "price":
		query = query.Order(fmt.Sprintf("price %s", strings.ToUpper(sortOrder)))
	case "created_at":
		query = query.Order(fmt.Sprintf("created_at %s", strings.ToUpper(sortOrder)))
	default:
		query = query.Order("created_at DESC")
	}

	// Применяем пагинацию
	query = query.Offset(offset).Limit(pageSize)

	// Выполняем запрос
	var results []models.ProductCard
	if err := query.Scan(&results).Error; err != nil {
		logger.Logger.Error("productCardRepository.ListProductCard (at the place where the request was executed)", zap.Error(err))
		return nil, fmt.Errorf("failed to list product cards: %w", err)
	}

	logger.Logger.Debug("Product Card List", zap.Any("results", results))
	return results, nil
}

// Валидация параметров сортировки
func validateSortParams(field, order string) (string, string) {
	// Проверяем допустимые поля сортировки
	switch strings.ToLower(field) {
	case "price", "created_at":
	default:
		field = "created_at"
	}

	// Проверяем допустимые направления
	switch strings.ToLower(order) {
	case "asc", "desc":
	default:
		order = "desc"
	}

	return field, order
}
