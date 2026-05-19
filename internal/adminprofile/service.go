package adminprofile

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/storage"

	"github.com/jackc/pgx/v5"
)

const maxAvatarBytes = 5 << 20 // 5 MiB

// Service handles admin profile settings.
type Service struct {
	repo  Repository
	files *storage.Local
}

// NewService creates an admin profile service.
func NewService(repo Repository, files *storage.Local) *Service {
	return &Service{repo: repo, files: files}
}

func (s *Service) requireAdmin(ctx context.Context) (string, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return "", auth.ErrUnauthorized
	}
	if claims.Role != "admin" && claims.Role != "super_admin" {
		return "", auth.ErrForbidden
	}
	return claims.Subject, nil
}

// Get returns the current admin profile.
func (s *Service) Get(ctx context.Context) (*Profile, error) {
	userID, err := s.requireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, pgx.ErrNoRows
	}
	if p.Email == "" {
		claims, _ := auth.ClaimsFromContext(ctx)
		if claims != nil {
			p.Email = claims.Email
		}
	}
	return p, nil
}

// Update patches allowed profile fields.
func (s *Service) Update(ctx context.Context, req UpdateRequest) (*Profile, error) {
	userID, err := s.requireAdmin(ctx)
	if err != nil {
		return nil, err
	}

	if req.FullName != nil {
		trimmed := strings.TrimSpace(*req.FullName)
		if len(trimmed) < 2 {
			return nil, errors.New("full_name must be at least 2 characters")
		}
		if len(trimmed) > 100 {
			return nil, errors.New("full_name is too long")
		}
		req.FullName = &trimmed
	}

	req.Phone = trimOptional(req.Phone, 20)
	req.Position = trimOptional(req.Position, 100)
	req.Department = trimOptional(req.Department, 100)
	req.Bio = trimOptional(req.Bio, 500)

	return s.repo.Update(ctx, userID, req)
}

// UploadAvatar stores an image and updates avatar_url.
func (s *Service) UploadAvatar(ctx context.Context, file *multipart.FileHeader) (*Profile, error) {
	userID, err := s.requireAdmin(ctx)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, errors.New("file is required")
	}
	if file.Size <= 0 || file.Size > maxAvatarBytes {
		return nil, errors.New("invalid file size")
	}
	ct := file.Header.Get("Content-Type")
	if ct != "" && !strings.HasPrefix(ct, "image/") {
		return nil, errors.New("file must be an image")
	}

	current, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	publicURL, _, err := s.files.SaveAvatar(userID, file.Filename, src)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.SetAvatarURL(ctx, userID, publicURL)
	if err != nil {
		_ = s.files.Delete(s.files.RelPathFromPublicURL(publicURL))
		return nil, err
	}

	if current != nil && current.AvatarURL != nil {
		if rel := s.files.RelPathFromPublicURL(*current.AvatarURL); rel != "" {
			_ = s.files.Delete(rel)
		}
	}

	return updated, nil
}

func trimOptional(s *string, maxLen int) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	if len(v) > maxLen {
		v = v[:maxLen]
	}
	return &v
}
