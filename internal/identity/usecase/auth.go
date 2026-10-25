package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/pkg/logger"
)

type RegisterInput struct {
	FullName string
	Phone    string
	Email    *string
	Password string
}

type TokenPair struct {
	AccessToken      string
	ExpiresAt        time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

func (s *Service) Register(c context.Context, in RegisterInput) (*domain.User, error) {
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("identity.Register hash: %w", err)
	}

	u := &domain.User{
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		PasswordHash: &hash,
		IsActive:     true,
		Roles:        []domain.Role{domain.RoleCustomer},
	}

	var created *domain.User

	err = s.tx.WithinTx(c, func(ctx context.Context) error {
		var err error
		created, err = s.users.Create(ctx, u)
		return err
	})
	if err != nil {
		return nil, err
	}

	logger.FromContext(c).Info("user registered", slog.Int64("user_id", created.ID))
	return created, nil
}

func (s *Service) Login(c context.Context, phone, password string) (*TokenPair, error) {
	user, err := s.users.GetByPhone(c, phone)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash == nil {
		return nil, domain.ErrInvalidCredentials
	}

	ok, err := s.hasher.Verify(password, *user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("identity.Login verify: %w", err)
	}

	if !ok {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	return s.issuePair(c, user, s.now().Add(s.refreshTTL))
}

func (s *Service) Refresh(c context.Context, raw string) (*TokenPair, error) {
	var pair *TokenPair

	err := s.tx.WithinTx(c, func(ctx context.Context) error {
		rt, err := s.refresh.GetByHash(ctx, hashToken(raw))

		if err != nil {
			return err
		}

		if !rt.IsUsable(s.now()) {
			return domain.ErrRefreshTokenInvalid
		}

		u, err := s.users.GetByID(ctx, rt.UserID)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				return domain.ErrRefreshTokenInvalid
			}
			return err
		}

		if !u.IsActive {
			return domain.ErrUserInactive
		}

		if err := s.refresh.Revoke(ctx, rt.ID); err != nil {
			return err
		}

		pair, err = s.issuePair(ctx, u, rt.ExpiresAt)

		return err
	})

	if err != nil {
		return nil, err
	}

	return pair, nil
}

func (s *Service) Logout(c context.Context, raw string) error {
	rt, err := s.refresh.GetByHash(c, hashToken(raw))

	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenInvalid) {
			return nil
		}
		return err
	}

	return s.refresh.Revoke(c, rt.ID)
}

func (s *Service) Me(c context.Context, userID int64) (*domain.User, error) {
	return s.users.GetByID(c, userID)
}

func (s *Service) issuePair(c context.Context, user *domain.User, refreshExpiresAt time.Time) (*TokenPair, error) {
	access, expireAt, err := s.tokens.Issue(user.ID, user.Roles)

	if err != nil {
		return nil, fmt.Errorf("identity.issuePair access: %w", err)
	}

	raw, hash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.refresh.Create(c, user.ID, hash, refreshExpiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      access,
		ExpiresAt:        expireAt,
		RefreshToken:     raw,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

