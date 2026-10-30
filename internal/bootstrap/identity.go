package bootstrap

import (
	"github.com/Hello256World/cloud-native-shop/internal/config"
	identityauth "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/auth"
	"github.com/Hello256World/cloud-native-shop/internal/identity/adapters/hasher"
	identityhttp "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/http"
	identitypg "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/postgres"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/internal/platform/middleware"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newIdentity(pool *pgxpool.Pool, conf config.AuthConfig) (*usecase.Service, *auth.TokenManager) {
	tm := auth.NewTokenManager(conf.JWTSecret, conf.Issuer, conf.AccessTokenTTL)

	uc := usecase.NewService(
		identitypg.NewUserRepository(pool),
		identitypg.NewRefreshTokenRepository(pool),
		identitypg.NewAddressRepository(pool),
		hasher.NewArgon2(),
		identityauth.NewAccessTokenIssuer(tm),
		postgres.NewTxManager(pool),
		conf.RefreshTokenTTL,
	)

	return uc, tm
}

func NewIdentityService(pool *pgxpool.Pool, conf config.AuthConfig) *usecase.Service {
	uc, _ := newIdentity(pool, conf)
	return uc
}

func WireIdentity(api *gin.RouterGroup, pool *pgxpool.Pool, conf config.AuthConfig) {
	uc, tm := newIdentity(pool, conf)

	identityhttp.RegisterRoutes(api, identityhttp.NewAuthHandler(uc), identityhttp.NewAddressHandler(uc), middleware.Auth(tm))
}
