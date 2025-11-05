package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alexcolls/findme/internal/domain/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	SetEmailVerificationToken(ctx context.Context, id uuid.UUID, token string, expiresAt sql.NullTime) error
	VerifyEmail(ctx context.Context, token string) error
	SetPasswordResetToken(ctx context.Context, email string, token string, expiresAt sql.NullTime) error
	ResetPassword(ctx context.Context, token string, passwordHash string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, date_of_birth, gender, bio, verified, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		user.Email, user.PasswordHash, user.FullName, user.DateOfBirth,
		user.Gender, user.Bio, user.Verified, user.Active,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, date_of_birth, gender, bio,
		       video_id, verified, last_login_at, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.DateOfBirth, &user.Gender, &user.Bio, &user.VideoID,
		&user.Verified, &user.LastLoginAt, &user.Active,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, date_of_birth, gender, bio,
		       video_id, verified, last_login_at, active, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.DateOfBirth, &user.Gender, &user.Bio, &user.VideoID,
		&user.Verified, &user.LastLoginAt, &user.Active,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET full_name = $1, bio = $2, updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query, user.FullName, user.Bio, user.ID).Scan(&user.UpdatedAt)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) SetEmailVerificationToken(ctx context.Context, id uuid.UUID, token string, expiresAt sql.NullTime) error {
	query := `
		UPDATE users
		SET email_verification_token = $1, email_verification_expires_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt, id)
	return err
}

func (r *userRepository) VerifyEmail(ctx context.Context, token string) error {
	query := `
		UPDATE users
		SET verified = true, email_verification_token = NULL, email_verification_expires_at = NULL
		WHERE email_verification_token = $1 AND email_verification_expires_at > NOW()
	`
	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("invalid or expired token")
	}
	return nil
}

func (r *userRepository) SetPasswordResetToken(ctx context.Context, email string, token string, expiresAt sql.NullTime) error {
	query := `
		UPDATE users
		SET password_reset_token = $1, password_reset_expires_at = $2
		WHERE email = $3 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt, email)
	return err
}

func (r *userRepository) ResetPassword(ctx context.Context, token string, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, password_reset_token = NULL, password_reset_expires_at = NULL
		WHERE password_reset_token = $2 AND password_reset_expires_at > NOW()
	`
	result, err := r.db.ExecContext(ctx, query, passwordHash, token)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("invalid or expired token")
	}
	return nil
}
