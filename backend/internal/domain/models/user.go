package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                          uuid.UUID  `json:"id" db:"id"`
	Email                       string     `json:"email" db:"email"`
	PasswordHash                string     `json:"-" db:"password_hash"`
	FullName                    string     `json:"full_name" db:"full_name"`
	DateOfBirth                 time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Gender                      string     `json:"gender" db:"gender"`
	Bio                         *string    `json:"bio,omitempty" db:"bio"`
	VideoID                     *uuid.UUID `json:"video_id,omitempty" db:"video_id"`
	Verified                    bool       `json:"verified" db:"verified"`
	EmailVerificationToken      *string    `json:"-" db:"email_verification_token"`
	EmailVerificationExpiresAt  *time.Time `json:"-" db:"email_verification_expires_at"`
	PasswordResetToken          *string    `json:"-" db:"password_reset_token"`
	PasswordResetExpiresAt      *time.Time `json:"-" db:"password_reset_expires_at"`
	LastLoginAt                 *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	Active                      bool       `json:"active" db:"active"`
	CreatedAt                   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                   *time.Time `json:"-" db:"deleted_at"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required,min=2"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
	Gender      string `json:"gender" binding:"required,oneof=male female other"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AuthResponse struct {
	User   *User          `json:"user"`
	Tokens *TokenResponse `json:"tokens"`
}
