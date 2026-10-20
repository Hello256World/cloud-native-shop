package auth

import (
	"time"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/ports"
	platformauth "github.com/Hello256World/cloud-native-shop/internal/platform/auth"
)

type issuer struct {
	tm *platformauth.TokenManager
}

func NewAccessTokenIssuer(tm *platformauth.TokenManager) ports.AccessTokenIssuer {
	return &issuer{tm: tm}
}

func (i *issuer) Issue(userID int64, roles []domain.Role) (string, time.Time, error) {
	names := make([]string, len(roles))
	for k, r := range roles {
		names[k] = string(r)
	}
	return i.tm.Issue(userID, names)
}
