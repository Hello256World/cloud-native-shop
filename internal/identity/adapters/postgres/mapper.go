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

// func toDBPasswordHash(hash *string) pgtype.Text {
// 	if hash == nil {
// 		return pgtype.Text{}
// 	}

// 	return pgtype.Text{String: *hash, Valid: true}
// }

// func fromDBPasswordHash(hash pgtype.Text) *string {
// 	if !hash.Valid {
// 		return nil
// 	}

// 	return &hash.String
// }

// func toDBBirthday(date *time.Time) pgtype.Date {
// 	if date == nil {
// 		return pgtype.Date{}
// 	}

// 	return pgtype.Date{Time: *date, Valid: true}
// }

// func fromDBBirthday(date pgtype.Date) *time.Time {
// 	if !date.Valid {
// 		return nil
// 	}

// 	return &date.Time
// }
