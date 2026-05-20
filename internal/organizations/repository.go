package organizations

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

// Repository persists organizations and org members.
type Repository interface {
	List(ctx context.Context, search string) ([]Organization, error)
	GetByID(ctx context.Context, id string) (*Organization, error)
	GetWithUsers(ctx context.Context, id string) (*OrganizationWithUsers, error)
	Create(ctx context.Context, o Organization) (*Organization, error)
	Update(ctx context.Context, id string, fields map[string]interface{}) (*Organization, error)
	Delete(ctx context.Context, id string) error
	ListMembers(ctx context.Context, orgID string) ([]OrgMember, error)
	CreateAuthUser(ctx context.Context, email, passwordHash string) (string, error)
	UpsertOrgProfile(ctx context.Context, userID, email, orgID string, fullName *string, role string, isActive bool) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a PostgreSQL organizations repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) List(ctx context.Context, search string) ([]Organization, error) {
	const q = `
		SELECT o.id::text, o.name, o.bin, o.phone, o.email, o.address, o.contact_person,
		       o.is_active, o.created_at, o.updated_at,
		       COALESCE(COUNT(p.id), 0)::int AS user_count
		FROM public.organizations o
		LEFT JOIN public.profiles p ON p.organization_id = o.id
		WHERE (
			$1 = ''
			OR o.name ILIKE '%' || $1 || '%'
			OR o.bin ILIKE '%' || $1 || '%'
			OR o.phone ILIKE '%' || $1 || '%'
		)
		GROUP BY o.id
		ORDER BY o.created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, strings.TrimSpace(search))
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	defer rows.Close()

	out := make([]Organization, 0)
	for rows.Next() {
		var item Organization
		if err := rows.Scan(
			&item.ID, &item.Name, &item.BIN, &item.Phone, &item.Email, &item.Address,
			&item.ContactPerson, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &item.UserCount,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Organization, error) {
	const q = `
		SELECT o.id::text, o.name, o.bin, o.phone, o.email, o.address, o.contact_person,
		       o.is_active, o.created_at, o.updated_at,
		       COALESCE((
		         SELECT COUNT(*)::int FROM public.profiles p WHERE p.organization_id = o.id
		       ), 0)
		FROM public.organizations o
		WHERE o.id = $1::uuid
	`
	var item Organization
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&item.ID, &item.Name, &item.BIN, &item.Phone, &item.Email, &item.Address,
		&item.ContactPerson, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &item.UserCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return &item, nil
}

func (r *postgresRepository) GetWithUsers(ctx context.Context, id string) (*OrganizationWithUsers, error) {
	org, err := r.GetByID(ctx, id)
	if err != nil || org == nil {
		return nil, err
	}
	users, err := r.ListMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	return &OrganizationWithUsers{
		Organization: *org,
		Users:        users,
	}, nil
}

func (r *postgresRepository) Create(ctx context.Context, o Organization) (*Organization, error) {
	const q = `
		INSERT INTO public.organizations (name, bin, phone, email, address, contact_person, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text, name, bin, phone, email, address, contact_person, is_active, created_at, updated_at
	`
	var out Organization
	err := r.pool.QueryRow(ctx, q,
		o.Name, o.BIN, o.Phone, o.Email, o.Address, o.ContactPerson, o.IsActive,
	).Scan(
		&out.ID, &out.Name, &out.BIN, &out.Phone, &out.Email, &out.Address,
		&out.ContactPerson, &out.IsActive, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("organization with this BIN already exists")
		}
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Update(ctx context.Context, id string, fields map[string]interface{}) (*Organization, error) {
	setParts := make([]string, 0, 8)
	args := make([]interface{}, 0, 10)
	argPos := 1

	for _, key := range []string{"name", "bin", "phone", "email", "address", "contact_person", "is_active"} {
		if v, ok := fields[key]; ok {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", key, argPos))
			args = append(args, v)
			argPos++
		}
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	q := fmt.Sprintf(
		`UPDATE public.organizations SET %s WHERE id = $%d::uuid
		 RETURNING id::text, name, bin, phone, email, address, contact_person, is_active, created_at, updated_at`,
		strings.Join(setParts, ", "),
		argPos,
	)

	var out Organization
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&out.ID, &out.Name, &out.BIN, &out.Phone, &out.Email, &out.Address,
		&out.ContactPerson, &out.IsActive, &out.CreatedAt, &out.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("organization with this BIN already exists")
		}
		return nil, fmt.Errorf("update organization: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM public.organizations WHERE id = $1::uuid`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete organization: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) ListMembers(ctx context.Context, orgID string) ([]OrgMember, error) {
	const q = `
		SELECT id::text, COALESCE(email, ''), full_name, role, is_active, created_at
		FROM public.profiles
		WHERE organization_id = $1::uuid
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, fmt.Errorf("list org members: %w", err)
	}
	defer rows.Close()

	out := make([]OrgMember, 0)
	for rows.Next() {
		var m OrgMember
		if err := rows.Scan(&m.ID, &m.Email, &m.FullName, &m.Role, &m.IsActive, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
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

func (r *postgresRepository) UpsertOrgProfile(
	ctx context.Context,
	userID, email, orgID string,
	fullName *string,
	role string,
	isActive bool,
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
	_, err := r.pool.Exec(ctx, q, userID, email, fullName, role, isActive, orgID)
	if err != nil {
		return fmt.Errorf("upsert org profile: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
