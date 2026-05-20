package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines admin user persistence.
type Repository interface {
	FindProfileRole(ctx context.Context, userID string) (string, error)
	CreateAuthUser(ctx context.Context, email, passwordHash string) (string, error)
	UpsertProfile(ctx context.Context, userID, email string, fullName *string, role string, isActive bool, organizationID *string) error
	OrganizationExists(ctx context.Context, orgID string) (bool, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateProfile(ctx context.Context, userID string, fullName *string, role *string, isActive *bool) error
	DeleteAuthUsers(ctx context.Context, userIDs []string) error
	ListProfiles(ctx context.Context, search string) ([]Profile, error)
	GetProfileByID(ctx context.Context, userID string) (*Profile, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a PostgreSQL users repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) FindProfileRole(ctx context.Context, userID string) (string, error) {
	const q = `SELECT role FROM public.profiles WHERE id = $1::uuid`
	var role string
	err := r.pool.QueryRow(ctx, q, userID).Scan(&role)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("find profile role: %w", err)
	}
	return role, nil
}

func (r *postgresRepository) CreateAuthUser(ctx context.Context, email, passwordHash string) (string, error) {
	const q = `
		INSERT INTO auth.users (
			instance_id, id, aud, role, email, encrypted_password,
			email_confirmed_at, raw_app_meta_data, raw_user_meta_data,
			created_at, updated_at
		) VALUES (
			'00000000-0000-0000-0000-000000000000',
			$1, 'authenticated', 'authenticated', $2, $3,
			NOW(),
			'{"provider":"email","providers":["email"]}'::jsonb,
			'{"email_verified":true}'::jsonb,
			NOW(), NOW()
		)
		RETURNING id::text
	`

	id := uuid.New().String()
	err := r.pool.QueryRow(ctx, q, id, email, passwordHash).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return "", ErrEmailExists
		}
		return "", fmt.Errorf("create auth user: %w", err)
	}
	return id, nil
}

func (r *postgresRepository) OrganizationExists(ctx context.Context, orgID string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM public.organizations WHERE id = $1::uuid)`
	var exists bool
	if err := r.pool.QueryRow(ctx, q, orgID).Scan(&exists); err != nil {
		return false, fmt.Errorf("organization exists: %w", err)
	}
	return exists, nil
}

func (r *postgresRepository) UpsertProfile(
	ctx context.Context,
	userID, email string,
	fullName *string,
	role string,
	isActive bool,
	organizationID *string,
) error {
	const q = `
		INSERT INTO public.profiles (id, email, full_name, role, is_active, organization_id)
		VALUES ($1::uuid, $2, $3, $4, $5, $6::uuid)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			full_name = EXCLUDED.full_name,
			role = EXCLUDED.role,
			is_active = EXCLUDED.is_active,
			organization_id = EXCLUDED.organization_id,
			updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, q, userID, email, fullName, role, isActive, organizationID)
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	return nil
}

func (r *postgresRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	const q = `
		UPDATE auth.users
		SET encrypted_password = $2, updated_at = NOW()
		WHERE id = $1::uuid AND deleted_at IS NULL
	`
	tag, err := r.pool.Exec(ctx, q, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) UpdateProfile(
	ctx context.Context,
	userID string,
	fullName *string,
	role *string,
	isActive *bool,
) error {
	setParts := make([]string, 0, 4)
	args := make([]interface{}, 0, 5)
	argPos := 1

	if fullName != nil {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", argPos))
		args = append(args, *fullName)
		argPos++
	}
	if role != nil {
		setParts = append(setParts, fmt.Sprintf("role = $%d", argPos))
		args = append(args, *role)
		argPos++
	}
	if isActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *isActive)
		argPos++
	}

	if len(setParts) == 0 {
		return nil
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, userID)

	q := fmt.Sprintf(
		`UPDATE public.profiles SET %s WHERE id = $%d::uuid`,
		strings.Join(setParts, ", "),
		argPos,
	)

	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) DeleteAuthUsers(ctx context.Context, userIDs []string) error {
	const q = `
		DELETE FROM auth.users
		WHERE id = ANY($1::uuid[])
	`
	_, err := r.pool.Exec(ctx, q, userIDs)
	if err != nil {
		return fmt.Errorf("delete auth users: %w", err)
	}
	return nil
}

func (r *postgresRepository) ListProfiles(ctx context.Context, search string) ([]Profile, error) {
	const q = `
		SELECT p.id::text, COALESCE(p.email, ''), p.full_name, p.role, p.is_active, p.created_at,
		       COALESCE(COUNT(ca.id) FILTER (WHERE ca.status = 'active'), 0)::int,
		       p.organization_id::text, o.name
		FROM public.profiles p
		LEFT JOIN public.course_assignments ca ON ca.user_id = p.id
		LEFT JOIN public.organizations o ON o.id = p.organization_id
		WHERE (
			$1 = ''
			OR p.email ILIKE '%' || $1 || '%'
			OR p.full_name ILIKE '%' || $1 || '%'
		)
		GROUP BY p.id, p.organization_id, o.name
		ORDER BY p.created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, strings.TrimSpace(search))
	if err != nil {
		return nil, fmt.Errorf("list profiles: %w", err)
	}
	defer rows.Close()

	out := make([]Profile, 0)
	for rows.Next() {
		var p Profile
		if err := rows.Scan(
			&p.ID, &p.Email, &p.FullName, &p.Role, &p.IsActive, &p.CreatedAt, &p.AssignmentCount,
			&p.OrganizationID, &p.OrganizationName,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetProfileByID(ctx context.Context, userID string) (*Profile, error) {
	const q = `
		SELECT p.id::text, COALESCE(p.email, ''), p.full_name, p.role, p.is_active, p.created_at,
		       COALESCE((
		         SELECT COUNT(*)::int FROM public.course_assignments ca
		         WHERE ca.user_id = p.id AND ca.status = 'active'
		       ), 0),
		       p.organization_id::text, o.name
		FROM public.profiles p
		LEFT JOIN public.organizations o ON o.id = p.organization_id
		WHERE p.id = $1::uuid
	`
	var p Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&p.ID, &p.Email, &p.FullName, &p.Role, &p.IsActive, &p.CreatedAt, &p.AssignmentCount,
		&p.OrganizationID, &p.OrganizationName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
