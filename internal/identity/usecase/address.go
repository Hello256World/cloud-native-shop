package usecase

import (
	"context"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

type AddressInput struct {
	ReceiverName string
	Phone        string
	Address      string
	PostalCode   *string
	No           *string
	Unit         *string
}

func (i AddressInput) toDomain(userID, id int64) *domain.Address {
	return &domain.Address{
		ID:           id,
		UserID:       userID,
		ReceiverName: i.ReceiverName,
		Phone:        i.Phone,
		Address:      i.Address,
		PostalCode:   i.PostalCode,
		No:           i.No,
		Unit:         i.Unit,
	}
}

func (s *Service) AddAddress(c context.Context, userID int64, in AddressInput) (*domain.Address, error) {
	return s.address.Create(c, in.toDomain(userID, 0))
}

func (s *Service) ListAddresses(c context.Context, userID int64) ([]domain.Address, error) {
	return s.address.ListByUser(c, userID)
}

func (s *Service) UpdateAddress(c context.Context, userID, id int64, in AddressInput) (*domain.Address, error) {
	return s.address.Update(c, in.toDomain(userID, id))
}

func (s *Service) DeleteAddress(c context.Context, userID, id int64) error {
	return s.address.Delete(c, userID, id)
}
