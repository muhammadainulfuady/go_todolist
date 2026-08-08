package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go_todolist/internal/entity"
)

type PriorityRepository interface {
	FindAll(ctx context.Context) ([]entity.Priority, error)
}

type priorityRepository struct {
	db *sql.DB
}

func NewPriorityRepository(db *sql.DB) PriorityRepository {
	return &priorityRepository{db: db}
}

func (r *priorityRepository) FindAll(ctx context.Context) ([]entity.Priority, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id_priorities, name, description FROM priorities ORDER BY id_priorities")
	if err != nil {
		return nil, fmt.Errorf("gagal query prioritas: %w", err)
	}
	defer rows.Close()

	priorities := make([]entity.Priority, 0)
	for rows.Next() {
		var p entity.Priority
		if err := rows.Scan(&p.IDPriorities, &p.Name, &p.Description); err != nil {
			return nil, fmt.Errorf("gagal scan prioritas: %w", err)
		}
		priorities = append(priorities, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal iterasi prioritas: %w", err)
	}
	return priorities, nil
}
