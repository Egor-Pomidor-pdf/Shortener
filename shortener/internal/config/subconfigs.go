package config

type PostgresConfig struct {
	MasterDSN                    string   `env:"MASTER_DSN"`
	SlaveDSNs                    []string `env:"SLAVE_DSNS" envSeparator:","`
	MaxOpenConnections           int      `env:"MAX_OPEN_CONNECTIONS" envDefault:"3"`
	MaxIdleConnections           int      `env:"MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnectionMaxLifetimeSeconds int      `env:"CONNECTION_MAX_LIFETIME_SECONDS" envDefault:"0"`
	PostgresRetryConfig RetryConfig 
}

type ClickHouseConfig struct {
	Addr         string `env:"CLICKHOUSE_ADDR"`     
	Database     string `env:"CLICKHOUSE_DATABASE"`
	Username     string `env:"CLICKHOUSE_USER"`
	Password     string `env:"CLICKHOUSE_PASSWORD"`
	DialTimeout  int    `env:"CLICKHOUSE_DIAL_TIMEOUT"`
	MaxOpenConns int    `env:"CLICKHOUSE_MAX_OPEN_CONNS"`
	MaxIdleConns int    `env:"CLICKHOUSE_MAX_IDLE_CONNS"`
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST"` 
	Port int    `yaml:"SERVER_PORT"` 
}

type RetryConfig struct {
	Attempts          int     `env:"ATTEMPTS"`
	DelayMilliseconds int     `env:"DELAY_MS"`
	Backoff           float64 `env:"BACKOFF"`
}

// type RedisConfig struct {
// 	Host       string `yaml:"host" env:"HOST"`             // Адрес Redis (например, "localhost")
// 	Port       int    `yaml:"port" env:"PORT"`             // Порт Redis (обычно 6379)
// 	Password   string `yaml:"password" env:"PASSWORD"`     // Пароль, если настроена аутентификация
// 	DB         int    `yaml:"db" env:"DB"`                 // Номер базы Redis (по умолчанию 0)
// 	Expiration int    `yaml:"expiration" env:"EXPIRATION"` // Время жизни ключей (TTL)
// }
