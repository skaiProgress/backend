package employee

import (
	"context"
	"errors"
	"strings"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
)

// Service handles employee cabinet business logic.
type Service struct {
	repo Repository
}

// NewService creates an employee service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) userID(ctx context.Context) (string, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return "", auth.ErrUnauthorized
	}
	return claims.Subject, nil
}

// ListCourses returns assigned published courses.
func (s *Service) ListCourses(ctx context.Context) ([]MyCourse, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListMyCourses(ctx, userID)
}

// GetCourseDetail returns one assigned course.
func (s *Service) GetCourseDetail(ctx context.Context, courseID string) (*MyCourseDetail, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.GetMyCourseDetail(ctx, userID, courseID)
}

// ListLessons returns lessons when the user has assignment.
func (s *Service) ListLessons(ctx context.Context, courseID string) ([]MyLesson, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.ListMyLessons(ctx, userID, courseID)
}

// ListMaterials returns materials for an assigned course.
func (s *Service) ListMaterials(ctx context.Context, courseID string) ([]Material, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.ListMyMaterials(ctx, userID, courseID)
}

// GetProfile returns the current user's profile.
func (s *Service) GetProfile(ctx context.Context) (*Profile, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, pgx.ErrNoRows
	}
	if !p.IsActive {
		return nil, auth.ErrForbidden
	}
	return p, nil
}

// UpdateProfile updates full_name for the current user.
func (s *Service) UpdateProfile(ctx context.Context, fullName string) error {
	userID, err := s.userID(ctx)
	if err != nil {
		return err
	}
	fullName = strings.TrimSpace(fullName)
	return s.repo.UpdateFullName(ctx, userID, fullName)
}
