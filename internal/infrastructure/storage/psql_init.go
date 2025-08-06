package storage

import (
	"context"
	"fmt"
	infra_config "gddd/internal/infrastructure/config"
	infra_logging "gddd/internal/infrastructure/logging"
	db "gddd/pkg/database"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func PostgresInit(ctx *context.Context, cfg *infra_config.Config, log *infra_logging.Logger) *db.Database {
	connectionString := buildConnectionString(cfg)
	pgxConfig, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		log.Fatalf("Unable to parse database cfg.Storage.Psql 💊: %v\n", err)
	}

	pgxConfig.MaxConns = cfg.Storage.Psql.MaxConnectionPoolSize
	pgxConfig.MaxConnLifetime = time.Duration(cfg.Storage.Psql.MaxConnectionLifetimeMinutes)
	pool, err := db.NewDB(ctx, pgxConfig)

	if err != nil {
		log.Fatalf("❌ not connect to db: %v\n", err)
	}

	return pool
}

func buildConnectionString(cfg *infra_config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Storage.Psql.Username, cfg.Storage.Psql.Password,
		cfg.Storage.Psql.Host, cfg.Storage.Psql.Port, cfg.Storage.Psql.Database,
	)
}
