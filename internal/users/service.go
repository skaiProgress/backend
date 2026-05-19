package users

import (
	"context"
	"errors"
	"strings"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service handles admin user management.
type Service struct {
	repo Repository
}

// NewService creates an admin users service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) callerFromContext(ctx context.Context) (*Caller, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if claims.Role != "admin" && claims.Role != "super_admin" {
		return nil, ErrForbidden
	}
	return &Caller{ID: claims.Subject, Role: claims.Role}, nil
}

// ListUsers returns all profiles for admin UI.
func (s *Service) ListUsers(ctx context.Context, search string) ([]Profile, error) {
	if _, err := s.callerFromContext(ctx); err != nil {
		return nil, err
	}
	return s.repo.ListProfiles(ctx, search)
}

// GetUser returns one profile by id.
func (s *Service) GetUser(ctx context.Context, userID string) (*Profile, error) {
	if _, err := s.callerFromContext(ctx); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, ErrInvalidInput
	}
	p, err := s.repo.GetProfileByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, pgx.ErrNoRows
	}
	return p, nil
}

// AddUser creates auth.users + profiles row.
func (s *Service) AddUser(ctx context.Context, req AddUserRequest) (*AddUserResponse, error) {
	caller, err := s.callerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(req.Email)
	password := req.Password
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	role := req.Role
	if role == "" {
		role = "user"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if role == "super_admin" {
		return nil, ErrForbidden
	}
	if role == "admin" && caller.Role != "super_admin" {
		return nil, errors.New("only super_admin can create admin users")
	}
	if role != "user" && role != "admin" {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := s.repo.CreateAuthUser(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpsertProfile(ctx, userID, email, req.FullName, role, isActive); err != nil {
		return nil, err
	}

	return &AddUserResponse{UserID: userID}, nil
}

// UpdateUser updates profile fields and/or password.
func (s *Service) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	caller, err := s.callerFromContext(ctx)
	if err != nil {
		return err
	}

	if req.UserID == "" {
		return errors.New("user_id is required")
	}

	if req.Role != nil {
		if *req.Role == "super_admin" && caller.Role != "super_admin" {
			return errors.New("only super_admin can grant super_admin role")
		}
		if *req.Role != "user" && *req.Role != "admin" && *req.Role != "super_admin" {
			return ErrInvalidInput
		}
	}

	hasProfileUpdate := req.FullName != nil || req.Role != nil || req.IsActive != nil
	if hasProfileUpdate {
		if err := s.repo.UpdateProfile(ctx, req.UserID, req.FullName, req.Role, req.IsActive); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvalidInput
			}
			return err
		}
	}

	if req.NewPassword != nil {
		pw := *req.NewPassword
		if len(pw) < 6 {
			return errors.New("password must be at least 6 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := s.repo.UpdatePassword(ctx, req.UserID, string(hash)); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrInvalidInput
			}
			return err
		}
	}

	if !hasProfileUpdate && req.NewPassword == nil {
		return ErrInvalidInput
	}

	return nil
}

// DeleteUsers removes users from auth.users (profiles cascade).
func (s *Service) DeleteUsers(ctx context.Context, req DeleteUserRequest) (*DeleteUserResponse, error) {
	caller, err := s.callerFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.UserIDs) == 0 {
		return nil, errors.New("user_ids must be a non-empty array")
	}

	for _, id := range req.UserIDs {
		if id == caller.ID {
			return nil, ErrCannotDeleteSelf
		}
	}

	if err := s.repo.DeleteAuthUsers(ctx, req.UserIDs); err != nil {
		return nil, err
	}

	return &DeleteUserResponse{OK: true, Deleted: len(req.UserIDs)}, nil
}
