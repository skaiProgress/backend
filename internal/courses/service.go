package courses

import (
	"context"
	"errors"
	"strings"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
)

// Service handles course business logic.
type Service struct {
	repo Repository
}

// NewService creates a courses service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) requireAdmin(ctx context.Context) (*auth.Claims, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, auth.ErrUnauthorized
	}
	if claims.Role != "admin" && claims.Role != "super_admin" {
		return nil, auth.ErrForbidden
	}
	return claims, nil
}

// List returns courses with lesson counts.
func (s *Service) List(ctx context.Context, search string) ([]CourseWithCount, error) {
	if _, err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, search)
}

// Create adds a new course.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Course, error) {
	claims, err := s.requireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}
	status := req.Status
	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "published" {
		return nil, errors.New("invalid status")
	}
	uid := claims.Subject
	return s.repo.Create(ctx, Course{
		Title:       title,
		Description: req.Description,
		Status:      status,
		CoverURL:    req.CoverURL,
		CreatedBy:   &uid,
	})
}

// Update patches a course.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Course, error) {
	if _, err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	fields := map[string]interface{}{}
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			return nil, errors.New("title cannot be empty")
		}
		fields["title"] = t
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Status != nil {
		if *req.Status != "draft" && *req.Status != "published" {
			return nil, errors.New("invalid status")
		}
		fields["status"] = *req.Status
	}
	if req.CoverURL != nil {
		fields["cover_url"] = *req.CoverURL
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

// Delete removes a course.
func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.requireAdmin(ctx); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
