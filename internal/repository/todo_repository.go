package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go_todolist/internal/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) (int64, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
	FindBySlug(ctx context.Context, idUsers int64, slug string) (*entity.Todo, error)
	List(ctx context.Context, idUsers int64, search string, idPriority *int) ([]entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	UpdateImage(ctx context.Context, id int64, image *string) error
	ToggleCompleted(ctx context.Context, idUsers int64, slug string) (*entity.Todo, error)
	Delete(ctx context.Context, id int64) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

const todoSelect = `
	SELECT t.id_todos, t.id_users, t.id_priorities, p.name, t.title, t.slug,
	       t.description, t.image, t.is_completed, t.created_at, t.updated_at
	FROM todos t
	JOIN priorities p ON p.id_priorities = t.id_priorities
`

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO todos (id_users, id_priorities, title, slug, description) VALUES (?, ?, ?, ?, ?)",
		todo.IDUsers, todo.IDPriorities, todo.Title, todo.Slug, todo.Description,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat tugas: %w", err)
	}
	return result.LastInsertId()
}

func (r *todoRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, "SELECT id_todos FROM todos WHERE slug = ?", slug).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("gagal cek slug: %w", err)
	}
	return true, nil
}

func (r *todoRepository) FindBySlug(ctx context.Context, idUsers int64, slug string) (*entity.Todo, error) {
	todo := &entity.Todo{}
	err := r.db.QueryRowContext(ctx, todoSelect+`
		WHERE t.id_users = ? AND t.slug = ?
	`, idUsers, slug).Scan(
		&todo.IDTodos,
		&todo.IDUsers,
		&todo.IDPriorities,
		&todo.PriorityName,
		&todo.Title,
		&todo.Slug,
		&todo.Description,
		&todo.Image,
		&todo.IsCompleted,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query tugas by slug: %w", err)
	}
	return todo, nil
}

func (r *todoRepository) List(ctx context.Context, idUsers int64, search string, idPriority *int) ([]entity.Todo, error) {
	query := todoSelect + " WHERE t.id_users = ?"
	args := []any{idUsers}

	if search != "" {
		query += " AND t.title LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if idPriority != nil {
		query += " AND t.id_priorities = ?"
		args = append(args, *idPriority)
	}

	query += " ORDER BY t.id_todos DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar tugas: %w", err)
	}
	defer rows.Close()

	todos := make([]entity.Todo, 0)
	for rows.Next() {
		var todo entity.Todo
		if err := rows.Scan(
			&todo.IDTodos,
			&todo.IDUsers,
			&todo.IDPriorities,
			&todo.PriorityName,
			&todo.Title,
			&todo.Slug,
			&todo.Description,
			&todo.Image,
			&todo.IsCompleted,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("gagal scan tugas: %w", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal iterasi tugas: %w", err)
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE todos
		SET id_priorities = ?, title = ?, slug = ?, description = ?, image = ?, updated_at = NOW()
		WHERE id_todos = ?
	`, todo.IDPriorities, todo.Title, todo.Slug, todo.Description, todo.Image, todo.IDTodos)
	if err != nil {
		return fmt.Errorf("gagal update tugas: %w", err)
	}
	return nil
}

func (r *todoRepository) UpdateImage(ctx context.Context, id int64, image *string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE todos SET image = ? WHERE id_todos = ?", image, id)
	if err != nil {
		return fmt.Errorf("gagal menyimpan gambar tugas: %w", err)
	}
	return nil
}

func (r *todoRepository) ToggleCompleted(ctx context.Context, idUsers int64, slug string) (*entity.Todo, error) {
	todo := &entity.Todo{}
	err := r.db.QueryRowContext(ctx, todoSelect+`
		WHERE t.id_users = ? AND t.slug = ?
	`, idUsers, slug).Scan(
		&todo.IDTodos,
		&todo.IDUsers,
		&todo.IDPriorities,
		&todo.PriorityName,
		&todo.Title,
		&todo.Slug,
		&todo.Description,
		&todo.Image,
		&todo.IsCompleted,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query tugas untuk toggle: %w", err)
	}

	if _, err := r.db.ExecContext(ctx,
		"UPDATE todos SET is_completed = ?, updated_at = NOW() WHERE id_todos = ?",
		!todo.IsCompleted, todo.IDTodos,
	); err != nil {
		return nil, fmt.Errorf("gagal toggle tugas: %w", err)
	}

	todo.IsCompleted = !todo.IsCompleted
	return todo, nil
}

func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM todos WHERE id_todos = ?", id)
	if err != nil {
		return fmt.Errorf("gagal hapus tugas: %w", err)
	}
	return nil
}
