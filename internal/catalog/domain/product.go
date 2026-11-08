package domain

import "time"

type ProductSort string

const (
	SortNewest    ProductSort = "newest"
	SortPriceAsc  ProductSort = "price_asc"
	SortPriceDesc ProductSort = "price_desc"
)

func (s ProductSort) Valid() bool {
	return s == SortNewest || s == SortPriceAsc || s == SortPriceDesc
}

type Product struct {
	ID              int64
	CategoryID      int64
	Name            string
	Slug            string
	Description     *string
	Thumbnail       string
	ShipmentWeightG int
	IsActive        bool
	Variants        []Variant
	Images          []Image
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ProductSummary struct {
	ID         int64
	CategoryID int64
	Name       string
	Slug       string
	Thumbnail  string
	MinPrice   int64
	IsActive   bool
}

type Image struct {
	ID        int64
	ImageName string
	Priority  int
	IsActive  bool
}

type Variant struct {
	ID        int64
	ProductID int64
	SKU       string
	Price     int64
	IsActive  bool
	Options   []VariantOption
}

type VariantOption struct {
	TypeID  int64
	Type    string
	ValueID int64
	Value   string
}
