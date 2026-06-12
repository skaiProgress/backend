package courses

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists courses.
type Repository interface {
	List(ctx context.Context, search string) ([]CourseWithCount, error)
	GetByID(ctx context.Context, id string) (*Course, error)
	Create(ctx context.Context, c Course) (*Course, error)
	Update(ctx context.Context, id string, fields map[string]interface{}) (*Course, error)
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a courses repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) List(ctx context.Context, search string) ([]CourseWithCount, error) {
	const q = `
		SELECT c.id::text, c.title, c.description, c.status, c.cover_url, c.created_by::text,
		       c.is_briefing_course, c.created_at, c.updated_at,
		       COALESCE(COUNT(l.id), 0)::int AS lesson_count
		FROM public.courses c
		LEFT JOIN public.lessons l ON l.course_id = c.id
		WHERE ($1 = '' OR c.title ILIKE '%' || $1 || '%')
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, strings.TrimSpace(search))
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	out := make([]CourseWithCount, 0)
	for rows.Next() {
		var item CourseWithCount
		if err := rows.Scan(
			&item.ID, &item.Title, &item.Description, &item.Status, &item.CoverURL,
			&item.CreatedBy, &item.IsBriefingCourse, &item.CreatedAt, &item.UpdatedAt, &item.LessonCount,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Course, error) {
	const q = `
		SELECT id::text, title, description, status, cover_url, created_by::text, is_briefing_course, created_at, updated_at
		FROM public.courses WHERE id = $1::uuid
	`
	var c Course
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&c.ID, &c.Title, &c.Description, &c.Status, &c.CoverURL,
		&c.CreatedBy, &c.IsBriefingCourse, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course: %w", err)
	}
	return &c, nil
}

func (r *postgresRepository) Create(ctx context.Context, c Course) (*Course, error) {
	const q = `
		INSERT INTO public.courses (title, description, status, cover_url, created_by, is_briefing_course)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, $6)
		RETURNING id::text, title, description, status, cover_url, created_by::text, is_briefing_course, created_at, updated_at
	`
	createdBy := ""
	if c.CreatedBy != nil {
		createdBy = *c.CreatedBy
	}
	var out Course
	err := r.pool.QueryRow(ctx, q, c.Title, c.Description, c.Status, c.CoverURL, createdBy, c.IsBriefingCourse).Scan(
		&out.ID, &out.Title, &out.Description, &out.Status, &out.CoverURL,
		&out.CreatedBy, &out.IsBriefingCourse, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create course: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*Course, error) {
	setParts := make([]string, 0, 5)
	args := make([]interface{}, 0, 6)
	pos := 1
	for col, val := range fields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, pos))
		args = append(args, val)
		pos++
	}
	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}
	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE public.courses SET %s WHERE id = $%d::uuid
		RETURNING id::text, title, description, status, cover_url, created_by::text, is_briefing_course, created_at, updated_at
	`, strings.Join(setParts, ", "), pos)

	var out Course
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&out.ID, &out.Title, &out.Description, &out.Status, &out.CoverURL,
		&out.CreatedBy, &out.IsBriefingCourse, &out.CreatedAt, &out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update course: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM public.courses WHERE id = $1::uuid`, id)
	if err != nil {
		return fmt.Errorf("delete course: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
