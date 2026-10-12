package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/config"
	"github.com/Hello256World/cloud-native-shop/internal/platform/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Application interface {
	Run(context.Context) error
}

type application struct {
	conf        *config.Config
	log         *slog.Logger
	httpEngine  *gin.Engine
	pgxPool     *pgxpool.Pool
	redisClient *redis.Client
}

func NewApplication(config *config.Config, pool *pgxpool.Pool, redis *redis.Client, log *slog.Logger) Application {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(),
	)

	app := &application{
		conf:        config,
		log:         log,
		httpEngine:  engine,
		pgxPool:     pool,
		redisClient: redis,
	}

	app.mapHandlers()
	return app
}

func (a *application) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.conf.HTTPConfig.Port),
		Handler: a.httpEngine,
	}

	errChan := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("server listening on %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed to start: %w", err)
	case <-ctx.Done():
		slog.Info("received signal shutting down", "signal", context.Cause(ctx))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	a.pgxPool.Close()

	if err := a.redisClient.Close(); err != nil {
		return fmt.Errorf("error closing redis: %w", err)
	}

	return nil
}
