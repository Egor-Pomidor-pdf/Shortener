package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	clickhouseConn "github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/clickHouse"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/config"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/handler"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/repository"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/service"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/internal/utils"
	postgres "github.com/Egor-Pomidor-pdf/Shortener/shortener/pkg/db"
	"github.com/Egor-Pomidor-pdf/Shortener/shortener/pkg/server"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

func main() {

	// make context
	ctx := context.Background()
	ctx, ctxStop := signal.NotifyContext(ctx, os.Interrupt)

	// init config
	cfg, err := config.NewConfig("../config/.env", "")
	if err != nil {
		log.Fatal(err)
	}

	// init logger
	zlog.InitConsole()
	err = zlog.SetLevel(cfg.Env)
	if err != nil {
		log.Fatal(fmt.Errorf("error setting log level to '%s': %w", cfg.Env, err))
	}
	zlog.Logger.Info().
		Str("env", cfg.Env).
		Msg("Start app...")

	// strategies
	postgresRetryStrategy := config.MakeStrategy(cfg.Database.PostgresRetryConfig)

	// connect to db
	var postgresDB *dbpg.DB
	err = retry.DoContext(ctx, postgresRetryStrategy, func() error {
		var postgresConnErr error
		postgresDB, postgresConnErr = dbpg.New(cfg.Database.MasterDSN, cfg.Database.SlaveDSNs,
			&dbpg.Options{
				MaxOpenConns:    cfg.Database.MaxOpenConnections,
				MaxIdleConns:    cfg.Database.MaxIdleConnections,
				ConnMaxLifetime: time.Duration(cfg.Database.ConnectionMaxLifetimeSeconds) * time.Second,
			})
		return postgresConnErr
	})

	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	zlog.Logger.Info().Msg("Successfully connected to PostgreSQL")

	migrationsPathClickHouse := "file://./internal/migrations/clickhouse" //"../internal/migrations/clickhouse"
	err = postgres.MigrateUpClickHouse(cfg.ClickHouse.Addr, migrationsPathClickHouse)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("couldn't migrate migrations fo ClickHouse")
	}

	zlog.Logger.Info().Msg("Successfully connected to create migrations fo ClickHouse")

	// connect to ClickHouse
	conn, err := clickhouseConn.InitClickHouseConn(&cfg.ClickHouse)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("failed to connect to cHouse")
	}

	zlog.Logger.Info().Msg("Successfully connected to ClickHouse")

	// create migrations
	migrationsPathPostgress := "file://./internal/migrations/db" // "../internal/migrations/db"
	err = postgres.MigrateUp(cfg.Database.MasterDSN, migrationsPathPostgress)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("couldn't migrate postgres on master DSN")
	}

	zlog.Logger.Info().Msg("Successfully connected to create migrations fo PSQL")


	// init repo
	store := repository.NewRepository(postgresDB, postgresRetryStrategy)
	anal := repository.NewAnalyticsRepo(conn)

	// inint crud service
	srv := service.NewService(store, anal, func() (string, error) { return utils.GenerateShortCode(8) })
	handl := handler.NewHandler(srv)
	router := handler.NewRouter(handl)

	// running server
	zlog.Logger.Info().Msg("server start")
	httpServer := server.NewHTTPServer(router)
	err = httpServer.GracefulRun(ctx, cfg.Server.Host, cfg.Server.Port)

	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Msg("failed GracefulRun server")
	}

	zlog.Logger.Info().Msg("server gracefully stopped")
	ctxStop()
	zlog.Logger.Info().Msg("background operations gracefully stopped")
}
