package domain

import "time"

type Category struct {
	ID        int64
	ParentID  *int64
	Name      string
	Slug      string
	ImageName *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
