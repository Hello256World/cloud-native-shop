package usecase

import "github.com/Hello256World/cloud-native-shop/internal/catalog/ports"

type Service struct {
	categories ports.CategoryRepository
	products   ports.ProductRepository
	types      ports.VariantTypeRepository
	attrs      ports.AttributeRepository
	tx         ports.TxManager
}

func NewService(categories ports.CategoryRepository, products ports.ProductRepository, types ports.VariantTypeRepository, attrs ports.AttributeRepository, tx ports.TxManager) *Service {
	return &Service{
		categories: categories,
		products:   products,
		types:      types,
		attrs:      attrs,
		tx:         tx,
	}
}
