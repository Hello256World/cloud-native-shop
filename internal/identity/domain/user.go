package domain

import "time"

type Role string

const (
	RoleCustomer   Role = "customer"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "super_admin"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	ID           int64
	FullName     string
	Email        *string
	Phone        string
	PasswordHash *string
	Birthday     *time.Time
	Gender       *Gender
	IsActive     bool
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) HasRole(r Role) bool {
	for _, val := range u.Roles {
		if val == r {
			return true
		}
	}
	return false
}

