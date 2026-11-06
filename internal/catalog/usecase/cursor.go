package usecase

import (
	"encoding/base64"
	"strconv"

	"github.com/Hello256World/cloud-native-shop/internal/catalog/domain"
)

func encodeCursor(id int64) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

func decodeCursor(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, domain.ErrInvalidCursor
	}

	id, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.ErrInvalidCursor
	}
	return id, nil
}
