package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Local stores uploaded files on disk.
type Local struct {
	rootDir string
	baseURL string
}

// NewLocal creates a local file storage.
func NewLocal(rootDir, publicBaseURL string) (*Local, error) {
	if err := os.MkdirAll(filepath.Join(rootDir, "materials"), 0o755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(rootDir, "avatars"), 0o755); err != nil {
		return nil, fmt.Errorf("create avatars dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(rootDir, "videos"), 0o755); err != nil {
		return nil, fmt.Errorf("create videos dir: %w", err)
	}
	return &Local{
		rootDir: rootDir,
		baseURL: strings.TrimRight(publicBaseURL, "/"),
	}, nil
}

// SaveMaterial stores a course material file and returns public URL + relative path.
func (s *Local) SaveMaterial(courseID, originalName string, r io.Reader) (publicURL, relPath string, err error) {
	safeName := safeObjectName(originalName)
	relPath = filepath.ToSlash(filepath.Join("materials", courseID, safeName))
	absPath := filepath.Join(s.rootDir, relPath)

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", "", err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(absPath)
		return "", "", err
	}

	publicURL = fmt.Sprintf("%s/files/%s", s.baseURL, relPath)
	return publicURL, relPath, nil
}

// SaveLessonVideo stores a lesson video file and returns public URL + relative path.
func (s *Local) SaveLessonVideo(courseID, originalName string, r io.Reader) (publicURL, relPath string, err error) {
	safeName := safeObjectName(originalName)
	relPath = filepath.ToSlash(filepath.Join("videos", courseID, safeName))
	absPath := filepath.Join(s.rootDir, relPath)

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", "", err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(absPath)
		return "", "", err
	}

	publicURL = fmt.Sprintf("%s/files/%s", s.baseURL, relPath)
	return publicURL, relPath, nil
}

// SaveBriefingVideo stores a fire-safety briefing video for a course/kind.
func (s *Local) SaveBriefingVideo(courseID, kind, originalName string, r io.Reader) (publicURL, relPath string, err error) {
	safeName := safeObjectName(originalName)
	relPath = filepath.ToSlash(filepath.Join("videos", "briefings", courseID, kind, safeName))
	absPath := filepath.Join(s.rootDir, relPath)

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", "", err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(absPath)
		return "", "", err
	}

	publicURL = fmt.Sprintf("%s/files/%s", s.baseURL, relPath)
	return publicURL, relPath, nil
}

// SaveAvatar stores a profile avatar and returns public URL + relative path.
func (s *Local) SaveAvatar(userID, originalName string, r io.Reader) (publicURL, relPath string, err error) {
	safeName := safeObjectName(originalName)
	relPath = filepath.ToSlash(filepath.Join("avatars", userID, safeName))
	absPath := filepath.Join(s.rootDir, relPath)

	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return "", "", err
	}

	f, err := os.Create(absPath)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		_ = os.Remove(absPath)
		return "", "", err
	}

	publicURL = fmt.Sprintf("%s/files/%s", s.baseURL, relPath)
	return publicURL, relPath, nil
}

// Delete removes a file by relative path under storage root.
func (s *Local) Delete(relPath string) error {
	abs := filepath.Join(s.rootDir, filepath.FromSlash(relPath))
	if err := os.Remove(abs); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// RelPathFromPublicURL extracts storage relative path from a public file URL.
func (s *Local) RelPathFromPublicURL(fileURL string) string {
	marker := "/files/"
	pos := strings.Index(fileURL, marker)
	if pos == -1 {
		return ""
	}
	return fileURL[pos+len(marker):]
}

func safeObjectName(original string) string {
	ext := ""
	if i := strings.LastIndex(original, "."); i >= 0 && i < len(original)-1 {
		ext = original[i:]
		if len(ext) > 13 {
			ext = ""
		}
	}
	unique := uuid.New().String()
	return fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), strings.ReplaceAll(unique, "-", ""), ext)
}
