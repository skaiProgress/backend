package auth

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userClaimsKey contextKey = "userClaims"

// Service handles authentication business logic.
type Service struct {
	repo   Repository
	tokens *TokenManager
}

// NewService creates a new auth service.
func NewService(repo Repository, tokens *TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

// Login verifies credentials and returns an access token.
func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.FindAuthUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.EncryptedPassword == "" {
		return nil, ErrInvalidCredentials
	}

	if user.BannedUntil != nil && user.BannedUntil.After(time.Now()) {
		return nil, ErrUserBanned
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	profile, err := s.repo.FindProfileByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	role := "user"
	var fullName *string
	if profile != nil {
		role = profile.Role
		fullName = profile.FullName
		if !profile.IsActive {
			return nil, ErrUserBanned
		}
	}

	token, expiresIn, err := s.tokens.GenerateAccessToken(user.ID, user.Email, role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   expiresIn,
		User: UserPublic{
			ID:       user.ID,
			Email:    user.Email,
			Role:     role,
			FullName: fullName,
		},
	}, nil
}

// Me returns the current user from validated JWT claims.
func (s *Service) Me(ctx context.Context) (*MeResponse, error) {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := s.repo.FindProfileByUserID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	resp := &MeResponse{
		ID:    claims.Subject,
		Email: claims.Email,
		Role:  claims.Role,
	}

	if profile != nil {
		resp.Email = profile.Email
		resp.Role = profile.Role
		resp.FullName = profile.FullName
		resp.IsActive = profile.IsActive
		resp.AvatarURL = profile.AvatarURL
	} else {
		resp.IsActive = true
	}

	return resp, nil
}

// ChangePassword verifies the current password and sets a new one.
func (s *Service) ChangePassword(ctx context.Context, currentPassword, newPassword string) error {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		return ErrUnauthorized
	}

	currentPassword = strings.TrimSpace(currentPassword)
	newPassword = strings.TrimSpace(newPassword)
	if currentPassword == "" || newPassword == "" {
		return errors.New("current_password and new_password are required")
	}
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !hasDigit(newPassword) {
		return errors.New("password must contain at least one digit")
	}

	user, err := s.repo.FindAuthUserByID(ctx, claims.Subject)
	if err != nil {
		return err
	}
	if user == nil || user.EncryptedPassword == "" {
		return ErrWrongPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(currentPassword)); err != nil {
		return ErrWrongPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, claims.Subject, string(hash))
}

func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// ValidateToken parses a bearer token and returns claims.
func (s *Service) ValidateToken(token string) (*Claims, error) {
	return s.tokens.ParseAccessToken(token)
}

// ContextWithClaims stores JWT claims in context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// ClaimsFromContext reads JWT claims from context.
func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(userClaimsKey).(*Claims)
	if !ok || claims == nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}
