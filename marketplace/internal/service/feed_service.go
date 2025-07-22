package service

import (
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/models"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FeedService struct {
	feedRepo repository.ProductCardRepository
}

func NewFeedService(feedRepo repository.ProductCardRepository) *FeedService {
	return &FeedService{feedRepo: feedRepo}
}

// GetFeed возвращает ленту объявлений с пагинацией и фильтрами
func (s *FeedService) GetFeed(
	page, limit int,
	sortField, sortOrder string,
	minPrice, maxPrice float64,
	currentUserID uuid.UUID,
) ([]models.ProductCard, int, error) {

	// Валидация параметров
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 20 {
		limit = 5
	}

	feed, err := s.feedRepo.ListProductCard(
		currentUserID,
		&minPrice,
		&maxPrice,
		sortField,
		sortOrder,
		page,
		limit,
	)
	if err != nil {
		logger.Logger.Error("feed_service.GetFeed", zap.Error(err))
		return nil, 0, err
	}
	// Получение данных из репозитория
	return feed, len(feed), nil
}
