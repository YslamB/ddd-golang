package storage

import (
	"context"
	"fmt"
	"gddd/internal/infrastructure/config"
	db "gddd/pkg/database"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func PostgresInit(cfg *config.Config) *pgxpool.Pool {
	connectionString := buildConnectionString(cfg)
	pgxConfig, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		log.Fatalf("Unable to parse database cfg.Storage.Psql 💊: %v\n", err)
	}

	pgxConfig.MaxConns = cfg.Storage.Psql.MaxConnectionPoolSize
	pgxConfig.MaxConnLifetime = time.Duration(cfg.Storage.Psql.MaxConnectionLifetimeMinutes)
	_, err = db.NewDB(context.Background(), pgxConfig)

	if err != nil {
		log.Fatalf("❌ Could not create database connection pool: %v\n", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)

	if err != nil {
		log.Fatalf("failed to create connection poolpool🏊: %v\n", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("❌ Could not ping postgres🫙 database: %v", err))
	}

	log.Println("✅ Database connection pool initialized successfully")
	return pool
}

func buildConnectionString(cfg *config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Storage.Psql.Username, cfg.Storage.Psql.Password,
		cfg.Storage.Psql.Host, cfg.Storage.Psql.Port, cfg.Storage.Psql.Database,
	)
}
