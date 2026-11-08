package usecase

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

type cursor struct {
	Sort  domain.ProductSort `json:"s"`
	Price int64              `json:"p,omitempty"`
	ID    int64              `json:"i"`
}

func encodeCursor(sort domain.ProductSort, price, id int64) string {
	b, _ := json.Marshal(cursor{Sort: sort, Price: price, ID: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

func decodeCursor(s string, sort domain.ProductSort) (price, id int64, err error) {
	if s == "" {
		return 0, 0, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, 0, domain.ErrInvalidCursor
	}
	var c cursor
	if err := json.Unmarshal(raw, &c); err != nil || c.Sort != sort || c.ID <= 0 { 
		return 0, 0, domain.ErrInvalidCursor
	}
	if sort != domain.SortNewest && c.Price <= 0 {
		return 0, 0, domain.ErrInvalidCursor
	}
	return c.Price, c.ID, nil
}