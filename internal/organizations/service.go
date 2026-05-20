package organizations

import (
	"context"
	"errors"
	"strings"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service handles organization management (super_admin only).
type Service struct {
	repo Repository
}

// NewService creates an organizations service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) requireSuperAdmin(ctx context.Context) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return ErrUnauthorized
	}
	if claims.Role != "super_admin" {
		return ErrForbidden
	}
	return nil
}

// List returns all organizations.
func (s *Service) List(ctx context.Context, search string) ([]Organization, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, search)
}

// Get returns one organization with its users.
func (s *Service) Get(ctx context.Context, id string) (*OrganizationWithUsers, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, ErrInvalidInput
	}
	out, err := s.repo.GetWithUsers(ctx, id)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, pgx.ErrNoRows
	}
	return out, nil
}

// Create adds a new organization.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Organization, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	o := Organization{
		Name:          name,
		BIN:           trimOptional(req.BIN),
		Phone:         trimOptional(req.Phone),
		Email:         trimOptional(req.Email),
		Address:       trimOptional(req.Address),
		ContactPerson: trimOptional(req.ContactPerson),
		IsActive:      isActive,
	}

	return s.repo.Create(ctx, o)
}

// Update patches organization fields.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Organization, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, ErrInvalidInput
	}

	fields := make(map[string]interface{})
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		if n == "" {
			return nil, errors.New("name cannot be empty")
		}
		fields["name"] = n
	}
	if req.BIN != nil {
		fields["bin"] = trimOptional(req.BIN)
	}
	if req.Phone != nil {
		fields["phone"] = trimOptional(req.Phone)
	}
	if req.Email != nil {
		fields["email"] = trimOptional(req.Email)
	}
	if req.Address != nil {
		fields["address"] = trimOptional(req.Address)
	}
	if req.ContactPerson != nil {
		fields["contact_person"] = trimOptional(req.ContactPerson)
	}
	if req.IsActive != nil {
		fields["is_active"] = *req.IsActive
	}

	if len(fields) == 0 {
		return nil, ErrInvalidInput
	}

	out, err := s.repo.Update(ctx, id, fields)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, pgx.ErrNoRows
	}
	return out, nil
}

// Delete removes an organization (profiles.organization_id becomes NULL via FK).
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return err
	}
	if id == "" {
		return ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}

// AddMember creates a user bound to the organization with role user or org_admin.
func (s *Service) AddMember(ctx context.Context, orgID string, req AddMemberRequest) (*AddMemberResponse, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}

	org, err := s.repo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, pgx.ErrNoRows
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
	if role != "user" && role != "org_admin" {
		return nil, errors.New("role must be user or org_admin")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := s.repo.CreateAuthUser(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpsertOrgProfile(ctx, userID, email, orgID, req.FullName, role, isActive); err != nil {
		return nil, err
	}

	return &AddMemberResponse{UserID: userID}, nil
}

func trimOptional(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}
