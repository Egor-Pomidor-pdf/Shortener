package dto

import (
	"fmt"
	"net/url"

	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/models"
	"github.com/google/uuid"
)

type ShortenRequest struct {
	OriginalURL string     `json:"original_url"`
	ClientID    *uuid.UUID `json:"client_id"`
}

func (b ShortenRequest) ToEntity() (*models.ShortURL, error) {
	if b.OriginalURL == "" {
		return nil, fmt.Errorf("original_url is required")
	}

	parsed, err := url.ParseRequestURI(b.OriginalURL)
	if err != nil {
		return nil, fmt.Errorf("invalid original_url: %w", err)
	}

	return &models.ShortURL{
		Original: parsed.String(),
		ClientID: b.ClientID,
	}, nil
}
