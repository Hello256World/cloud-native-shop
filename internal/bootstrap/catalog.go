package bootstrap

import (
	"github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/http"
	catalogpg "github.com/Hello256World/cloud-native-shop/internal/catalog/adapters/postgres"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
	"github.com/Hello256World/cloud-native-shop/internal/platform/middleware"
	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WireCatalog(api *gin.RouterGroup, pool *pgxpool.Pool, tv middleware.TokenVerifier) {
	uc := usecase.NewService(catalogpg.NewCategoryRepository(pool), postgres.NewTxManager(pool))

	http.RegisterRoutes(api, http.NewCategoryHandler(uc), middleware.Auth(tv), middleware.RequireRole("admin"))
}
