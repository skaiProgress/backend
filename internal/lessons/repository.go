package lessons

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists lessons.
type Repository interface {
	ListByCourse(ctx context.Context, courseID string) ([]Lesson, error)
	Create(ctx context.Context, l Lesson) (*Lesson, error)
	Update(ctx context.Context, id string, fields map[string]interface{}) (*Lesson, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, courseID string, orderedIDs []string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a lessons repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) ListByCourse(ctx context.Context, courseID string) ([]Lesson, error) {
	const q = `
		SELECT id::text, course_id::text, title, description, youtube_url, youtube_video_id,
		       order_index, is_free, created_at, updated_at
		FROM public.lessons
		WHERE course_id = $1::uuid
		ORDER BY order_index ASC, created_at ASC
	`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, fmt.Errorf("list lessons: %w", err)
	}
	defer rows.Close()

	out := make([]Lesson, 0)
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(
			&l.ID, &l.CourseID, &l.Title, &l.Description, &l.YoutubeURL, &l.YoutubeVideoID,
			&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *postgresRepository) Reorder(ctx context.Context, courseID string, orderedIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, id := range orderedIDs {
		tag, err := tx.Exec(ctx, `
			UPDATE public.lessons SET order_index = $1, updated_at = NOW()
			WHERE id = $2::uuid AND course_id = $3::uuid
		`, i+1, id, courseID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
	}
	return tx.Commit(ctx)
}

func (r *postgresRepository) Create(ctx context.Context, l Lesson) (*Lesson, error) {
	const q = `
		INSERT INTO public.lessons (
			course_id, title, description, youtube_url, youtube_video_id, order_index, is_free
		) VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, course_id::text, title, description, youtube_url, youtube_video_id,
		          order_index, is_free, created_at, updated_at
	`
	var out Lesson
	err := r.pool.QueryRow(ctx, q,
		l.CourseID, l.Title, l.Description, l.YoutubeURL, l.YoutubeVideoID, l.OrderIndex, l.IsFree,
	).Scan(
		&out.ID, &out.CourseID, &out.Title, &out.Description, &out.YoutubeURL, &out.YoutubeVideoID,
		&out.OrderIndex, &out.IsFree, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create lesson: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*Lesson, error) {
	setParts := make([]string, 0, 8)
	args := make([]interface{}, 0, 9)
	pos := 1
	for col, val := range fields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, pos))
		args = append(args, val)
		pos++
	}
	if len(setParts) == 0 {
		return r.getByID(ctx, id)
	}
	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE public.lessons SET %s WHERE id = $%d::uuid
		RETURNING id::text, course_id::text, title, description, youtube_url, youtube_video_id,
		          order_index, is_free, created_at, updated_at
	`, strings.Join(setParts, ", "), pos)

	var out Lesson
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&out.ID, &out.CourseID, &out.Title, &out.Description, &out.YoutubeURL, &out.YoutubeVideoID,
		&out.OrderIndex, &out.IsFree, &out.CreatedAt, &out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update lesson: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM public.lessons WHERE id = $1::uuid`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) getByID(ctx context.Context, id string) (*Lesson, error) {
	const q = `
		SELECT id::text, course_id::text, title, description, youtube_url, youtube_video_id,
		       order_index, is_free, created_at, updated_at
		FROM public.lessons WHERE id = $1::uuid
	`
	var out Lesson
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&out.ID, &out.CourseID, &out.Title, &out.Description, &out.YoutubeURL, &out.YoutubeVideoID,
		&out.OrderIndex, &out.IsFree, &out.CreatedAt, &out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
