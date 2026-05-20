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
	GetByID(ctx context.Context, id string) (*Lesson, error)
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

const lessonSelectCols = `
	id::text, course_id::text, title, description, video_source,
	youtube_url, youtube_video_id, video_url,
	order_index, is_free, created_at, updated_at
`

func scanLesson(row pgx.Row) (Lesson, error) {
	var l Lesson
	err := row.Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Description, &l.VideoSource,
		&l.YoutubeURL, &l.YoutubeVideoID, &l.VideoURL,
		&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
	)
	return l, err
}

func (r *postgresRepository) ListByCourse(ctx context.Context, courseID string) ([]Lesson, error) {
	q := fmt.Sprintf(`
		SELECT %s
		FROM public.lessons
		WHERE course_id = $1::uuid
		ORDER BY order_index ASC, created_at ASC
	`, lessonSelectCols)
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, fmt.Errorf("list lessons: %w", err)
	}
	defer rows.Close()

	out := make([]Lesson, 0)
	for rows.Next() {
		l, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Lesson, error) {
	q := fmt.Sprintf(`SELECT %s FROM public.lessons WHERE id = $1::uuid`, lessonSelectCols)
	l, err := scanLesson(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
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
			course_id, title, description, video_source,
			youtube_url, youtube_video_id, video_url,
			order_index, is_free
		) VALUES ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING ` + lessonSelectCols
	row := r.pool.QueryRow(ctx, q,
		l.CourseID, l.Title, l.Description, l.VideoSource,
		l.YoutubeURL, l.YoutubeVideoID, l.VideoURL,
		l.OrderIndex, l.IsFree,
	)
	out, err := scanLesson(row)
	if err != nil {
		return nil, fmt.Errorf("create lesson: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*Lesson, error) {
	setParts := make([]string, 0, 10)
	args := make([]interface{}, 0, 11)
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
		UPDATE public.lessons SET %s WHERE id = $%d::uuid
		RETURNING %s
	`, strings.Join(setParts, ", "), pos, lessonSelectCols)

	out, err := scanLesson(r.pool.QueryRow(ctx, q, args...))
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
