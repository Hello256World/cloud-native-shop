package domain

import "time"

type Address struct {
	ID           int64
	UserID       int64
	ReceiverName string
	Phone        string
	Address      string
	PostalCode   *string
	No           *string
	Unit         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
