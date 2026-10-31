package http

import "time"

type CategoryRequest struct {
	ParentID  *int64  `json:"parent_id" binding:"omitempty,gt=0"`
	Name      string  `json:"name" binding:"required,min=2,max=100"`
	Slug      string  `json:"slug" binding:"required,min=2,max=120,slug"`
	ImageName *string `json:"image_name" binding:"omitempty,max=255"`
	IsActive  bool    `json:"is_active"`
}

type CategoryResponse struct {
	ID        int64     `json:"id"`
	ParentID  *int64    `json:"parent_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	ImageName *string   `json:"image_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
