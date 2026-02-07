package clickhouseConn

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/config"
)

func InitClickHouseConn(cfg *config.ClickHouseConfig) (driver.Conn, error) {
    // Опции для подключения
    opts := &clickhouse.Options{
        Addr: []string{"localhost:9000"},
        Auth: clickhouse.Auth{
            Database: cfg.Database,
            Username: cfg.Username,
            Password: cfg.Password,
        },
        DialTimeout:      time.Duration(cfg.DialTimeout) * time.Second,
        MaxOpenConns:     int(cfg.MaxOpenConns),
        MaxIdleConns:     int(cfg.MaxIdleConns),
    }

    // Создаём соединение
    conn, err := clickhouse.Open(opts)
    if err != nil {
        return nil, fmt.Errorf("clickhouse open: %w", err)
    }

    // Проверяем соединение
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := conn.Ping(ctx); err != nil {
        conn.Close()
        return nil, fmt.Errorf("clickhouse ping: %w", err)
    }

    return conn, nil
}
