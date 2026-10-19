package postgres

import (
	"context"
	"fmt"

	"github.com/Hello256World/cloud-native-shop/internal/identity/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/ports"
	pg "github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) q(ctx context.Context) *db.Queries {
	return db.New(pg.Executor(ctx, r.pool))
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	q := r.q(ctx)

	row, err := q.InsertUser(ctx, db.InsertUserParams{
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone,
		PasswordHash: user.PasswordHash,
		Birthday:     user.Birthday,
		Gender:       toDBGender(user.Gender),
		IsActive:     user.IsActive,
	})

	if err != nil {
		switch {
		case constraintContains(err, "phone"):
			return nil, domain.ErrPhoneTaken
		case constraintContains(err, "email"):
			return nil, domain.ErrEmailTaken
		}

		return nil, fmt.Errorf("identity.repo.User.Create: %w", err)
	}

	for _, role := range user.Roles {
		if err := q.AddUserRoleByName(ctx, db.AddUserRoleByNameParams{
			UserID: row.ID,
			Name:   string(role),
		}); err != nil {
			return nil, fmt.Errorf("identity.repo.User.Create role=%s: %w", role, err)
		}
	}

	return r.toDomain(row, user.Roles), nil
}

func (r *userRepository) GetByID(c context.Context, id int64) (*domain.User, error) {
	q := r.q(c)

	row, err := q.GetUserByID(c, id)
	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("identity.repo.User.GetByID id=%d: %w", id, err)
	}
	return r.withRoles(c, q, row)
}

func (r *userRepository) GetByPhone(c context.Context, phone string) (*domain.User, error) {
	q := r.q(c)

	row, err := q.GetUserByPhone(c, phone)

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("identity.repo.User.GetByPhone: %w", err)
	}
	return r.withRoles(c, q, row)
}

func (r *userRepository) withRoles(c context.Context, q *db.Queries, row db.IdentityUser) (*domain.User, error) {
	name, err := q.ListRoleNamesByUser(c, row.ID)

	if err != nil {
		return nil, fmt.Errorf("identity.repo.User.roles id=%d: %w", row.ID, err)
	}

	roles := make([]domain.Role, len(name))
	for i, n := range name {
		roles[i] = domain.Role(n)
	}
	return r.toDomain(row, roles), nil
}

func (r *userRepository) toDomain(row db.IdentityUser, roles []domain.Role) *domain.User {
	return &domain.User{
		ID:           row.ID,
		FullName:     row.FullName,
		Email:        row.Email,
		Phone:        row.Phone,
		PasswordHash: row.PasswordHash,
		Birthday:     row.Birthday,
		Gender:       fromDBGender(row.Gender),
		IsActive:     row.IsActive,
		Roles:        roles,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
