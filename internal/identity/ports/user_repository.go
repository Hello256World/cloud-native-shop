package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) (*domain.User, error)
	GetByID(context.Context, int64) (*domain.User, error)
	GetByPhone(context.Context, string) (*domain.User, error)
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) (bool, error)
}
