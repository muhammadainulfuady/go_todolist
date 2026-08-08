package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"go_todolist/internal/entity"
	"go_todolist/internal/helper"
	"go_todolist/internal/model"
	"go_todolist/internal/repository"
)

var ErrUserNotFound = errors.New("user tidak ditemukan")

type ProfileUsecase interface {
	GetProfile(ctx context.Context, userID int64) (*model.UserResponse, error)
	UpdateProfile(ctx context.Context, userID int64, nama string, foto *multipart.FileHeader) (*model.UserResponse, error)
}

type profileUsecase struct {
	userRepo  repository.UserRepository
	baseURL   string
	uploadDir string
}

func NewProfileUsecase(userRepo repository.UserRepository, baseURL, uploadDir string) ProfileUsecase {
	return &profileUsecase{userRepo: userRepo, baseURL: baseURL, uploadDir: uploadDir}
}

func (uc *profileUsecase) GetProfile(ctx context.Context, userID int64) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil profil: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return toUserResponse(uc.baseURL, user), nil
}

func (uc *profileUsecase) UpdateProfile(ctx context.Context, userID int64, nama string, foto *multipart.FileHeader) (*model.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil profil: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if nama == "" {
		nama = user.Nama
	} else if verr := helper.ValidateStruct(&model.ProfileUpdateRequest{Nama: nama}); verr != nil {
		return nil, verr
	}

	newPhoto := user.FotoProfile
	if foto != nil {
		if verr := validateProfileImage(foto); verr != nil {
			return nil, verr
		}
		path, err := helper.SaveImage(foto, uc.uploadDir, fmt.Sprintf("profile_%d", user.IDUsers))
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan foto profil: %w", err)
		}
		newPhoto = &path
	}

	if err := uc.userRepo.UpdateProfile(ctx, user.IDUsers, nama, newPhoto); err != nil {
		return nil, err
	}

	if foto != nil {
		removeOldPhoto(user.FotoProfile, *newPhoto)
	}

	user.Nama = nama
	user.FotoProfile = newPhoto
	return toUserResponse(uc.baseURL, user), nil
}

func validateProfileImage(f *multipart.FileHeader) *helper.ValidationError {
	switch err := helper.ValidateImage(f); {
	case errors.Is(err, helper.ErrImageTooLarge):
		return &helper.ValidationError{FieldErrors: map[string]string{
			"foto_profile": "Ukuran gambar maksimal 2MB",
		}}
	case errors.Is(err, helper.ErrInvalidImageType):
		return &helper.ValidationError{FieldErrors: map[string]string{
			"foto_profile": "Format gambar harus .jpg atau .png",
		}}
	case err != nil:
		return &helper.ValidationError{FieldErrors: map[string]string{
			"foto_profile": "File foto tidak valid",
		}}
	default:
		return nil
	}
}

func removeOldPhoto(old *string, newPath string) {
	if old == nil || *old == "" || *old == newPath {
		return
	}
	helper.RemoveFile(*old)
}

func toUserResponse(baseURL string, user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		IDUsers:     user.IDUsers,
		Nama:        user.Nama,
		Email:       user.Email,
		FotoProfile: photoURL(baseURL, user.FotoProfile),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func photoURL(baseURL string, path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	url := baseURL + "/" + strings.TrimPrefix(*path, "/")
	return &url
}
