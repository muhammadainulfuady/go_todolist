package model

type ProfileUpdateRequest struct {
	Nama string `json:"nama" validate:"omitempty,min=3,max=100"`
}
