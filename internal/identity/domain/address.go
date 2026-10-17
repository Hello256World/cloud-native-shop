package domain

import "time"

type Address struct {
	ID           uint64
	UserID       uint64
	ReceiverName string
	Phone        string
	Address      string
	PostalCode   *string
	No           *string
	Unit         *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}