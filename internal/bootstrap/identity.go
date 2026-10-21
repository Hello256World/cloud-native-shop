package bootstrap

import (
	"github.com/Hello256World/cloud-native-shop/internal/config"
	identityauth "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/auth"
	"github.com/Hello256World/cloud-native-shop/internal/identity/adapters/hasher"
	identityhttp "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/http"
	identitypg "github.com/Hello256World/cloud-native-shop/internal/identity/adapters/postgres"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
	platformauth "github.com/Hello256World/cloud-native-shop/internal/platform/auth"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WireIdentity(api *gin.RouterGroup, pool *pgxpool.Pool, tm *platformauth.TokenManager, cfg config.AuthConfig) {
	uc := usecase.NewService(
		identitypg.NewUserRepository(pool),
		identitypg.NewRefreshTokenRepository(pool),
		hasher.NewArgon2(),
		identityauth.NewAccessTokenIssuer(tm),
		postgres.NewTxManager(pool), cfg.RefreshTokenTTL,
	)

	identityhttp.RegisterRoutes(api, identityhttp.NewHandler(uc))
}
