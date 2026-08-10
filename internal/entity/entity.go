package entity

import "time"

type User struct {
	IDUsers      int64      `json:"id_users"`
	Nama         string     `json:"nama"`
	Email        string     `json:"email"`
	FotoProfile  *string    `json:"foto_profile"`
	OTPCode      *string    `json:"-"`
	OTPExpiresAt *time.Time `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Priority struct {
	IDPriorities int    `json:"id_priorities"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type Todo struct {
	IDTodos      int64     `json:"id_todos"`
	IDUsers      int64     `json:"id_users"`
	IDPriorities int       `json:"id_priorities"`
	PriorityName string    `json:"-"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description"`
	Image        *string   `json:"image"`
	IsCompleted  bool      `json:"is_completed"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
