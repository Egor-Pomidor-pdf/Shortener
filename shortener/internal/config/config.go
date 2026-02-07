package config

import (
	"fmt"
	"time"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/retry"
)

type Config struct {
	Env        string         `yaml:"env" env:"ENV"`
	Database   PostgresConfig 
	Server     ServerConfig   
	ClickHouse ClickHouseConfig
	// Redis           RedisConfig    

}



func NewConfig(envFilePath string, configFilePath string) (*Config, error) {
	myConfig := &Config{}

	cfg := config.New()

	if envFilePath != "" {
		if err := cfg.LoadEnvFiles(envFilePath); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}
	cfg.EnableEnv("")

	if configFilePath != "" {
		if err := cfg.LoadConfigFiles(configFilePath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	myConfig.Env = cfg.GetString("ENV")
	// Postgres
	myConfig.Database.MasterDSN = cfg.GetString("POSTGRES_MASTER_DSN")
	myConfig.Database.SlaveDSNs = cfg.GetStringSlice("POSTGRES_SLAVE_DSNS")
	myConfig.Database.MaxOpenConnections = cfg.GetInt("POSTGRES_MAX_OPEN_CONNECTIONS")
	myConfig.Database.MaxIdleConnections = cfg.GetInt("POSTGRES_MAX_IDLE_CONNECTIONS")
	myConfig.Database.ConnectionMaxLifetimeSeconds = cfg.GetInt("POSTGRES_CONNECTION_MAX_LIFETIME_SECONDS")
	// Postgres retry
	myConfig.Database.PostgresRetryConfig.Attempts = cfg.GetInt("RETRY_POSTGRES_ATTEMPTS")
	myConfig.Database.PostgresRetryConfig.DelayMilliseconds = cfg.GetInt("RETRY_POSTGRES_DELAY_MS")
	myConfig.Database.PostgresRetryConfig.Backoff = cfg.GetFloat64("RETRY_POSTGRES_BACKOFF")
	//server
	myConfig.Server.Host = cfg.GetString("SERVER_HOST")
	myConfig.Server.Port = cfg.GetInt("SERVER_PORT")

	// ClickHouse
	myConfig.ClickHouse.Addr = cfg.GetString("CLICKHOUSE_ADDR")
	myConfig.ClickHouse.Database = cfg.GetString("CLICKHOUSE_DATABASE")
	myConfig.ClickHouse.Username = cfg.GetString("CLICKHOUSE_USER")
	myConfig.ClickHouse.Password = cfg.GetString("CLICKHOUSE_PASSWORD")
	myConfig.ClickHouse.DialTimeout = cfg.GetInt("CLICKHOUSE_DIAL_TIMEOUT")
	myConfig.ClickHouse.MaxOpenConns = cfg.GetInt("CLICKHOUSE_MAX_OPEN_CONNS")
	myConfig.ClickHouse.MaxIdleConns = cfg.GetInt("CLICKHOUSE_MAX_IDLE_CONNS")
	// Redis
	// myConfig.Redis.Host = cfg.GetString("DELAYED_NOTIFIER_REDIS_HOST")
	// myConfig.Redis.Port = cfg.GetInt("DELAYED_NOTIFIER_REDIS_PORT")
	// myConfig.Redis.Password = cfg.GetString("DELAYED_NOTIFIER_REDIS_PASSWORD")
	// myConfig.Redis.DB = cfg.GetInt("DELAYED_NOTIFIER_REDIS_DB")
	// myConfig.Redis.Expiration = cfg.GetInt("DELAYED_NOTIFIER_REDIS_EXPIRATION")
	// Retry
	// // RedisRepository retry
	// myConfig.RedisRepoRetry.Attempts = cfg.GetInt("DELAYED_NOTIFIER_RETRY_REDIS_REPO_ATTEMPTS")
	// myConfig.RedisRepoRetry.DelayMilliseconds = cfg.GetInt("DELAYED_NOTIFIER_RETRY_REDIS_REPO_DELAY_MS")
	// myConfig.RedisRepoRetry.Backoff = cfg.GetFloat64("DELAYED_NOTIFIER_RETRY_REDIS_REPO_BACKOFF")

	return myConfig, nil
}

func MakeStrategy(c RetryConfig) retry.Strategy {
	return retry.Strategy{
		Attempts: c.Attempts,
		Delay:    time.Duration(c.DelayMilliseconds) * time.Millisecond,
		Backoff:  c.Backoff,
	}
}
