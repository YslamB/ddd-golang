package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ DB = (*Database)(nil)

type Querier interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type DB interface {
	Querier
	Get(ctx context.Context, dest any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	Close()
}

type Database struct {
	pool *pgxpool.Pool
}

func NewDB(ctx *context.Context, cfg *pgxpool.Config) (*Database, error) {
	pool, err := pgxpool.NewWithConfig(*ctx, cfg)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(*ctx); err != nil {
		return nil, err
	}

	return &Database{pool: pool}, nil
}

func (d *Database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, query, args...)
}

func (d *Database) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args...)
}

func (d *Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	result, err := d.pool.Exec(ctx, query, args...)

	if err != nil {
		return pgconn.CommandTag{}, err
	}

	if result.RowsAffected() == 0 {
		return pgconn.CommandTag{}, pgx.ErrNoRows
	}

	return result, nil
}

func (d *Database) Get(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, d.pool, dest, query, args...)
}

func (d *Database) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, d.pool, dest, query, args...)
}

func (d *Database) Close() {
	d.pool.Close()
}
