package quizzes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists lesson quizzes.
type Repository interface {
	LessonBelongsToCourse(ctx context.Context, lessonID string) (courseID string, err error)
	UpsertQuiz(ctx context.Context, lessonID, fileName string, questions []QuestionRow) (*QuizMeta, error)
	DeleteByLesson(ctx context.Context, lessonID string) error
	GetAdminQuiz(ctx context.Context, lessonID string) (*AdminQuiz, error)
	GetEmployeeQuiz(ctx context.Context, lessonID, userID string) (*EmployeeQuizPayload, error)
	SaveAttempt(ctx context.Context, userID, lessonID, quizID string, answers json.RawMessage, score int, passed bool) (*AttemptSummary, error)
	HasAssignment(ctx context.Context, userID, courseID string) (bool, error)
}

// QuestionRow is a question insert row.
type QuestionRow struct {
	OrderIndex    int
	QuestionText  string
	OptionA       string
	OptionB       string
	OptionC       string
	CorrectOption string
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a quizzes repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) LessonBelongsToCourse(ctx context.Context, lessonID string) (string, error) {
	const q = `SELECT course_id::text FROM public.lessons WHERE id = $1::uuid`
	var courseID string
	err := r.pool.QueryRow(ctx, q, lessonID).Scan(&courseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", pgx.ErrNoRows
	}
	return courseID, err
}

func (r *postgresRepository) HasAssignment(ctx context.Context, userID, courseID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.course_assignments ca
			JOIN public.courses c ON c.id = ca.course_id
			WHERE ca.user_id = $1::uuid
			  AND ca.course_id = $2::uuid
			  AND ca.status = 'active'
			  AND c.status = 'published'
			  AND (ca.expires_at IS NULL OR ca.expires_at > NOW())
		)
	`
	var ok bool
	err := r.pool.QueryRow(ctx, q, userID, courseID).Scan(&ok)
	return ok, err
}

func (r *postgresRepository) UpsertQuiz(
	ctx context.Context,
	lessonID, fileName string,
	questions []QuestionRow,
) (*QuizMeta, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var quizID string
	err = tx.QueryRow(ctx, `
		INSERT INTO public.lesson_quizzes (lesson_id, source_file_name)
		VALUES ($1::uuid, $2)
		ON CONFLICT (lesson_id) DO UPDATE
		SET source_file_name = EXCLUDED.source_file_name,
		    updated_at = NOW()
		RETURNING id::text
	`, lessonID, fileName).Scan(&quizID)
	if err != nil {
		return nil, fmt.Errorf("upsert quiz: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM public.lesson_quiz_questions WHERE quiz_id = $1::uuid`, quizID); err != nil {
		return nil, err
	}

	for _, qn := range questions {
		_, err := tx.Exec(ctx, `
			INSERT INTO public.lesson_quiz_questions (
				quiz_id, order_index, question_text,
				option_a, option_b, option_c, correct_option
			) VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)
		`, quizID, qn.OrderIndex, qn.QuestionText, qn.OptionA, qn.OptionB, qn.OptionC, qn.CorrectOption)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.getMeta(ctx, lessonID)
}

func (r *postgresRepository) getMeta(ctx context.Context, lessonID string) (*QuizMeta, error) {
	const q = `
		SELECT q.id::text, q.lesson_id::text, q.source_file_name,
		       COUNT(qq.id)::int, q.created_at, q.updated_at
		FROM public.lesson_quizzes q
		LEFT JOIN public.lesson_quiz_questions qq ON qq.quiz_id = q.id
		WHERE q.lesson_id = $1::uuid
		GROUP BY q.id
	`
	var meta QuizMeta
	err := r.pool.QueryRow(ctx, q, lessonID).Scan(
		&meta.ID, &meta.LessonID, &meta.SourceFileName,
		&meta.QuestionCount, &meta.CreatedAt, &meta.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &meta, nil
}

func (r *postgresRepository) DeleteByLesson(ctx context.Context, lessonID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM public.lesson_quizzes WHERE lesson_id = $1::uuid`, lessonID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) GetAdminQuiz(ctx context.Context, lessonID string) (*AdminQuiz, error) {
	meta, err := r.getMeta(ctx, lessonID)
	if err != nil || meta == nil {
		return nil, err
	}

	const q = `
		SELECT id::text, order_index, question_text,
		       option_a, option_b, option_c, correct_option
		FROM public.lesson_quiz_questions
		WHERE quiz_id = $1::uuid
		ORDER BY order_index ASC
	`
	rows, err := r.pool.Query(ctx, q, meta.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]AdminQuestion, 0)
	for rows.Next() {
		var item AdminQuestion
		if err := rows.Scan(
			&item.ID, &item.OrderIndex, &item.QuestionText,
			&item.OptionA, &item.OptionB, &item.OptionC, &item.CorrectOption,
		); err != nil {
			return nil, err
		}
		questions = append(questions, item)
	}
	return &AdminQuiz{QuizMeta: *meta, Questions: questions}, rows.Err()
}

func (r *postgresRepository) GetEmployeeQuiz(ctx context.Context, lessonID, userID string) (*EmployeeQuizPayload, error) {
	admin, err := r.GetAdminQuiz(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, nil
	}

	questions := make([]EmployeeQuestion, 0, len(admin.Questions))
	for _, qn := range admin.Questions {
		questions = append(questions, EmployeeQuestion{
			ID:           qn.ID,
			OrderIndex:   qn.OrderIndex,
			QuestionText: qn.QuestionText,
			Options:      []string{qn.OptionA, qn.OptionB, qn.OptionC},
		})
	}

	attempt, err := r.latestAttempt(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}

	var attemptView *EmployeeAttemptView
	if attempt != nil {
		attemptView = &EmployeeAttemptView{
			Submitted:   true,
			Passed:      attempt.Passed,
			CompletedAt: attempt.CompletedAt,
		}
	}

	return &EmployeeQuizPayload{
		LessonID:  lessonID,
		Questions: questions,
		Attempt:   attemptView,
	}, nil
}

func (r *postgresRepository) latestAttempt(ctx context.Context, userID, lessonID string) (*AttemptSummary, error) {
	const q = `
		SELECT score, max_score, passed, completed_at
		FROM public.lesson_quiz_attempts
		WHERE user_id = $1::uuid AND lesson_id = $2::uuid
		ORDER BY completed_at DESC
		LIMIT 1
	`
	var a AttemptSummary
	err := r.pool.QueryRow(ctx, q, userID, lessonID).Scan(&a.Score, &a.MaxScore, &a.Passed, &a.CompletedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *postgresRepository) SaveAttempt(
	ctx context.Context,
	userID, lessonID, quizID string,
	answers json.RawMessage,
	score int,
	passed bool,
) (*AttemptSummary, error) {
	const q = `
		INSERT INTO public.lesson_quiz_attempts (
			user_id, lesson_id, quiz_id, answers, score, max_score, passed
		) VALUES ($1::uuid, $2::uuid, $3::uuid, $4::jsonb, $5, 5, $6)
		RETURNING score, max_score, passed, completed_at
	`
	var a AttemptSummary
	err := r.pool.QueryRow(ctx, q, userID, lessonID, quizID, answers, score, passed).Scan(
		&a.Score, &a.MaxScore, &a.Passed, &a.CompletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
