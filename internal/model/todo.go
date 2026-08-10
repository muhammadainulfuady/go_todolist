package model

import "time"

type TodoCreateRequest struct {
	Title        string `json:"title" validate:"required,min=1,max=150"`
	Description  string `json:"description"`
	IDPriorities int    `json:"id_priorities" validate:"required,min=1,max=4"`
}

type TodoUpdateRequest struct {
	Title        string `json:"title" validate:"omitempty,min=1,max=150"`
	Description  string `json:"description"`
	IDPriorities int    `json:"id_priorities" validate:"omitempty,min=1,max=4"`
}

type TodoPriority struct {
	IDPriorities int    `json:"id_priorities"`
	Name         string `json:"name"`
}

type TodoResponse struct {
	IDTodos      int64         `json:"id_todos"`
	IDUsers      int64         `json:"id_users"`
	IDPriorities int           `json:"id_priorities"`
	Priority     TodoPriority  `json:"priority"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Description  *string       `json:"description"`
	Image        *string       `json:"image"`
	IsCompleted  bool          `json:"is_completed"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type TodoUpdateResponse struct {
	IDTodos   int64     `json:"id_todos"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoToggleResponse struct {
	IDTodos     int64     `json:"id_todos"`
	Slug        string    `json:"slug"`
	IsCompleted bool      `json:"is_completed"`
	UpdatedAt   time.Time `json:"updated_at"`
}
