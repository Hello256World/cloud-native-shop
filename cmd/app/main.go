package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/app"
	"github.com/Hello256World/cloud-native-shop/internal/config"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conf, err := config.GetConfig()

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postgresPool, err := postgres.Connect(ctx, postgres.Options{
		URI:         conf.PostgresConfig.GetConnectionURI(),
		AppName:     conf.PostgresConfig.AppName,
		MinPoolSize: int32(conf.PostgresConfig.MinPoolSize),
		MaxPoolSize: int32(conf.PostgresConfig.MaxPoolSize),
	})

	if err != nil {
		return err
	}

	app := app.NewApplication(conf, postgresPool)

	runCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(runCtx); err != nil {
		return err
	}

	return nil
}
