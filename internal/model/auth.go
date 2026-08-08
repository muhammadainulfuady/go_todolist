package model

import (
	"time"
)

type RequestOtpRequest struct {
	Nama  string `json:"nama" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email,max=100"`
}

type VerifyOtpRequest struct {
	Email   string `json:"email" validate:"required,email,max=100"`
	OtpCode string `json:"otp_code" validate:"required,len=6,numeric"`
}

type OtpSentResponse struct {
	Email string `json:"email"`
}

type UserResponse struct {
	IDUsers     int64     `json:"id_users"`
	Nama        string    `json:"nama"`
	Email       string    `json:"email"`
	FotoProfile *string   `json:"foto_profile"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
