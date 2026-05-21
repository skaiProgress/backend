package contactrequests

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"aiqadam-backend/internal/auth"
)

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Service handles contact form submissions.
type Service struct {
	repo Repository
}

// NewService creates a contact requests service.
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

func validateCreate(req CreateRequest) (CreateRequest, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	phone := strings.TrimSpace(req.Phone)
	message := trimOptional(req.Message)
	company := trimOptional(req.Company)

	if name == "" {
		return CreateRequest{}, errors.New("name is required")
	}
	if len(name) > 200 {
		return CreateRequest{}, errors.New("name is too long")
	}
	if email == "" {
		return CreateRequest{}, errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return CreateRequest{}, errors.New("invalid email")
	}
	if phone == "" {
		return CreateRequest{}, errors.New("phone is required")
	}
	if len(phone) > 30 {
		return CreateRequest{}, errors.New("phone is too long")
	}
	if message != nil && len(*message) > 5000 {
		return CreateRequest{}, errors.New("message is too long")
	}
	if company != nil && len(*company) > 200 {
		return CreateRequest{}, errors.New("company is too long")
	}

	return CreateRequest{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Company: company,
		Message: message,
	}, nil
}

func validateStatus(status string) (string, error) {
	switch status {
	case StatusNew, StatusInProgress, StatusDone:
		return status, nil
	default:
		return "", errors.New("invalid status")
	}
}

// Create stores a public contact form submission.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*ContactRequest, error) {
	validated, err := validateCreate(req)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, ContactRequest{
		Name:    validated.Name,
		Email:   validated.Email,
		Phone:   validated.Phone,
		Company: validated.Company,
		Message: validated.Message,
		Status:  StatusNew,
	})
}

// List returns contact requests for super admin.
func (s *Service) List(ctx context.Context, status string) ([]ContactRequest, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	if status != "" {
		if _, err := validateStatus(status); err != nil {
			return nil, err
		}
	}
	return s.repo.List(ctx, status)
}

// Get returns one contact request for super admin.
func (s *Service) Get(ctx context.Context, id string) (*ContactRequest, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	if id == "" {
		return nil, ErrInvalidInput
	}
	out, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateStatus changes request status for super admin.
func (s *Service) UpdateStatus(ctx context.Context, id string, req UpdateRequest) (*ContactRequest, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return nil, err
	}
	if id == "" || req.Status == nil {
		return nil, ErrInvalidInput
	}
	status, err := validateStatus(strings.TrimSpace(*req.Status))
	if err != nil {
		return nil, err
	}
	out, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CountNew returns the number of new requests for super admin dashboard badge.
func (s *Service) CountNew(ctx context.Context) (int, error) {
	if err := s.requireSuperAdmin(ctx); err != nil {
		return 0, err
	}
	return s.repo.CountByStatus(ctx, StatusNew)
}

func trimOptional(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}
