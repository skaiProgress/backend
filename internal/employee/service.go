package employee

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"aiqadam-backend/internal/ai"
	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
)

// AIRepository persists and runs AI analytics for training completion.
type AIRepository interface {
	ListPendingForEmployeeCourse(ctx context.Context, employeeID, courseID string) ([]ai.PendingItem, error)
	GetQuizResultForAnalysis(ctx context.Context, employeeID, quizResultID string) (*ai.EmployeeAnalysis, error)
	SaveAnalysis(ctx context.Context, employeeID, orgID, quizResultID, courseName string, score float64, a *ai.Analysis) error
}

// Service handles employee cabinet business logic.
type Service struct {
	repo   Repository
	aiRepo AIRepository
}

// NewService creates an employee service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// SetAIRepository injects AI analytics (after construction).
func (s *Service) SetAIRepository(ar AIRepository) {
	s.aiRepo = ar
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

// GetCourseProgress returns training completion status.
func (s *Service) GetCourseProgress(ctx context.Context, courseID string) (*CourseProgress, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	return s.repo.GetCourseProgress(ctx, userID, courseID)
}

// CompleteTraining marks course done and runs AI analysis for all quiz attempts.
func (s *Service) CompleteTraining(ctx context.Context, courseID string) (*CompleteTrainingResult, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}
	if s.aiRepo == nil {
		return nil, errors.New("AI module is not configured")
	}

	progress, err := s.repo.GetCourseProgress(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if progress == nil {
		return nil, pgx.ErrNoRows
	}
	if !progress.TrainingCompleted && !progress.CanComplete {
		return nil, errors.New("сдайте все тесты по урокам, чтобы завершить обучение")
	}

	analyzed := 0
	failed := 0
	var lastErr string

	pending, err := s.aiRepo.ListPendingForEmployeeCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	for _, item := range pending {
		data, err := s.aiRepo.GetQuizResultForAnalysis(ctx, item.EmployeeID, item.QuizResultID)
		if err != nil {
			failed++
			lastErr = err.Error()
			log.Printf("ai: get quiz result %s: %v", item.QuizResultID, err)
			continue
		}
		result, err := ai.Analyze(*data)
		if err != nil {
			failed++
			lastErr = err.Error()
			log.Printf("ai: openai analyze %s: %v", item.QuizResultID, err)
			continue
		}
		if err := s.aiRepo.SaveAnalysis(ctx,
			item.EmployeeID, item.OrganizationID, item.QuizResultID,
			data.CourseName, data.Score, result,
		); err != nil {
			failed++
			lastErr = err.Error()
			log.Printf("ai: save analysis %s: %v", item.QuizResultID, err)
			continue
		}
		analyzed++
	}

	if !progress.TrainingCompleted {
		if err := s.repo.MarkTrainingComplete(ctx, userID, courseID); err != nil {
			return nil, err
		}
	}

	out := &CompleteTrainingResult{
		TrainingCompleted: true,
		AnalyzedCount:     analyzed,
		FailedCount:       failed,
	}

	switch {
	case analyzed == 0 && failed > 0:
		out.Message = fmt.Sprintf("AI-анализ не выполнен (OpenAI): %s", lastErr)
	case analyzed == 0 && len(pending) == 0:
		out.Message = "Все результаты уже были проанализированы ранее"
	case failed > 0:
		out.Message = fmt.Sprintf("Частично: %d успешно, %d с ошибкой", analyzed, failed)
	default:
		out.Message = fmt.Sprintf("Проанализировано тестов: %d", analyzed)
	}

	return out, nil
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
