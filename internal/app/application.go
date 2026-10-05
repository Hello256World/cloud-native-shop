package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application interface {
	Run(context.Context) error
}

type application struct {
	conf       *config.Config
	httpEngine *gin.Engine
	pgxPool    *pgxpool.Pool
}

func NewApplication(config *config.Config, pool *pgxpool.Pool) Application {
	app := &application{
		conf:       config,
		httpEngine: gin.Default(),
		pgxPool:    pool,
	}

	app.mapHnadler()

	return app
}

func (a *application) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.conf.HTTPConfig.Port),
		Handler: a.httpEngine,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed to start: %w", err)
	case <-ctx.Done():
		log.Printf("received signal %s, shutting down", context.Cause(ctx).Error())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	a.pgxPool.Close()

	return nil
}
