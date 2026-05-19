package materials

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/storage"

	"github.com/jackc/pgx/v5"
)

const maxUploadBytes = 50 << 20 // 50 MiB

// Service handles course materials.
type Service struct {
	repo    Repository
	files   *storage.Local
}

// NewService creates a materials service.
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

// List returns materials for a course.
func (s *Service) List(ctx context.Context, courseID string) ([]Material, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.ListByCourse(ctx, courseID)
}

// Upload saves file to disk and inserts DB row.
func (s *Service) Upload(
	ctx context.Context,
	courseID, originalName, contentType string,
	size int64,
	body io.Reader,
) (*Material, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if courseID == "" || originalName == "" {
		return nil, errors.New("course_id and file are required")
	}
	if size <= 0 || size > maxUploadBytes {
		return nil, errors.New("invalid file size")
	}

	publicURL, _, err := s.files.SaveMaterial(courseID, originalName, body)
	if err != nil {
		return nil, err
	}

	var ft *string
	if contentType != "" {
		ft = &contentType
	}
	sz := size

	row, err := s.repo.Insert(ctx, Material{
		CourseID: courseID,
		Name:     originalName,
		FileURL:  publicURL,
		FileType: ft,
		FileSize: &sz,
	})
	if err != nil {
		_ = s.files.Delete(s.files.RelPathFromPublicURL(publicURL))
		return nil, err
	}
	return row, nil
}

// UploadFromMultipart handles multipart file upload.
func (s *Service) UploadFromMultipart(ctx context.Context, courseID string, fh *multipart.FileHeader) (*Material, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return s.Upload(ctx, courseID, fh.Filename, fh.Header.Get("Content-Type"), fh.Size, f)
}

// Delete removes DB row and stored file.
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
	if rel := s.files.RelPathFromPublicURL(row.FileURL); rel != "" {
		_ = s.files.Delete(rel)
	}
	return s.repo.Delete(ctx, id)
}
