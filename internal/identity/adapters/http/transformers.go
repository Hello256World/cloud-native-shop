package http

import (
	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/internal/identity/usecase"
)

func toUserResponse(u *domain.User) UserResponse {
	roles := make([]string, len(u.Roles))

	for i, r := range u.Roles {
		roles[i] = string(r)
	}

	return UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Email:     u.Email,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
	}
}

func toTokenResponse(p *usecase.TokenPair) TokenResponse {
	return TokenResponse{
		AccessToken:      p.AccessToken,
		TokenType:        "Bearer",
		ExpiresAt:        p.ExpiresAt,
		RefreshToken:     p.RefreshToken,
		RefreshExpiresAt: p.RefreshExpiresAt,
	}
}
