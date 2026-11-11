package http

import (
	"strconv"

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

func toVariantTypeResponse(t domain.VariantType) VariantTypeResponse {
	vals := make([]VariantValueResponse, len(t.Values))
	for i, v := range t.Values {
		vals[i] = VariantValueResponse{ID: v.ID, Value: v.Value, DisplayOrder: v.DisplayOrder}
	}
	return VariantTypeResponse{
		ID: t.ID, Name: t.Name, DataType: t.DataType, DisplayOrder: t.DisplayOrder, Values: vals,
	}
}

func toVariantInput(r VariantRequest) usecase.VariantInput {
	return usecase.VariantInput{SKU: r.SKU, Price: r.Price, IsActive: r.IsActive, ValueIDs: r.ValueIDs}
}

func toImageInput(r ImageRequest) usecase.ImageInput {
	return usecase.ImageInput{ImageName: r.ImageName, Priority: r.Priority, IsActive: r.IsActive}
}

func toCreateProductInput(r CreateProductRequest) usecase.CreateProductInput {
	in := usecase.CreateProductInput{ProductInput: toProductInput(r.ProductRequest)}
	for _, v := range r.Variants {
		in.Variants = append(in.Variants, toVariantInput(v))
	}
	for _, im := range r.Images {
		in.Images = append(in.Images, toImageInput(im))
	}
	return in
}

func toVariantResponse(v domain.Variant) VariantResponse {
	opts := make([]OptionResponse, len(v.Options))
	for i, o := range v.Options {
		opts[i] = OptionResponse{TypeID: o.TypeID, Type: o.Type, ValueID: o.ValueID, Value: o.Value}
	}
	return VariantResponse{ID: v.ID, SKU: v.SKU, Price: v.Price, IsActive: v.IsActive, Options: opts}
}

func toImageResponse(im domain.Image) ImageResponse {
	return ImageResponse{ID: im.ID, ImageName: im.ImageName, Priority: im.Priority, IsActive: im.IsActive}
}

func toProductResponse(p *domain.Product) ProductResponse {
	variants := make([]VariantResponse, len(p.Variants))
	for i, v := range p.Variants {
		variants[i] = toVariantResponse(v)
	}
	images := make([]ImageResponse, len(p.Images))
	for i, im := range p.Images {
		images[i] = toImageResponse(im)
	}
	specs := make([]SpecResponse, len(p.Attributes))
	for i, a := range p.Attributes {
		specs[i] = toSpecResponse(a)
	}
	return ProductResponse{
		ID: p.ID, CategoryID: p.CategoryID, Name: p.Name, Slug: p.Slug, Description: p.Description,
		Thumbnail: p.Thumbnail, ShipmentWeightG: p.ShipmentWeightG, IsActive: p.IsActive,
		Variants: variants, Images: images, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,Specs: specs,
	}
}

func toListProductsInput(q ProductListQuery) usecase.ListProductsInput {
	return usecase.ListProductsInput{
		Cursor: q.Cursor, Limit: q.Limit, CategoryID: q.CategoryID, Search: q.Q,
		Sort: domain.ProductSort(q.Sort), PriceMin: q.PriceMin, PriceMax: q.PriceMax,
	}
}
func toProductPageResponse(p *usecase.ProductPage) ProductPageResponse {
	items := make([]ProductSummaryResponse, len(p.Items))
	for i, v := range p.Items {
		items[i] = ProductSummaryResponse{
			ID: v.ID, CategoryID: v.CategoryID, Name: v.Name, Slug: v.Slug,
			Thumbnail: v.Thumbnail, MinPrice: v.MinPrice, IsActive: v.IsActive,
		}
	}
	out := ProductPageResponse{Items: items}
	if p.NextCursor != "" {
		next := p.NextCursor
		out.NextCursor = &next
	}
	return out
}

func toAdminListProductsInput(q AdminProductListQuery) usecase.ListProductsInput {
	return usecase.ListProductsInput{Cursor: q.Cursor, Limit: q.Limit, CategoryID: q.CategoryID, Search: q.Q}
}

func toProductInput(r ProductRequest) usecase.ProductInput {
	in := usecase.ProductInput{
		CategoryID: r.CategoryID, Name: r.Name, Slug: r.Slug, Description: r.Description,
		Thumbnail: r.Thumbnail, ShipmentWeightG: r.ShipmentWeightG, IsActive: r.IsActive,
	}
	if r.Attributes != nil {
		in.Attributes = make([]usecase.AttributeInput, len(r.Attributes))
		for i, a := range r.Attributes {
			in.Attributes[i] = usecase.AttributeInput{AttributeID: a.AttributeID, Value: a.Value}
		}
	}
	return in
}

func toAttributeDefinitionResponse(d domain.AttributeDefinition) AttributeDefinitionResponse {
	return AttributeDefinitionResponse{
		ID: d.ID, Code: d.Code, Name: d.Name, ValueType: string(d.ValueType),
		Unit: d.Unit, IsFilterable: d.IsFilterable, DisplayOrder: d.DisplayOrder,
	}
}

func toSpecResponse(a domain.ProductAttribute) SpecResponse {
	val := ""
	switch {
	case a.ValueNumber != nil:
		val = strconv.FormatFloat(*a.ValueNumber, 'f', -1, 64)
	case a.ValueText != nil:
		val = *a.ValueText
	}
	return SpecResponse{
		AttributeID: a.AttributeID, Code: a.Code, Name: a.Name,
		ValueType: string(a.ValueType), Value: val, Unit: a.Unit,
	}
}
