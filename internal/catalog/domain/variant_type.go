package domain

type VariantType struct {
	ID           int64
	Name         string
	DataType     string
	DisplayOrder int
	Values       []VariantValue
}

type VariantValue struct {
	ID           int64
	TypeID       int64
	Value        string
	DisplayOrder int
}
