package models

import (
	"time"

	"github.com/google/uuid"
)

// ShortURL доменная модель хранения короткой ссылки
type ShortURL struct {
	ID        *uuid.UUID
	Original  string
	ShortCode string
	ClientID  *uuid.UUID
	CreatedAt time.Time
}

// ClickEvent событие перехода (ClickHouse)
type ClickEvent struct {
	ShortCode string
	ClientID  uuid.UUID
	UserAgent string
	IP        string
	At        time.Time
}

// AggPoint агрегированная точка для графиков/сводок
type AggPoint struct {
	Key   string
	Count int64
}
