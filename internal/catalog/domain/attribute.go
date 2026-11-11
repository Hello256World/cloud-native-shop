package domain

type AttributeValueType string

const (
	AttrText   AttributeValueType = "text"
	AttrNumber AttributeValueType = "number"
)

func (t AttributeValueType) IsValid() bool {
	return t == AttrText || t == AttrNumber
}

type AttributeDefinition struct {
	ID           int64
	Code         string
	Name         string
	ValueType    AttributeValueType
	Unit         *string
	IsFilterable bool
	DisplayOrder int
}

type ProductAttribute struct {
	AttributeID  int64
	Code         string
	Name         string
	ValueType    AttributeValueType
	Unit         *string
	DisplayOrder int
	ValueText    *string
	ValueNumber  *float64
}
