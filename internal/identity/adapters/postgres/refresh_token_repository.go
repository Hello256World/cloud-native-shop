package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/ports"
	pg "github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type refreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) ports.RefreshTokenRepository {
	return &refreshTokenRepository{
		pool: pool,
	}
}

func (r *refreshTokenRepository) q(c context.Context) *db.Queries {
	return db.New(pg.Executor(c, r.pool))
}

func (r *refreshTokenRepository) Create(c context.Context, userID int64, tokenHash string, expiredAt time.Time) error {
	q := r.q(c)

	err := q.InsertRefreshToken(c, db.InsertRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiredAt,
	})

	if err != nil {
		return fmt.Errorf("identity.repo.RefreshToken.Create user=%d: %w", userID, err)
	}

	return nil
}

func (r *refreshTokenRepository) GetByHash(c context.Context, tokenHash string) (*domain.RefreshToken, error) {
	q := r.q(c)

	row, err := q.GetRefreshTokenByHash(c, tokenHash)

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrRefreshTokenInvalid
		}

		return nil, fmt.Errorf("identity.repo.RefreshToken.GetByHash: %w", err)
	}

	return &domain.RefreshToken{
		ID: row.ID, UserID: row.UserID, TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt, RevokedAt: row.RevokedAt, CreatedAt: row.CreatedAt,
	}, nil
}

func (r *refreshTokenRepository) Revoke(c context.Context, id int64) error {
	if err := r.q(c).RevokeRefreshToken(c, id); err != nil {
		return fmt.Errorf("identity.repo.RefreshToken.Revoke id=%d: %w", id, err)
	}
	return nil
}
