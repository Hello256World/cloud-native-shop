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

type addressRepository struct {
	pool *pgxpool.Pool
}

func NewAddressRepository(pool *pgxpool.Pool) ports.AddressRepository {
	return &addressRepository{
		pool: pool,
	}
}

func (r *addressRepository) q(c context.Context) *db.Queries {
	return db.New(pg.Executor(c, r.pool))
}

func (r *addressRepository) Create(c context.Context, address *domain.Address) (*domain.Address, error) {
	q := r.q(c)

	row, err := q.InsertAddress(c, db.InsertAddressParams{
		UserID:       address.UserID,
		ReceiverName: address.ReceiverName,
		Phone:        address.Phone,
		Address:      address.Address,
		PostalCode:   address.PostalCode,
		No:           address.No,
		Unit:         address.Unit,
	})

	if err != nil {
		return nil, fmt.Errorf("identity.repo.Address.Create user=%d: %w", address.UserID, err)
	}

	return addressToDomain(row), nil
}

func (r *addressRepository) ListByUser(c context.Context, userID int64) ([]domain.Address, error) {
	q := r.q(c)

	rows, err := q.ListAddressesByUser(c, userID)

	if err != nil {
		return nil, fmt.Errorf("identity.repo.Address.ListByUser user=%d: %w", userID, err)
	}

	addresses := make([]domain.Address, len(rows))

	for i, a := range rows {
		addresses[i] = *addressToDomain(a)
	}

	return addresses, nil
}

func (r *addressRepository) Update(c context.Context, address *domain.Address) (*domain.Address, error) {
	q := r.q(c)

	row, err := q.UpdateAddress(c, db.UpdateAddressParams{
		ID:           address.ID,
		UserID:       address.UserID,
		ReceiverName: address.ReceiverName,
		Phone:        address.Phone,
		Address:      address.Address,
		PostalCode:   address.PostalCode,
		No:           address.No,
		Unit:         address.Unit,
	})

	if err != nil {
		if isNoRows(err) {
			return nil, domain.ErrAddressNotFound
		}
		return nil, fmt.Errorf("identity.repo.Address.Update id=%d: %w", address.ID, err)
	}

	return addressToDomain(row), nil
}

func (r *addressRepository) Delete(c context.Context, userID, id int64) error {
	n, err := r.q(c).DeleteAddress(c, db.DeleteAddressParams{
		ID:     id,
		UserID: userID,
	})

	if err != nil {
		return fmt.Errorf("identity.repo.Address.Delete id=%d: %w", id, err)
	}

	if n == 0 {
		return domain.ErrAddressNotFound
	}

	return nil
}

func addressToDomain(row db.IdentityAddress) *domain.Address {
	return &domain.Address{
		ID:           row.ID,
		UserID:       row.UserID,
		ReceiverName: row.ReceiverName,
		Phone:        row.Phone,
		Address:      row.Address,
		PostalCode:   row.PostalCode,
		No:           row.No,
		Unit:         row.Unit,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
