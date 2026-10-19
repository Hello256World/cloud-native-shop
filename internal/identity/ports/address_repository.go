package ports

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type AddressRepository interface {
	Create(context.Context, *domain.Address) (*domain.Address, error)
	ListByUser(context.Context, int64) ([]domain.Address, error)
	Update(context.Context, *domain.Address) (*domain.Address, error)
	Delete(c context.Context, userID, id int64) error
}
