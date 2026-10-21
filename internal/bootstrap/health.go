package bootstrap

import (
	"context"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/platform/health"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func WireHealth(engine *gin.Engine, pool *pgxpool.Pool, rds *redis.Client) {
	healthHandler := health.NewHandler(2*time.Second,
		health.NewChecker("postgres", func(ctx context.Context) error {
			return pool.Ping(ctx)
		}),
		health.NewChecker("redis", func(ctx context.Context) error {
			return rds.Ping(ctx).Err()
		}),
	)
	healthHandler.Register(engine)
}
