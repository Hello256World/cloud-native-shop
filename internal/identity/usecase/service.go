package usecase

import (
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/ports"
)

type Service struct {
	users      ports.UserRepository
	refresh    ports.RefreshTokenRepository
	address    ports.AddressRepository
	hasher     ports.PasswordHasher
	tokens     ports.AccessTokenIssuer
	tx         ports.TxManager
	refreshTTL time.Duration
	now        func() time.Time
}

func NewService(users ports.UserRepository, refresh ports.RefreshTokenRepository, address ports.AddressRepository, hasher ports.PasswordHasher, tokens ports.AccessTokenIssuer, tx ports.TxManager, refreshTTL time.Duration) *Service {
	return &Service{
		users:      users,
		refresh:    refresh,
		address:    address,
		hasher:     hasher,
		tokens:     tokens,
		tx:         tx,
		refreshTTL: refreshTTL,
		now:        time.Now,
	}
}
