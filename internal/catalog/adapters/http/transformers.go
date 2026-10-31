package http

import (
	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
	"github.com/Hello256World/cloud-native-shop/internal/catalog/usecase"
)

func toCategoryInput(r CategoryRequest) usecase.CategoryInput {
	return usecase.CategoryInput{
		ParentID: r.ParentID, Name: r.Name, Slug: r.Slug,
		ImageName: r.ImageName, IsActive: r.IsActive,
	}
}

func toCategoryResponse(c *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID: c.ID, ParentID: c.ParentID, Name: c.Name, Slug: c.Slug,
		ImageName: c.ImageName, IsActive: c.IsActive,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryListResponse(list []domain.Category) []CategoryResponse {
	out := make([]CategoryResponse, len(list)) 
	for i := range list {
		out[i] = toCategoryResponse(&list[i])
	}
	return out
}
