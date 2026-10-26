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

func toAddressResponse(a *domain.Address) AddressResponse {
	return AddressResponse{
		ID:           a.ID,
		ReceiverName: a.ReceiverName,
		Phone:        a.Phone,
		Address:      a.Address,
		PostalCode:   a.PostalCode,
		No:           a.No,
		Unit:         a.Unit,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func toAddressListResponse(list []domain.Address) []AddressResponse {
	addresses := make([]AddressResponse, len(list))

	for i, val := range list {
		addresses[i] = toAddressResponse(&val)
	}

	return addresses
}
