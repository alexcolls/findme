package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/alexcolls/findme/internal/domain/models"
	"github.com/alexcolls/findme/internal/repository/postgres"
	"github.com/alexcolls/findme/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	userRepo   postgres.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo postgres.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		DateOfBirth:  dob,
		Gender:       req.Gender,
		Verified:     false,
		Active:       true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate email verification token
	verificationToken, err := generateToken()
	if err == nil {
		expiresAt := sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true}
		s.userRepo.SetEmailVerificationToken(ctx, user.ID, verificationToken, expiresAt)
		// TODO: Send verification email
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: user,
		Tokens: &models.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    s.jwtManager.GetAccessTokenDuration(),
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.Active {
		return nil, fmt.Errorf("account is inactive")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: user,
		Tokens: &models.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    s.jwtManager.GetAccessTokenDuration(),
			TokenType:    "Bearer",
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Generate new access token
	accessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtManager.GetAccessTokenDuration(),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	return s.userRepo.VerifyEmail(ctx, token)
}

func (s *authService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Check if user exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return "", nil
	}

	// Generate reset token
	resetToken, err := generateToken()
	if err != nil {
		return "", err
	}

	expiresAt := sql.NullTime{Time: time.Now().Add(1 * time.Hour), Valid: true}
	if err := s.userRepo.SetPasswordResetToken(ctx, email, resetToken, expiresAt); err != nil {
		return "", err
	}

	// TODO: Send password reset email
	return resetToken, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.ResetPassword(ctx, token, string(hashedPassword))
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
