package lessons

import (
	"context"
	"errors"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/youtube"

	"github.com/jackc/pgx/v5"
)

// Service handles lesson business logic.
type Service struct {
	repo Repository
}

// NewService creates a lessons service.
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

// List returns lessons for a course.
func (s *Service) List(ctx context.Context, courseID string) ([]Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.ListByCourse(ctx, courseID)
}

// Reorder updates order_index for lessons in a course.
func (s *Service) Reorder(ctx context.Context, req ReorderRequest) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	if req.CourseID == "" || len(req.OrderedIDs) == 0 {
		return errors.New("course_id and ordered_ids are required")
	}
	return s.repo.Reorder(ctx, req.CourseID, req.OrderedIDs)
}

// Create adds a lesson with YouTube metadata.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if req.CourseID == "" || title == "" || strings.TrimSpace(req.YoutubeURL) == "" {
		return nil, errors.New("course_id, title and youtube_url are required")
	}
	videoID := req.YoutubeVideoID
	if videoID == "" {
		videoID = youtube.ParseVideoID(req.YoutubeURL)
	}
	if videoID == "" {
		return nil, errors.New("invalid youtube url")
	}
	order := req.OrderIndex
	if order <= 0 {
		order = 1
	}
	return s.repo.Create(ctx, Lesson{
		CourseID:       req.CourseID,
		Title:          title,
		Description:    req.Description,
		YoutubeURL:     strings.TrimSpace(req.YoutubeURL),
		YoutubeVideoID: videoID,
		OrderIndex:     order,
		IsFree:         req.IsFree,
	})
}

// Update patches a lesson.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
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
	if req.YoutubeURL != nil {
		url := strings.TrimSpace(*req.YoutubeURL)
		fields["youtube_url"] = url
		vid := ""
		if req.YoutubeVideoID != nil && *req.YoutubeVideoID != "" {
			vid = *req.YoutubeVideoID
		} else {
			vid = youtube.ParseVideoID(url)
		}
		if vid == "" {
			return nil, errors.New("invalid youtube url")
		}
		fields["youtube_video_id"] = vid
	} else if req.YoutubeVideoID != nil {
		fields["youtube_video_id"] = *req.YoutubeVideoID
	}
	if req.OrderIndex != nil {
		fields["order_index"] = *req.OrderIndex
	}
	if req.IsFree != nil {
		fields["is_free"] = *req.IsFree
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

// Delete removes a lesson.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
