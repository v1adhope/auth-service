package postgresql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

type Config struct {
	ConnStr string
}

func Build(ctx context.Context, cfg Config) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("postgresql: postgresql: Build: New: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgresql: postgresql: Build: Ping: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{pool, builder}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}
