package app

import (
	"github.com/Hello256World/cloud-native-shop/internal/bootstrap"
	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
)

func (a *application) mapHandlers() {
	bootstrap.WireHealth(a.httpEngine, a.pgxPool, a.redisClient)

	if !a.conf.IsProduction() {
		bootstrap.WireSwagger(a.httpEngine)
	}

	tm := auth.NewTokenManager(a.conf.AuthConfig.JWTSecret, a.conf.AuthConfig.Issuer, a.conf.AuthConfig.AccessTokenTTL)

	// version one
	api := a.httpEngine.Group("/api/v1")
	bootstrap.WireIdentity(api, a.pgxPool, tm, a.conf.AuthConfig)
}
