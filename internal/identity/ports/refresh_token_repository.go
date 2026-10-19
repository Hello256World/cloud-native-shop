package ports

import (
	"context"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type RefreshTokenRepository interface {
	Create(context.Context, int64, string, time.Time) error
	GetByHash(context.Context, string) (*domain.RefreshToken, error)
	Revoke(context.Context, int64) error
}
