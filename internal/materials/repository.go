package materials

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists course materials.
type Repository interface {
	ListByCourse(ctx context.Context, courseID string) ([]Material, error)
	Insert(ctx context.Context, m Material) (*Material, error)
	GetByID(ctx context.Context, id string) (*Material, error)
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a materials repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) ListByCourse(ctx context.Context, courseID string) ([]Material, error) {
	const q = `
		SELECT id::text, course_id::text, name, file_url, file_type, file_size, created_at
		FROM public.course_materials
		WHERE course_id = $1::uuid
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, fmt.Errorf("list materials: %w", err)
	}
	defer rows.Close()

	out := make([]Material, 0)
	for rows.Next() {
		var m Material
		if err := rows.Scan(
			&m.ID, &m.CourseID, &m.Name, &m.FileURL, &m.FileType, &m.FileSize, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *postgresRepository) Insert(ctx context.Context, m Material) (*Material, error) {
	const q = `
		INSERT INTO public.course_materials (course_id, name, file_url, file_type, file_size)
		VALUES ($1::uuid, $2, $3, $4, $5)
		RETURNING id::text, course_id::text, name, file_url, file_type, file_size, created_at
	`
	var out Material
	err := r.pool.QueryRow(ctx, q, m.CourseID, m.Name, m.FileURL, m.FileType, m.FileSize).Scan(
		&out.ID, &out.CourseID, &out.Name, &out.FileURL, &out.FileType, &out.FileSize, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert material: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Material, error) {
	const q = `
		SELECT id::text, course_id::text, name, file_url, file_type, file_size, created_at
		FROM public.course_materials WHERE id = $1::uuid
	`
	var out Material
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&out.ID, &out.CourseID, &out.Name, &out.FileURL, &out.FileType, &out.FileSize, &out.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM public.course_materials WHERE id = $1::uuid`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
