package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/models"
)
// toDo: add in ports
type AnalyticsRepository interface {
	InsertClick(ctx context.Context, e models.ClickEvent) error
	CountByShortCode(ctx context.Context, shortCode string) (uint64, error)
	AggregateDaily(ctx context.Context, shortCode string) ([]models.AggPoint, error)
	AggregateByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error)
}


type AnalyticsRepo struct {
	conn   driver.Conn 
}

func NewAnalyticsRepo(conn driver.Conn) *AnalyticsRepo {
	return &AnalyticsRepo{conn: conn}
}

func (r *AnalyticsRepo) InsertClick(ctx context.Context, e models.ClickEvent) error {
    query := `
        INSERT INTO click_analytics (short_code, client_id, user_agent, ip, timestamp)
        VALUES (?, ?, ?, ?, ?)
    `

	
	return r.conn.Exec(ctx, query,
		e.ShortCode,
		e.ClientID,
		e.UserAgent,
		e.IP,
		e.At, 
	)
}

func (r *AnalyticsRepo) CountByShortCode(ctx context.Context, shortCode string) (uint64, error) {
    query := `SELECT count() FROM click_analytics WHERE short_code = ?`

    var count uint64
    err := r.conn.QueryRow(ctx, query, shortCode).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (r *AnalyticsRepo) AggregateDaily(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
    query := `
        SELECT formatDateTime(timestamp, '%%Y-%%m-%%d') AS day, count() AS c
        FROM click_analytics
        WHERE short_code = ?
        GROUP BY day
        ORDER BY day ASC
    `

    rows, err := r.conn.Query(ctx, query, shortCode)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []models.AggPoint
    for rows.Next() {
        var p models.AggPoint
        if err := rows.Scan(&p.Key, &p.Count); err != nil {
            return nil, err
        }
        result = append(result, p)
    }
    return result, nil
}

func (r *AnalyticsRepo) AggregateByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
    query := `
        SELECT user_agent AS ua, count() AS c
        FROM click_analytics
        WHERE short_code = ?
        GROUP BY ua
        ORDER BY c DESC
    `

    rows, err := r.conn.Query(ctx, query, shortCode)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []models.AggPoint
    for rows.Next() {
        var p models.AggPoint
        if err := rows.Scan(&p.Key, &p.Count); err != nil {
            return nil, err
        }
        result = append(result, p)
    }
    return result, nil
}
