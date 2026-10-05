package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	URI      string
	Password string
	Database int
}

func NewClient(cfg Options) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.URI,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error ping redis: %w", err)
	}

	return client, nil
}
