package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"go_todolist/internal/entity"
	mailgate "go_todolist/internal/gateway/email"
	"go_todolist/internal/helper"
	"go_todolist/internal/model"
	"go_todolist/internal/repository"
	"go_todolist/internal/security"
)

var (
	ErrEmailNotFound = errors.New("email tidak terdaftar")
	ErrInvalidOtp    = errors.New("kode OTP salah atau telah kedaluwarsa")
)

type AuthUsecase interface {
	RequestOtp(ctx context.Context, req *model.RequestOtpRequest) (*model.OtpSentResponse, error)
	VerifyOtp(ctx context.Context, req *model.VerifyOtpRequest) (*model.AuthResponse, error)
}

type authUsecase struct {
	userRepo   repository.UserRepository
	mailer     *mailgate.Mailer
	jwtSecret  string
	jwtExpires time.Duration
	otpExpires time.Duration
}

func NewAuthUsecase(userRepo repository.UserRepository, mailer *mailgate.Mailer,
	jwtSecret string, jwtExpires, otpExpires time.Duration) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		mailer:     mailer,
		jwtSecret:  jwtSecret,
		jwtExpires: jwtExpires,
		otpExpires: otpExpires,
	}
}

func (uc *authUsecase) RequestOtp(ctx context.Context, req *model.RequestOtpRequest) (*model.OtpSentResponse, error) {
	if verr := helper.ValidateStruct(req); verr != nil {
		return nil, verr
	}

	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari user: %w", err)
	}
	if user == nil {
		id, err := uc.userRepo.Create(ctx, req.Nama, req.Email)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat user: %w", err)
		}
		user = &entity.User{IDUsers: id, Nama: req.Nama, Email: req.Email}
	}

	code, err := generateOTP()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat kode OTP: %w", err)
	}

	if err := uc.userRepo.UpdateOTP(ctx, user.IDUsers, code, time.Now().Add(uc.otpExpires)); err != nil {
		return nil, fmt.Errorf("gagal menyimpan OTP: %w", err)
	}

	if err := uc.mailer.SendOTP(user.Email, code); err != nil {
		return nil, fmt.Errorf("gagal mengirim email OTP: %w", err)
	}

	return &model.OtpSentResponse{Email: user.Email}, nil
}

func (uc *authUsecase) VerifyOtp(ctx context.Context, req *model.VerifyOtpRequest) (*model.AuthResponse, error) {
	if verr := helper.ValidateStruct(req); verr != nil {
		return nil, verr
	}

	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari user: %w", err)
	}
	if user == nil {
		return nil, ErrEmailNotFound
	}

	if user.OTPCode == nil || user.OTPExpiresAt == nil ||
		*user.OTPCode != req.OtpCode || time.Now().After(*user.OTPExpiresAt) {
		return nil, ErrInvalidOtp
	}

	if err := uc.userRepo.ClearOTP(ctx, user.IDUsers); err != nil {
		return nil, fmt.Errorf("gagal menghapus OTP: %w", err)
	}

	token, err := security.GenerateToken(uc.jwtSecret, uc.jwtExpires, *user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	return &model.AuthResponse{
		Token: token,
		User: model.UserResponse{
			IDUsers:     user.IDUsers,
			Nama:        user.Nama,
			Email:       user.Email,
			FotoProfile: user.FotoProfile,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
	}, nil
}

func generateOTP() (string, error) {
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		sb.WriteByte('0' + byte(n.Int64()))
	}
	return sb.String(), nil
}
