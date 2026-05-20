package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Material is a course material row for employees.
type Material struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Name      string    `json:"name"`
	FileURL   string    `json:"file_url"`
	FileType  *string   `json:"file_type"`
	FileSize  *int64    `json:"file_size"`
	CreatedAt string    `json:"created_at"`
}

// Repository loads employee-scoped data.
type Repository interface {
	ListMyCourses(ctx context.Context, userID string) ([]MyCourse, error)
	GetMyCourseDetail(ctx context.Context, userID, courseID string) (*MyCourseDetail, error)
	ListMyLessons(ctx context.Context, userID, courseID string) ([]MyLesson, error)
	ListMyMaterials(ctx context.Context, userID, courseID string) ([]Material, error)
	GetProfile(ctx context.Context, userID string) (*Profile, error)
	UpdateFullName(ctx context.Context, userID, fullName string) error
	GetCourseProgress(ctx context.Context, userID, courseID string) (*CourseProgress, error)
	MarkTrainingComplete(ctx context.Context, userID, courseID string) error

}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates an employee repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) ListMyCourses(ctx context.Context, userID string) ([]MyCourse, error) {
	const q = `
		SELECT ca.id::text, c.id::text, c.title, c.description, c.cover_url,
		       ca.assigned_at, ca.expires_at,
		       (ca.expires_at IS NOT NULL AND ca.expires_at < NOW())
		FROM public.course_assignments ca
		JOIN public.courses c ON c.id = ca.course_id
		WHERE ca.user_id = $1::uuid
		  AND ca.status = 'active'
		  AND c.status = 'published'
		ORDER BY ca.assigned_at DESC
	`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list my courses: %w", err)
	}
	defer rows.Close()

	out := make([]MyCourse, 0)
	for rows.Next() {
		var c MyCourse
		if err := rows.Scan(
			&c.AssignmentID, &c.CourseID, &c.Title, &c.Description, &c.CoverURL,
			&c.AssignedAt, &c.ExpiresAt, &c.IsExpired,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetMyCourseDetail(ctx context.Context, userID, courseID string) (*MyCourseDetail, error) {
	const q = `
		SELECT c.id::text, c.title, c.description, c.cover_url,
		       ca.assigned_at, ca.expires_at,
		       (ca.expires_at IS NOT NULL AND ca.expires_at < NOW())
		FROM public.course_assignments ca
		JOIN public.courses c ON c.id = ca.course_id
		WHERE ca.user_id = $1::uuid
		  AND ca.course_id = $2::uuid
		  AND ca.status = 'active'
		  AND c.status = 'published'
		LIMIT 1
	`
	var d MyCourseDetail
	err := r.pool.QueryRow(ctx, q, userID, courseID).Scan(
		&d.CourseID, &d.Title, &d.Description, &d.CoverURL,
		&d.AssignedAt, &d.ExpiresAt, &d.IsExpired,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get my course detail: %w", err)
	}
	return &d, nil
}

func (r *postgresRepository) hasAssignment(ctx context.Context, userID, courseID string, requireNotExpired bool) (bool, error) {
	q := `
		SELECT EXISTS (
			SELECT 1 FROM public.course_assignments ca
			JOIN public.courses c ON c.id = ca.course_id
			WHERE ca.user_id = $1::uuid
			  AND ca.course_id = $2::uuid
			  AND ca.status = 'active'
			  AND c.status = 'published'
	`
	if requireNotExpired {
		q += ` AND (ca.expires_at IS NULL OR ca.expires_at > NOW())`
	}
	q += ` )`
	var ok bool
	err := r.pool.QueryRow(ctx, q, userID, courseID).Scan(&ok)
	return ok, err
}

func (r *postgresRepository) ListMyLessons(ctx context.Context, userID, courseID string) ([]MyLesson, error) {
	ok, err := r.hasAssignment(ctx, userID, courseID, false)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []MyLesson{}, nil
	}

	const q = `
		SELECT l.id::text, l.title, l.description, l.video_source,
		       l.youtube_url, l.youtube_video_id, l.video_url,
		       l.order_index, l.is_free,
		       EXISTS (SELECT 1 FROM public.lesson_quizzes qz WHERE qz.lesson_id = l.id) AS has_quiz,
		       COALESCE((
		           SELECT a.passed
		           FROM public.lesson_quiz_attempts a
		           WHERE a.user_id = $2::uuid AND a.lesson_id = l.id
		           ORDER BY a.completed_at DESC
		           LIMIT 1
		       ), FALSE) AS quiz_passed,
		       EXISTS (
		           SELECT 1 FROM public.lesson_quiz_attempts a
		           WHERE a.user_id = $2::uuid AND a.lesson_id = l.id
		       ) AS quiz_submitted
		FROM public.lessons l
		WHERE l.course_id = $1::uuid
		ORDER BY l.order_index ASC, l.created_at ASC
	`
	rows, err := r.pool.Query(ctx, q, courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("list my lessons: %w", err)
	}
	defer rows.Close()

	out := make([]MyLesson, 0)
	for rows.Next() {
		var l MyLesson
		if err := rows.Scan(
			&l.ID, &l.Title, &l.Description, &l.VideoSource,
			&l.YoutubeURL, &l.YoutubeVideoID, &l.VideoURL,
			&l.OrderIndex, &l.IsFree, &l.HasQuiz, &l.QuizPassed, &l.QuizSubmitted,
		); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *postgresRepository) ListMyMaterials(ctx context.Context, userID, courseID string) ([]Material, error) {
	ok, err := r.hasAssignment(ctx, userID, courseID, true)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []Material{}, nil
	}

	const q = `
		SELECT id::text, course_id::text, name, file_url, file_type, file_size, created_at::text
		FROM public.course_materials
		WHERE course_id = $1::uuid
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Material, 0)
	for rows.Next() {
		var m Material
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Name, &m.FileURL, &m.FileType, &m.FileSize, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	const q = `
		SELECT id::text, COALESCE(email, ''), full_name, role, is_active
		FROM public.profiles
		WHERE id = $1::uuid
	`
	var p Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(&p.ID, &p.Email, &p.FullName, &p.Role, &p.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postgresRepository) GetCourseProgress(ctx context.Context, userID, courseID string) (*CourseProgress, error) {
	ok, err := r.hasAssignment(ctx, userID, courseID, true)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	const q = `
		SELECT
			(SELECT COUNT(*)::int FROM public.lessons l
			 WHERE l.course_id = $2::uuid
			   AND EXISTS (SELECT 1 FROM public.lesson_quizzes qz WHERE qz.lesson_id = l.id)),
			(SELECT COUNT(DISTINCT l.id)::int FROM public.lessons l
			 JOIN public.lesson_quizzes qz ON qz.lesson_id = l.id
			 JOIN public.lesson_quiz_attempts a ON a.lesson_id = l.id AND a.user_id = $1::uuid
			 WHERE l.course_id = $2::uuid),
			(SELECT COUNT(DISTINCT l.id)::int FROM public.lessons l
			 JOIN public.lesson_quizzes qz ON qz.lesson_id = l.id
			 JOIN public.lesson_quiz_attempts a ON a.lesson_id = l.id AND a.user_id = $1::uuid AND a.passed = TRUE
			 WHERE l.course_id = $2::uuid),
			ca.training_completed_at::text
		FROM public.course_assignments ca
		WHERE ca.user_id = $1::uuid
		  AND ca.course_id = $2::uuid
		  AND ca.status = 'active'
	`
	var totalQuiz, submitted, passed int
	var completedAt *string
	err = r.pool.QueryRow(ctx, q, userID, courseID).Scan(&totalQuiz, &submitted, &passed, &completedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	canComplete := totalQuiz > 0 && submitted >= totalQuiz
	trainingDone := completedAt != nil && *completedAt != ""

	return &CourseProgress{
		TotalQuizLessons:     totalQuiz,
		SubmittedQuizLessons: submitted,
		PassedQuizLessons:    passed,
		CanComplete:          canComplete,
		TrainingCompleted:    trainingDone,
		CompletedAt:          completedAt,
	}, nil
}

func (r *postgresRepository) MarkTrainingComplete(ctx context.Context, userID, courseID string) error {
	const q = `
		UPDATE public.course_assignments
		SET training_completed_at = NOW(), updated_at = NOW()
		WHERE user_id = $1::uuid
		  AND course_id = $2::uuid
		  AND status = 'active'
		  AND training_completed_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, q, userID, courseID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// already completed — not an error
		const q2 = `
			SELECT 1 FROM public.course_assignments
			WHERE user_id = $1::uuid AND course_id = $2::uuid
			  AND status = 'active' AND training_completed_at IS NOT NULL
		`
		var one int
		if err := r.pool.QueryRow(ctx, q2, userID, courseID).Scan(&one); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return pgx.ErrNoRows
			}
			return err
		}
	}
	return nil
}

func (r *postgresRepository) UpdateFullName(ctx context.Context, userID, fullName string) error {
	const q = `
		UPDATE public.profiles
		SET full_name = $2, updated_at = NOW()
		WHERE id = $1::uuid
	`
	tag, err := r.pool.Exec(ctx, q, userID, fullName)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
