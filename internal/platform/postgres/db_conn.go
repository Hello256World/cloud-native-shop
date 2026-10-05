package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	URI         string
	AppName     string
	MinPoolSize int32
	MaxPoolSize int32
}

func Connect(ctx context.Context, opt Options) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(opt.URI)

	if err != nil {
		return nil, fmt.Errorf("error parsing config for postgres pool: %v", err)
	}

	poolConfig.MinConns = opt.MinPoolSize
	poolConfig.MaxConns = opt.MaxPoolSize

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, fmt.Errorf("error creating postgres pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error ping postgres pool: %v", err)
	}

	return pool, nil
}
