package assignments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListFilter filters assignment list queries.
type ListFilter struct {
	UserID   string
	CourseID string
	ActiveOnly bool
}

// Repository persists course assignments.
type Repository interface {
	List(ctx context.Context, f ListFilter) ([]ListItem, error)
	UpsertOne(ctx context.Context, a Assignment) (*Assignment, error)
	UpsertBulk(ctx context.Context, rows []Assignment) (int, error)
	Revoke(ctx context.Context, id string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates an assignments repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) List(ctx context.Context, f ListFilter) ([]ListItem, error) {
	const q = `
		SELECT ca.id::text, ca.user_id::text, ca.course_id::text,
		       ca.assigned_at, ca.expires_at, ca.status, ca.revoked_at,
		       COALESCE(p.email, ''), p.full_name,
		       c.id::text, c.title, c.status
		FROM public.course_assignments ca
		JOIN public.profiles p ON p.id = ca.user_id
		JOIN public.courses c ON c.id = ca.course_id
		WHERE ($1 = '' OR ca.user_id = $1::uuid)
		  AND ($2 = '' OR ca.course_id = $2::uuid)
		  AND ($3 = false OR ca.status = 'active')
		ORDER BY ca.assigned_at DESC
	`
	rows, err := r.pool.Query(ctx, q, f.UserID, f.CourseID, f.ActiveOnly)
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}
	defer rows.Close()

	out := make([]ListItem, 0)
	for rows.Next() {
		var item ListItem
		var email string
		var fullName *string
		var courseID, courseTitle, courseStatus string
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.CourseID,
			&item.AssignedAt, &item.ExpiresAt, &item.Status, &item.RevokedAt,
			&email, &fullName,
			&courseID, &courseTitle, &courseStatus,
		); err != nil {
			return nil, err
		}
		item.UserProfile = &UserProfileSnippet{Email: email, FullName: fullName}
		item.Courses = &CourseSnippet{ID: courseID, Title: courseTitle, Status: courseStatus}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *postgresRepository) UpsertOne(ctx context.Context, a Assignment) (*Assignment, error) {
	const q = `
		INSERT INTO public.course_assignments (
			user_id, course_id, assigned_by, assigned_at, expires_at, status, revoked_at
		) VALUES ($1::uuid, $2::uuid, NULLIF($3, '')::uuid, $4, $5, 'active', NULL)
		ON CONFLICT (user_id, course_id) DO UPDATE SET
			assigned_by = EXCLUDED.assigned_by,
			assigned_at = EXCLUDED.assigned_at,
			expires_at = EXCLUDED.expires_at,
			status = 'active',
			revoked_at = NULL,
			updated_at = NOW()
		RETURNING id::text, user_id::text, course_id::text, assigned_by::text,
		          assigned_at, expires_at, status, revoked_at
	`
	assignedBy := ""
	if a.AssignedBy != nil {
		assignedBy = *a.AssignedBy
	}
	var out Assignment
	err := r.pool.QueryRow(ctx, q,
		a.UserID, a.CourseID, assignedBy, a.AssignedAt, a.ExpiresAt,
	).Scan(
		&out.ID, &out.UserID, &out.CourseID, &out.AssignedBy,
		&out.AssignedAt, &out.ExpiresAt, &out.Status, &out.RevokedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert assignment: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) UpsertBulk(ctx context.Context, rows []Assignment) (int, error) {
	count := 0
	for _, row := range rows {
		if _, err := r.UpsertOne(ctx, row); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (r *postgresRepository) Revoke(ctx context.Context, id string) error {
	const q = `
		UPDATE public.course_assignments
		SET status = 'revoked', revoked_at = NOW(), updated_at = NOW()
		WHERE id = $1::uuid AND status = 'active'
	`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func parseExpiresAt(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		return nil, errors.New("invalid expires_at")
	}
	return &t, nil
}
