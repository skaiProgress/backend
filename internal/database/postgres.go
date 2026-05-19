package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres wraps a pgx connection pool.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres creates a connection pool to PostgreSQL.
func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	return &Postgres{pool: pool}, nil
}

// Ping verifies the database connection is alive.
func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// Pool exposes the underlying connection pool for repositories.
func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}

// Close releases all connections in the pool.
func (p *Postgres) Close() {
	p.pool.Close()
}
