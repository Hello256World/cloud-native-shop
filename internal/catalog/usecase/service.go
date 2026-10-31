package usecase

import "github.com/Hello256World/cloud-native-shop/internal/catalog/ports"

type Service struct {
	categories ports.CategoryRepository
	tx  ports.TxManager
}

func NewService(categories ports.CategoryRepository, tx ports.TxManager) *Service {
	return &Service{
		categories: categories,
		tx:  tx,
	}
}
