package postgres

import (
	"github.com/Hello256World/cloud-native-shop/internal/identity/adapters/postgres/db"
	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
)

func toDBGender(g *domain.Gender) db.NullIdentityGender {
	if g == nil {
		return db.NullIdentityGender{}
	}
	return db.NullIdentityGender{IdentityGender: db.IdentityGender(*g), Valid: true}
}

func fromDBGender(g db.NullIdentityGender) *domain.Gender {
	if !g.Valid {
		return nil
	}
	v := domain.Gender(g.IdentityGender)
	return &v
}
