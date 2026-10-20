package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("auth.randomID: %w", err)
	}
	return hex.EncodeToString(b), nil
}
