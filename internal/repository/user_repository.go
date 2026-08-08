package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go_todolist/internal/entity"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	Create(ctx context.Context, nama, email string) (int64, error)
	UpdateOTP(ctx context.Context, id int64, otpCode string, expiresAt time.Time) error
	ClearOTP(ctx context.Context, id int64) error
	UpdateProfile(ctx context.Context, id int64, nama string, fotoProfile *string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id_users, nama, email, foto_profile, otp_code, otp_expires_at, created_at, updated_at
		FROM users WHERE email = ?
	`, email).Scan(
		&user.IDUsers,
		&user.Nama,
		&user.Email,
		&user.FotoProfile,
		&user.OTPCode,
		&user.OTPExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query user by email: %w", err)
	}
	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id_users, nama, email, foto_profile, otp_code, otp_expires_at, created_at, updated_at
		FROM users WHERE id_users = ?
	`, id).Scan(
		&user.IDUsers,
		&user.Nama,
		&user.Email,
		&user.FotoProfile,
		&user.OTPCode,
		&user.OTPExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("gagal query user by id: %w", err)
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, nama, email string) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO users (nama, email) VALUES (?, ?)", nama, email)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat user: %w", err)
	}
	return result.LastInsertId()
}

func (r *userRepository) UpdateOTP(ctx context.Context, id int64, otpCode string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET otp_code = ?, otp_expires_at = ? WHERE id_users = ?",
		otpCode, expiresAt, id)
	if err != nil {
		return fmt.Errorf("gagal menyimpan OTP: %w", err)
	}
	return nil
}

func (r *userRepository) ClearOTP(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET otp_code = NULL, otp_expires_at = NULL WHERE id_users = ?", id)
	if err != nil {
		return fmt.Errorf("gagal menghapus OTP: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id int64, nama string, fotoProfile *string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET nama = ?, foto_profile = ?, updated_at = NOW() WHERE id_users = ?",
		nama, fotoProfile, id)
	if err != nil {
		return fmt.Errorf("gagal update profil: %w", err)
	}
	return nil
}
