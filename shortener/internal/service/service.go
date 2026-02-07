package service

import (
	"context"

	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/models"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/repository"
	"github.com/google/uuid"
)

// Service объединяет операции над ссылками и их аналитикой
type ServiceInterface interface {
	// Ссылки
	CreateShort(ctx context.Context, sUrl *models.ShortURL) (*models.ShortURL, error)
	Resolve(ctx context.Context, shortCode string) (*models.ShortURL, error)

	// Аналитика
	RecordClick(ctx context.Context, ev models.ClickEvent) error
	Count(ctx context.Context, shortCode string) (uint64, error)
	Daily(ctx context.Context, shortCode string) ([]models.AggPoint, error)
	ByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error)
}

type service struct {
	urls         repository.Repository
	analytics    repository.AnalyticsRepository
	genShortCode func() (string, error)
}

func NewService(urls repository.Repository, analytics repository.AnalyticsRepository, genShortCode func() (string, error)) *service {
	return &service{urls: urls, analytics: analytics, genShortCode: genShortCode}
}

func (s *service) CreateShort(ctx context.Context, sUrl *models.ShortURL) (*models.ShortURL, error) {
	id := uuid.New()
	sUrl.ID = &id


	const maxAttempts = 3 //hardCode(sorry)
	var shortCode string
	var err error

	// var sUrl models.ShortURL


	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortCode, err = s.genShortCode()
		if err != nil {
			return nil, err
		}
		sUrl.ShortCode = shortCode
		err = s.urls.Create(ctx, sUrl)
		if err == nil {
			break
		}
		// возможна коллизия по unique(short_code) — попробуем другой код
		if attempt == maxAttempts-1 {
			return nil, err
		}
	}
	return sUrl, nil
}

func (s *service) Resolve(ctx context.Context, shortCode string) (*models.ShortURL, error) {
	return s.urls.GetByShortCode(ctx, shortCode)
}

func (s *service) RecordClick(ctx context.Context, ev models.ClickEvent) error {
	return s.analytics.InsertClick(ctx, ev)
}

func (s *service) Count(ctx context.Context, shortCode string) (uint64, error) {
	return s.analytics.CountByShortCode(ctx, shortCode)
}

func (s *service) Daily(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	return s.analytics.AggregateDaily(ctx, shortCode)
}

func (s *service) ByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	return s.analytics.AggregateByUserAgent(ctx, shortCode)
}
