package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

func (s *Service) EnsureSuperAdmin(c context.Context, in RegisterInput) (created bool, err error) {
	if _, err := s.users.GetByPhone(c, in.Phone); err == nil {
		return false, err
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return false, err
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return false, fmt.Errorf("identity.EnsureAdmin hash: %w", err)
	}

	u := &domain.User{
		FullName:     in.FullName,
		Email:        in.Email,
		Phone:        in.Phone,
		PasswordHash: &hash,
		IsActive:     true,
		Roles:        []domain.Role{domain.RoleAdmin, domain.RoleSuperAdmin},
	}

	err = s.tx.WithinTx(c, func(ctx context.Context) error {
		_, err := s.users.Create(ctx, u)
		return err
	})

	if errors.Is(err, domain.ErrPhoneTaken) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
