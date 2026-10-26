package http

import "time"

type RegisterRequest struct {
	FullName string  `json:"full_name" form:"full_name" binding:"required,min=2,max=100"`
	Phone    string  `json:"phone" form:"phone" binding:"required,numeric,min=10,max=15"`
	Email    *string `json:"email" form:"email" binding:"omitempty,email,max=254"`
	Password string  `json:"password" form:"password" binding:"required,min=8,max=128"`
}

type LoginRequest struct {
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Email     *string   `json:"email,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type" example:"Bearer"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type AddressRequest struct {
	ReceiverName string  `json:"receiver_name" form:"receiver_name" binding:"required,min=2,max=100"`
	Phone        string  `json:"phone" form:"phone" binding:"required,numeric,min=10,max=15"`
	Address      string  `json:"address" form:"address" binding:"required,min=5,max=500"`
	PostalCode   *string `json:"postal_code" form:"postal_code" binding:"omitempty,numeric,len=10"`
	No           *string `json:"no" form:"no" binding:"omitempty,max=20"`
	Unit         *string `json:"unit" form:"unit" binding:"omitempty,max=20"`
}

type AddressResponse struct {
	ID           int64     `json:"id"`
	ReceiverName string    `json:"receiver_name"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	No           *string   `json:"no,omitempty"`
	Unit         *string   `json:"unit,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
