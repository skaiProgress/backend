package assignments

import (
	"context"
	"errors"
	"time"

	"aiqadam-backend/internal/auth"
)

// Service handles course assignments.
type Service struct {
	repo Repository
}

// NewService creates an assignments service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) requireAdmin(ctx context.Context) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return auth.ErrUnauthorized
	}
	if claims.Role != "admin" && claims.Role != "super_admin" {
		return auth.ErrForbidden
	}
	return nil
}

// List returns assignments for admin views.
func (s *Service) List(ctx context.Context, userID, courseID string, activeOnly bool) ([]ListItem, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, ListFilter{
		UserID:     userID,
		CourseID:   courseID,
		ActiveOnly: activeOnly,
	})
}

// Create assigns one course to one user.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Assignment, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if req.UserID == "" || req.CourseID == "" {
		return nil, errors.New("user_id and course_id are required")
	}
	expires, err := parseExpiresAt(req.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return s.repo.UpsertOne(ctx, Assignment{
		UserID:     req.UserID,
		CourseID:   req.CourseID,
		AssignedBy: req.AssignedBy,
		AssignedAt: time.Now().UTC(),
		ExpiresAt:  expires,
	})
}

// Bulk assigns courses to users (cartesian product).
func (s *Service) Bulk(ctx context.Context, req BulkRequest) (*BulkResponse, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 || len(req.CourseIDs) == 0 {
		return nil, errors.New("Выберите хотя бы одного пользователя и один курс")
	}
	expires, err := parseExpiresAt(req.ExpiresAt)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	rows := make([]Assignment, 0, len(req.UserIDs)*len(req.CourseIDs))
	for _, uid := range req.UserIDs {
		for _, cid := range req.CourseIDs {
			rows = append(rows, Assignment{
				UserID:     uid,
				CourseID:   cid,
				AssignedBy: req.AssignedBy,
				AssignedAt: now,
				ExpiresAt:  expires,
			})
		}
	}
	count, err := s.repo.UpsertBulk(ctx, rows)
	if err != nil {
		return nil, err
	}
	return &BulkResponse{Count: count}, nil
}

// Revoke soft-deletes an assignment.
func (s *Service) Revoke(ctx context.Context, id string) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	return s.repo.Revoke(ctx, id)
}
