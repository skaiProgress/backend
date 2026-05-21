package contactrequests

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists contact form submissions.
type Repository interface {
	Create(ctx context.Context, item ContactRequest) (*ContactRequest, error)
	List(ctx context.Context, status string) ([]ContactRequest, error)
	GetByID(ctx context.Context, id string) (*ContactRequest, error)
	UpdateStatus(ctx context.Context, id, status string) (*ContactRequest, error)
	CountByStatus(ctx context.Context, status string) (int, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a PostgreSQL contact requests repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func scanContactRequest(row pgx.Row) (*ContactRequest, error) {
	var out ContactRequest
	err := row.Scan(
		&out.ID, &out.Name, &out.Email, &out.Phone, &out.Company, &out.Message,
		&out.Status, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *postgresRepository) Create(ctx context.Context, item ContactRequest) (*ContactRequest, error) {
	const q = `
		INSERT INTO public.contact_requests (name, email, phone, company, message, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, name, email, phone, company, message, status, created_at, updated_at
	`
	return scanContactRequest(r.pool.QueryRow(
		ctx, q, item.Name, item.Email, item.Phone, item.Company, item.Message, item.Status,
	))
}

func (r *postgresRepository) List(ctx context.Context, status string) ([]ContactRequest, error) {
	status = strings.TrimSpace(status)
	q := `
		SELECT id::text, name, email, phone, company, message, status, created_at, updated_at
		FROM public.contact_requests
	`
	args := []any{}
	if status != "" {
		q += ` WHERE status = $1`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list contact requests: %w", err)
	}
	defer rows.Close()

	out := make([]ContactRequest, 0)
	for rows.Next() {
		var item ContactRequest
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Email, &item.Phone, &item.Company, &item.Message,
			&item.Status, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*ContactRequest, error) {
	const q = `
		SELECT id::text, name, email, phone, company, message, status, created_at, updated_at
		FROM public.contact_requests
		WHERE id = $1::uuid
	`
	out, err := scanContactRequest(r.pool.QueryRow(ctx, q, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("get contact request: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id, status string) (*ContactRequest, error) {
	const q = `
		UPDATE public.contact_requests
		SET status = $2
		WHERE id = $1::uuid
		RETURNING id::text, name, email, phone, company, message, status, created_at, updated_at
	`
	out, err := scanContactRequest(r.pool.QueryRow(ctx, q, id, status))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("update contact request: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	const q = `SELECT COUNT(*)::int FROM public.contact_requests WHERE status = $1`
	var count int
	if err := r.pool.QueryRow(ctx, q, status).Scan(&count); err != nil {
		return 0, fmt.Errorf("count contact requests: %w", err)
	}
	return count, nil
}
