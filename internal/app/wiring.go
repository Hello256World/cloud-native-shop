package app

import "github.com/Hello256World/cloud-native-shop/internal/bootstrap"

func (a *application) mapHandlers() {
	bootstrap.WireHealth(a.httpEngine, a.pgxPool, a.redisClient)

	if !a.conf.IsProduction() {
		bootstrap.WireSwagger(a.httpEngine)
	}
}
