package quizzes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/quiztxt"

	"github.com/jackc/pgx/v5"
)

const passThreshold = 4 // 4 из 5

// Service handles quiz business logic.
type Service struct {
	repo Repository
}

// NewService creates a quizzes service.
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

func (s *Service) userID(ctx context.Context) (string, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return "", auth.ErrUnauthorized
	}
	return claims.Subject, nil
}

// GetAdmin returns quiz for a lesson (admin).
func (s *Service) GetAdmin(ctx context.Context, lessonID string) (*AdminQuiz, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if lessonID == "" {
		return nil, errors.New("lesson_id is required")
	}
	return s.repo.GetAdminQuiz(ctx, lessonID)
}

// Upload parses .txt and saves quiz for a lesson.
func (s *Service) Upload(ctx context.Context, lessonID, fileName string, body io.Reader) (*AdminQuiz, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	if lessonID == "" || fileName == "" {
		return nil, errors.New("lesson_id and file are required")
	}

	raw, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	parsed, err := quiztxt.Parse(string(raw))
	if err != nil {
		return nil, err
	}

	rows := make([]QuestionRow, 0, len(parsed))
	for _, q := range parsed {
		rows = append(rows, QuestionRow{
			OrderIndex:    q.OrderIndex,
			QuestionText:  q.QuestionText,
			OptionA:       q.OptionA,
			OptionB:       q.OptionB,
			OptionC:       q.OptionC,
			CorrectOption: q.CorrectOption,
		})
	}

	if _, err := s.repo.UpsertQuiz(ctx, lessonID, fileName, rows); err != nil {
		return nil, err
	}
	return s.repo.GetAdminQuiz(ctx, lessonID)
}

// Delete removes quiz for a lesson.
func (s *Service) Delete(ctx context.Context, lessonID string) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	return s.repo.DeleteByLesson(ctx, lessonID)
}

// GetEmployee returns quiz without answers for assigned employee.
func (s *Service) GetEmployee(ctx context.Context, lessonID string) (*EmployeeQuizPayload, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if lessonID == "" {
		return nil, errors.New("lesson_id is required")
	}

	courseID, err := s.repo.LessonBelongsToCourse(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	ok, err := s.repo.HasAssignment(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return s.repo.GetEmployeeQuiz(ctx, lessonID, userID)
}

// Submit grades employee answers (one attempt per lesson, no score in response).
func (s *Service) Submit(ctx context.Context, lessonID string, req SubmitRequest) (*EmployeeSubmitResult, error) {
	userID, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	if lessonID == "" {
		return nil, errors.New("lesson_id is required")
	}
	if len(req.Answers) != quiztxt.RequiredQuestions {
		return nil, errors.New("нужно ответить на все 5 вопросов")
	}

	courseID, err := s.repo.LessonBelongsToCourse(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	ok, err := s.repo.HasAssignment(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, pgx.ErrNoRows
	}

	existing, err := s.repo.GetEmployeeQuiz(ctx, lessonID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Attempt != nil && existing.Attempt.Submitted {
		return nil, errors.New("повторная сдача теста недоступна")
	}

	admin, err := s.repo.GetAdminQuiz(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, pgx.ErrNoRows
	}

	correctByID := map[string]string{}
	for _, q := range admin.Questions {
		correctByID[q.ID] = q.CorrectOption
	}

	score := 0
	for _, ans := range req.Answers {
		want, ok := correctByID[ans.QuestionID]
		if !ok {
			return nil, errors.New("неизвестный question_id")
		}
		given := strings.ToUpper(strings.TrimSpace(ans.Answer))
		if given == want {
			score++
		}
	}

	passed := score >= passThreshold
	payload, err := json.Marshal(req.Answers)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.SaveAttempt(ctx, userID, lessonID, admin.ID, payload, score, passed); err != nil {
		return nil, err
	}

	return &EmployeeSubmitResult{
		Submitted: true,
		Passed:    passed,
	}, nil
}

// HasQuiz checks if lesson has a quiz (for list enrichment).
func (s *Service) HasQuiz(ctx context.Context, lessonID string) (bool, error) {
	admin, err := s.repo.GetAdminQuiz(ctx, lessonID)
	if err != nil {
		return false, err
	}
	return admin != nil, nil
}
