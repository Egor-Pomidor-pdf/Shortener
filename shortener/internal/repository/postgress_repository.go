package repository

import (
	"context"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/models"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

type StoreRepository struct {
	db       *dbpg.DB
	strategy retry.Strategy
}

func NewRepository(db *dbpg.DB, strategy retry.Strategy) *StoreRepository {
	return &StoreRepository{
		db:       db,
		strategy: strategy,
	}
}

type Repository interface {
	Create(ctx context.Context, s *models.ShortURL) error
	GetByShortCode(ctx context.Context, shortCode string) (*models.ShortURL, error)
	ExistsByOriginalURL(ctx context.Context, original string) (bool, error)
}

func (r *StoreRepository) Create(ctx context.Context, s *models.ShortURL) error {
	query := `
        INSERT INTO shortcuts (id, original_url, short_code, client_id, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecWithRetry(
		ctx,
		r.strategy,
		query,
		s.ID,
		s.Original,
		s.ShortCode,
		s.ClientID,
		s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *StoreRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.ShortURL, error) {
	query := `
		SELECT id, original_url, short_code, client_id, created_at
		FROM shortcuts
		WHERE short_code = $1
	`

	row, err := r.db.QueryRowWithRetry(
		ctx,
		r.strategy,
		query,
		shortCode,
	)
	if err != nil {
		return nil, err
	}

	var s models.ShortURL

	err = row.Scan(
		&s.ID,
		&s.Original,
		&s.ShortCode,
		&s.ClientID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StoreRepository) ExistsByOriginalURL(ctx context.Context, original string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shortcuts WHERE original_url = $1)`

	row, err := r.db.QueryRowWithRetry(
		ctx,
		r.strategy,
		query,
		original,
	)
	if err != nil {
		return false, err
	}

	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
