package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"go_todolist/internal/entity"
	"go_todolist/internal/helper"
	"go_todolist/internal/model"
	"go_todolist/internal/repository"
)

var ErrTodoNotFound = errors.New("tugas tidak ditemukan")

type TodoUsecase interface {
	Create(ctx context.Context, userID int64, req model.TodoCreateRequest, image *multipart.FileHeader) (*model.TodoResponse, error)
	List(ctx context.Context, userID int64, search string, idPriority *int) ([]model.TodoResponse, error)
	GetBySlug(ctx context.Context, userID int64, slug string) (*model.TodoResponse, error)
	Update(ctx context.Context, userID int64, slug string, req model.TodoUpdateRequest, image *multipart.FileHeader) (*model.TodoUpdateResponse, error)
	Toggle(ctx context.Context, userID int64, slug string) (*model.TodoToggleResponse, error)
	Delete(ctx context.Context, userID int64, slug string) error
}

type todoUsecase struct {
	todoRepo  repository.TodoRepository
	baseURL   string
	uploadDir string
}

func NewTodoUsecase(todoRepo repository.TodoRepository, baseURL, uploadDir string) TodoUsecase {
	return &todoUsecase{todoRepo: todoRepo, baseURL: baseURL, uploadDir: uploadDir}
}

func (uc *todoUsecase) Create(ctx context.Context, userID int64, req model.TodoCreateRequest, image *multipart.FileHeader) (*model.TodoResponse, error) {
	if verr := helper.ValidateStruct(&req); verr != nil {
		return nil, verr
	}
	if image != nil {
		if verr := validateTodoImage(image); verr != nil {
			return nil, verr
		}
	}

	slug, err := helper.EnsureUniqueSlug(ctx, uc.todoRepo, helper.GenerateSlug(req.Title))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat slug: %w", err)
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	todo := &entity.Todo{
		IDUsers:      userID,
		IDPriorities: req.IDPriorities,
		Title:        req.Title,
		Slug:         slug,
		Description:  description,
	}

	id, err := uc.todoRepo.Create(ctx, todo)
	if err != nil {
		return nil, err
	}
	todo.IDTodos = id

	if image != nil {
		path, err := helper.SaveImage(image, uc.uploadDir, fmt.Sprintf("todo_%d", id))
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan gambar tugas: %w", err)
		}
		todo.Image = &path
		if err := uc.todoRepo.UpdateImage(ctx, id, &path); err != nil {
			return nil, err
		}
	}

	fresh, err := uc.todoRepo.FindBySlug(ctx, userID, slug)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil tugas yang baru dibuat: %w", err)
	}
	return toTodoResponse(uc.baseURL, fresh), nil
}

func (uc *todoUsecase) List(ctx context.Context, userID int64, search string, idPriority *int) ([]model.TodoResponse, error) {
	todos, err := uc.todoRepo.List(ctx, userID, search, idPriority)
	if err != nil {
		return nil, err
	}

	resp := make([]model.TodoResponse, 0, len(todos))
	for i := range todos {
		resp = append(resp, *toTodoResponse(uc.baseURL, &todos[i]))
	}
	return resp, nil
}

func (uc *todoUsecase) GetBySlug(ctx context.Context, userID int64, slug string) (*model.TodoResponse, error) {
	todo, err := uc.todoRepo.FindBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrTodoNotFound
	}
	return toTodoResponse(uc.baseURL, todo), nil
}

func (uc *todoUsecase) Update(ctx context.Context, userID int64, slug string, req model.TodoUpdateRequest, image *multipart.FileHeader) (*model.TodoUpdateResponse, error) {
	if verr := helper.ValidateStruct(&req); verr != nil {
		return nil, verr
	}
	if image != nil {
		if verr := validateTodoImage(image); verr != nil {
			return nil, verr
		}
	}

	todo, err := uc.todoRepo.FindBySlug(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrTodoNotFound
	}

	if req.Title != "" && req.Title != todo.Title {
		newSlug := helper.GenerateSlug(req.Title)
		if newSlug != todo.Slug {
			unique, err := helper.EnsureUniqueSlug(ctx, uc.todoRepo, newSlug)
			if err != nil {
				return nil, fmt.Errorf("gagal membuat slug: %w", err)
			}
			todo.Slug = unique
		}
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = &req.Description
	}
	if req.IDPriorities != 0 {
		todo.IDPriorities = req.IDPriorities
	}

	oldImage := todo.Image
	if image != nil {
		path, err := helper.SaveImage(image, uc.uploadDir, fmt.Sprintf("todo_%d", todo.IDTodos))
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan gambar tugas: %w", err)
		}
		todo.Image = &path
	}

	if err := uc.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	if image != nil {
		removeOldPhoto(oldImage, *todo.Image)
	}

	updatedAt := time.Now()
	return &model.TodoUpdateResponse{
		IDTodos:   todo.IDTodos,
		Title:     todo.Title,
		Slug:      todo.Slug,
		UpdatedAt: updatedAt,
	}, nil
}

func (uc *todoUsecase) Toggle(ctx context.Context, userID int64, slug string) (*model.TodoToggleResponse, error) {
	todo, err := uc.todoRepo.ToggleCompleted(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, ErrTodoNotFound
	}
	return &model.TodoToggleResponse{
		IDTodos:     todo.IDTodos,
		Slug:        todo.Slug,
		IsCompleted: todo.IsCompleted,
		UpdatedAt:   time.Now(),
	}, nil
}

func (uc *todoUsecase) Delete(ctx context.Context, userID int64, slug string) error {
	todo, err := uc.todoRepo.FindBySlug(ctx, userID, slug)
	if err != nil {
		return err
	}
	if todo == nil {
		return ErrTodoNotFound
	}

	if err := uc.todoRepo.Delete(ctx, todo.IDTodos); err != nil {
		return err
	}
	removeOldPhoto(todo.Image, "")
	return nil
}

func validateTodoImage(f *multipart.FileHeader) *helper.ValidationError {
	switch err := helper.ValidateImage(f); {
	case errors.Is(err, helper.ErrImageTooLarge):
		return &helper.ValidationError{FieldErrors: map[string]string{
			"image": "Ukuran gambar maksimal 2MB",
		}}
	case errors.Is(err, helper.ErrInvalidImageType):
		return &helper.ValidationError{FieldErrors: map[string]string{
			"image": "Format gambar harus .jpg atau .png",
		}}
	case err != nil:
		return &helper.ValidationError{FieldErrors: map[string]string{
			"image": "File gambar tidak valid",
		}}
	default:
		return nil
	}
}

func toTodoResponse(baseURL string, todo *entity.Todo) *model.TodoResponse {
	return &model.TodoResponse{
		IDTodos:      todo.IDTodos,
		IDUsers:      todo.IDUsers,
		IDPriorities: todo.IDPriorities,
		Priority: model.TodoPriority{
			IDPriorities: todo.IDPriorities,
			Name:         todo.PriorityName,
		},
		Title:       todo.Title,
		Slug:        todo.Slug,
		Description: todo.Description,
		Image:       photoURL(baseURL, todo.Image),
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
