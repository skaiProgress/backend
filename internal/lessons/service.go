package lessons

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/youtube"
	"aiqadam-backend/internal/storage"

	"github.com/jackc/pgx/v5"
)

const maxVideoUploadBytes = 1 << 30 // 1 GiB

var allowedVideoExtensions = map[string]struct{}{
	".mp4":  {},
	".webm": {},
	".mov":  {},
	".m4v":  {},
}

// Service handles lesson business logic.
type Service struct {
	repo  Repository
	files *storage.Local
}

// NewService creates a lessons service.
func NewService(repo Repository, files *storage.Local) *Service {
	return &Service{repo: repo, files: files}
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
	url := strings.TrimSpace(req.YoutubeURL)
	return s.repo.Create(ctx, Lesson{
		CourseID:       req.CourseID,
		Title:          title,
		Description:    req.Description,
		VideoSource:    VideoSourceYouTube,
		YoutubeURL:     &url,
		YoutubeVideoID: &videoID,
		OrderIndex:     order,
		IsFree:         req.IsFree,
	})
}

// CreateFromMultipart adds a lesson with an uploaded video file.
func (s *Service) CreateFromMultipart(ctx context.Context, req CreateUploadRequest, fh *multipart.FileHeader) (*Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if req.CourseID == "" || title == "" {
		return nil, errors.New("course_id, title and video file are required")
	}
	if fh == nil {
		return nil, errors.New("video file is required")
	}
	if err := validateVideoFile(fh.Filename, fh.Size); err != nil {
		return nil, err
	}

	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	publicURL, _, err := s.files.SaveLessonVideo(req.CourseID, fh.Filename, f)
	if err != nil {
		return nil, err
	}

	order := req.OrderIndex
	if order <= 0 {
		order = 1
	}

	row, err := s.repo.Create(ctx, Lesson{
		CourseID:    req.CourseID,
		Title:       title,
		Description: req.Description,
		VideoSource: VideoSourceUpload,
		VideoURL:    &publicURL,
		OrderIndex:  order,
		IsFree:      req.IsFree,
	})
	if err != nil {
		_ = s.files.Delete(s.files.RelPathFromPublicURL(publicURL))
		return nil, err
	}
	return row, nil
}

// Update patches a YouTube lesson.
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, pgx.ErrNoRows
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
		fields["video_source"] = VideoSourceYouTube
		fields["video_url"] = nil
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

	if req.YoutubeURL != nil &&
		current.VideoSource == VideoSourceUpload &&
		current.VideoURL != nil {
		if rel := s.files.RelPathFromPublicURL(*current.VideoURL); rel != "" {
			_ = s.files.Delete(rel)
		}
	}
	return out, nil
}

// UpdateFromMultipart patches lesson metadata and optionally replaces uploaded video.
func (s *Service) UpdateFromMultipart(
	ctx context.Context,
	id string,
	req UpdateUploadRequest,
	fh *multipart.FileHeader,
) (*Lesson, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, pgx.ErrNoRows
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
	if req.OrderIndex != nil {
		fields["order_index"] = *req.OrderIndex
	}
	if req.IsFree != nil {
		fields["is_free"] = *req.IsFree
	}

	var oldVideoURL string
	if fh != nil {
		if err := validateVideoFile(fh.Filename, fh.Size); err != nil {
			return nil, err
		}
		if current.VideoURL != nil {
			oldVideoURL = *current.VideoURL
		}

		file, err := fh.Open()
		if err != nil {
			return nil, err
		}
		publicURL, _, err := s.files.SaveLessonVideo(current.CourseID, fh.Filename, file)
		file.Close()
		if err != nil {
			return nil, err
		}

		fields["video_source"] = VideoSourceUpload
		fields["video_url"] = publicURL
		fields["youtube_url"] = nil
		fields["youtube_video_id"] = nil
	}

	if len(fields) == 0 {
		return current, nil
	}

	out, err := s.repo.Update(ctx, id, fields)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, pgx.ErrNoRows
	}

	if oldVideoURL != "" && out.VideoURL != nil && oldVideoURL != *out.VideoURL {
		if rel := s.files.RelPathFromPublicURL(oldVideoURL); rel != "" {
			_ = s.files.Delete(rel)
		}
	}
	return out, nil
}

// Delete removes a lesson and its uploaded video file if present.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if row == nil {
		return pgx.ErrNoRows
	}
	if row.VideoSource == VideoSourceUpload && row.VideoURL != nil {
		if rel := s.files.RelPathFromPublicURL(*row.VideoURL); rel != "" {
			_ = s.files.Delete(rel)
		}
	}
	return s.repo.Delete(ctx, id)
}

func validateVideoFile(filename string, size int64) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if _, ok := allowedVideoExtensions[ext]; !ok {
		return errors.New("unsupported video format (allowed: mp4, webm, mov, m4v)")
	}
	if size <= 0 || size > maxVideoUploadBytes {
		return errors.New("invalid video file size")
	}
	return nil
}
