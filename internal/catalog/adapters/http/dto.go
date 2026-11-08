package http

import "time"

type CategoryRequest struct {
	ParentID  *int64  `json:"parent_id" binding:"omitempty,gt=0"`
	Name      string  `json:"name" binding:"required,min=2,max=100"`
	Slug      string  `json:"slug" binding:"required,min=2,max=120,slug"`
	ImageName *string `json:"image_name" binding:"omitempty,max=255"`
	IsActive  bool    `json:"is_active"`
}

type CategoryResponse struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ImageName *string   `json:"image_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantTypeRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=50"`
	DataType     string `json:"data_type" binding:"required,oneof=text color"`
	DisplayOrder int    `json:"display_order" binding:"min=0"`
}

type VariantValueRequest struct {
	Value        string `json:"value" binding:"required,min=1,max=100"`
	DisplayOrder int    `json:"display_order" binding:"min=0"`
}

type VariantValueResponse struct {
	ID           int64  `json:"id"`
	Value        string `json:"value"`
	DisplayOrder int    `json:"display_order"`
}

type VariantTypeResponse struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	DataType     string                 `json:"data_type"`
	DisplayOrder int                    `json:"display_order"`
	Values       []VariantValueResponse `json:"values"`
}

type ProductRequest struct {
	CategoryID      int64   `json:"category_id" binding:"required,gt=0"`
	Name            string  `json:"name" binding:"required,min=2,max=200"`
	Slug            string  `json:"slug" binding:"required,min=2,max=220,slug"`
	Description     *string `json:"description" binding:"omitempty,max=10000"`
	Thumbnail       string  `json:"thumbnail" binding:"required,max=255"`
	ShipmentWeightG int     `json:"shipment_weight_g" binding:"required,gt=0"`
	IsActive        bool    `json:"is_active"`
}

type VariantRequest struct {
	SKU      string  `json:"sku" binding:"required,min=2,max=64"`
	Price    int64   `json:"price" binding:"required,gt=0"`
	IsActive bool    `json:"is_active"`
	ValueIDs []int64 `json:"value_ids" binding:"omitempty,max=10,dive,gt=0"`
}

type VariantUpdateRequest struct {
	Price    int64 `json:"price" binding:"required,gt=0"`
	IsActive bool  `json:"is_active"`
}

type ImageRequest struct {
	ImageName string `json:"image_name" binding:"required,max=255"`
	Priority  int    `json:"priority" binding:"min=0"`
	IsActive  bool   `json:"is_active"`
}

type CreateProductRequest struct {
	ProductRequest
	Variants []VariantRequest `json:"variants" binding:"required,min=1,max=100,dive"`
	Images   []ImageRequest   `json:"images" binding:"omitempty,max=20,dive"`
}

type OptionResponse struct {
	TypeID  int64  `json:"type_id"`
	Type    string `json:"type"`
	ValueID int64  `json:"value_id"`
	Value   string `json:"value"`
}

type VariantResponse struct {
	ID       int64            `json:"id"`
	SKU      string           `json:"sku"`
	Price    int64            `json:"price"`
	IsActive bool             `json:"is_active"`
	Options  []OptionResponse `json:"options"`
}

type ImageResponse struct {
	ID        int64  `json:"id"`
	ImageName string `json:"image_name"`
	Priority  int    `json:"priority"`
	IsActive  bool   `json:"is_active"`
}

type ProductResponse struct {
	ID              int64             `json:"id"`
	CategoryID      int64             `json:"category_id"`
	Name            string            `json:"name"`
	Slug            string            `json:"slug"`
	Description     *string           `json:"description,omitempty"`
	Thumbnail       string            `json:"thumbnail"`
	ShipmentWeightG int               `json:"shipment_weight_g"`
	IsActive        bool              `json:"is_active"`
	Variants        []VariantResponse `json:"variants"`
	Images          []ImageResponse   `json:"images"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type ProductListQuery struct {
	Cursor     string `json:"cursor" form:"cursor" binding:"omitempty,max=128"`
	Limit      int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=50"`
	CategoryID *int64 `json:"category_id" form:"category_id" binding:"omitempty,gt=0"`
	Q          string `json:"q" form:"q" binding:"omitempty,max=100"`
	Sort       string `json:"sort" form:"sort" binding:"omitempty,oneof=newest price_asc price_desc"`
	PriceMin   int64  `json:"price_min" form:"price_min" binding:"omitempty,gt=0"`
	PriceMax   int64  `json:"price_max" form:"price_max" binding:"omitempty,gt=0,gtefield=PriceMin"`
}

type ProductSummaryResponse struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Thumbnail  string `json:"thumbnail"`
	MinPrice   int64  `json:"min_price"`
	IsActive   bool   `json:"is_active"`
}

type ProductPageResponse struct {
	Items      []ProductSummaryResponse `json:"items"`
	NextCursor *string                  `json:"next_cursor"`
}

type AdminProductListQuery struct {
	Cursor     string `json:"cursor" form:"cursor" binding:"omitempty,max=128"`
	Limit      int    `json:"limit" form:"limit" binding:"omitempty,min=1,max=50"`
	CategoryID *int64 `json:"category_id" form:"category_id" binding:"omitempty,gt=0"`
	Q          string `json:"q" form:"q" binding:"omitempty,max=100"`
}
