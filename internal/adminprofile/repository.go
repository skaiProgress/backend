package adminprofile

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository loads and updates admin profile rows.
type Repository interface {
	GetByID(ctx context.Context, userID string) (*Profile, error)
	Update(ctx context.Context, userID string, req UpdateRequest) (*Profile, error)
	SetAvatarURL(ctx context.Context, userID, avatarURL string) (*Profile, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates an admin profile repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) GetByID(ctx context.Context, userID string) (*Profile, error) {
	const q = `
		SELECT p.id::text, COALESCE(p.email, ''), p.full_name, p.phone, p.position,
		       p.department, p.bio, p.avatar_url, p.updated_at
		FROM public.profiles p
		WHERE p.id = $1::uuid
		LIMIT 1
	`
	var p Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&p.ID, &p.Email, &p.FullName, &p.Phone, &p.Position,
		&p.Department, &p.Bio, &p.AvatarURL, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get admin profile: %w", err)
	}
	return &p, nil
}

func (r *postgresRepository) Update(ctx context.Context, userID string, req UpdateRequest) (*Profile, error) {
	const q = `
		UPDATE public.profiles
		SET full_name = COALESCE($2, full_name),
		    phone = $3,
		    position = $4,
		    department = $5,
		    bio = $6,
		    updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING id::text, COALESCE(email, ''), full_name, phone, position,
		          department, bio, avatar_url, updated_at
	`
	var p Profile
	err := r.pool.QueryRow(ctx, q,
		userID, req.FullName, req.Phone, req.Position, req.Department, req.Bio,
	).Scan(
		&p.ID, &p.Email, &p.FullName, &p.Phone, &p.Position,
		&p.Department, &p.Bio, &p.AvatarURL, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update admin profile: %w", err)
	}
	return &p, nil
}

func (r *postgresRepository) SetAvatarURL(ctx context.Context, userID, avatarURL string) (*Profile, error) {
	const q = `
		UPDATE public.profiles
		SET avatar_url = $2, updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING id::text, COALESCE(email, ''), full_name, phone, position,
		          department, bio, avatar_url, updated_at
	`
	var p Profile
	err := r.pool.QueryRow(ctx, q, userID, avatarURL).Scan(
		&p.ID, &p.Email, &p.FullName, &p.Phone, &p.Position,
		&p.Department, &p.Bio, &p.AvatarURL, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("set avatar url: %w", err)
	}
	return &p, nil
}
