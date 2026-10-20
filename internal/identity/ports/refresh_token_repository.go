package ports

import (
	"context"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id int64) error
}
