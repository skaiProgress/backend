package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines data access for authentication.
type Repository interface {
	FindAuthUserByEmail(ctx context.Context, email string) (*AuthUser, error)
	FindAuthUserByID(ctx context.Context, userID string) (*AuthUser, error)
	FindProfileByUserID(ctx context.Context, userID string) (*Profile, error)
	UpdatePassword(ctx context.Context, userID, encryptedPassword string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a PostgreSQL auth repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) FindAuthUserByEmail(ctx context.Context, email string) (*AuthUser, error) {
	const q = `
		SELECT id::text, email, encrypted_password, banned_until, deleted_at
		FROM auth.users
		WHERE LOWER(email) = LOWER($1)
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var u AuthUser
	err := r.pool.QueryRow(ctx, q, email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.BannedUntil,
		&u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find auth user: %w", err)
	}
	return &u, nil
}

func (r *postgresRepository) FindAuthUserByID(ctx context.Context, userID string) (*AuthUser, error) {
	const q = `
		SELECT id::text, email, encrypted_password, banned_until, deleted_at
		FROM auth.users
		WHERE id = $1::uuid
		  AND deleted_at IS NULL
		LIMIT 1
	`

	var u AuthUser
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.BannedUntil,
		&u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find auth user by id: %w", err)
	}
	return &u, nil
}

func (r *postgresRepository) UpdatePassword(ctx context.Context, userID, encryptedPassword string) error {
	const q = `
		UPDATE auth.users
		SET encrypted_password = $2, updated_at = NOW()
		WHERE id = $1::uuid
		  AND deleted_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, q, userID, encryptedPassword)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) FindProfileByUserID(ctx context.Context, userID string) (*Profile, error) {
	const q = `
		SELECT id::text, COALESCE(email, ''), role, full_name, is_active, avatar_url
		FROM public.profiles
		WHERE id = $1::uuid
		LIMIT 1
	`

	var p Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&p.ID,
		&p.Email,
		&p.Role,
		&p.FullName,
		&p.IsActive,
		&p.AvatarURL,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find profile: %w", err)
	}
	return &p, nil
}
