package ports

import (
	"context"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type AccessTokenIssuer interface {
	Issue(userID int64, roles []domain.Role) (token string, expiresAt time.Time, err error)
}

type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}
