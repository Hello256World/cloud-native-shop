package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/bootstrap"
	"github.com/Hello256World/cloud-native-shop/internal/config"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/Hello256World/cloud-native-shop/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("seed failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	slog.SetDefault(logger.New(conf.Environment, conf.LogConfig.Level))

	phone := os.Getenv("SUPERADMIN_PHONE")
	password := os.Getenv("SUPERADMIN_PASSWORD")
	if phone == "" || len(password) < 12 {
		return fmt.Errorf("SUPERADMIN_PHONE is required and SUPERADMIN_PASSWORD must be >= 12 chars")
	}
	name := os.Getenv("SUPERADMIN_NAME")
	if name == "" {
		name = "Super Admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, postgres.Options{
		URI:         conf.PostgresConfig.GetConnectionURI(),
		AppName:     conf.PostgresConfig.AppName + "-seed",
		MinPoolSize: 1,
		MaxPoolSize: 2,
	})
	if err != nil {
		return err
	}

	defer pool.Close()

	uc := bootstrap.NewIdentityService(pool, conf.AuthConfig)

	created, err := uc.EnsureSuperAdmin(ctx, usecase.RegisterInput{
		FullName: name,
		Phone:    phone,
		Password: password,
	})

	if created {
		slog.Info("super admin created", "phone", phone)
	} else {
		slog.Info("super admin already exists, nothing to do", "phone", phone)
	}
	return nil
}
