package orgadmin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists org-admin scoped data.
type Repository interface {
	LoadOrgContext(ctx context.Context, userID string) (*OrgContext, error)
	GetStats(ctx context.Context, orgID, orgAdminUserID string) (*Stats, error)
	ListMembers(ctx context.Context, orgID string) ([]Member, error)
	CreateAuthUser(ctx context.Context, email, passwordHash string) (string, error)
	UpsertMemberProfile(ctx context.Context, userID, email, orgID string, fullName *string, isActive bool) error
	MemberInOrg(ctx context.Context, orgID, userID string) (bool, error)
	OrgAdminHasCourse(ctx context.Context, orgAdminUserID, courseID string) (bool, error)
	ListMyCourses(ctx context.Context, orgAdminUserID, orgID string) ([]Course, error)
	GetCourseDetail(ctx context.Context, orgAdminUserID, orgID, courseID string) (*CourseDetail, error)
	ListAssignments(ctx context.Context, orgID, orgAdminUserID, courseID string) ([]AssignmentRow, error)
	UpsertAssignment(ctx context.Context, userID, courseID, assignedBy string, expiresAt *string) (*AssignmentRow, error)
	RevokeAssignment(ctx context.Context, orgID, orgAdminUserID, assignmentID string) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates an org-admin repository.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) LoadOrgContext(ctx context.Context, userID string) (*OrgContext, error) {
	const q = `
		SELECT p.organization_id::text, o.name, COALESCE(p.email, ''), p.full_name
		FROM public.profiles p
		JOIN public.organizations o ON o.id = p.organization_id
		WHERE p.id = $1::uuid
		  AND p.role = 'org_admin'
		  AND p.is_active = TRUE
		  AND p.organization_id IS NOT NULL
	`
	var orgID, orgName, email string
	var fullName *string
	err := r.pool.QueryRow(ctx, q, userID).Scan(&orgID, &orgName, &email, &fullName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load org context: %w", err)
	}
	return &OrgContext{
		UserID:           userID,
		Email:            email,
		FullName:         fullName,
		OrganizationID:   orgID,
		OrganizationName: orgName,
	}, nil
}

func (r *postgresRepository) GetStats(ctx context.Context, orgID, orgAdminUserID string) (*Stats, error) {
	const q = `
		SELECT
			o.name,
			(SELECT COUNT(*)::int FROM public.profiles p
			 WHERE p.organization_id = $1::uuid AND p.role = 'user'),
			(SELECT COUNT(*)::int FROM public.profiles p
			 WHERE p.organization_id = $1::uuid AND p.role = 'user' AND p.is_active),
			(SELECT COUNT(*)::int FROM public.course_assignments ca
			 JOIN public.courses c ON c.id = ca.course_id
			 WHERE ca.user_id = $2::uuid AND ca.status = 'active' AND c.status = 'published'),
			(SELECT COUNT(*)::int FROM public.course_assignments ca
			 JOIN public.profiles p ON p.id = ca.user_id
			 WHERE p.organization_id = $1::uuid
			   AND p.role = 'user'
			   AND ca.status = 'active'
			   AND ca.course_id IN (
			     SELECT course_id FROM public.course_assignments
			     WHERE user_id = $2::uuid AND status = 'active'
			   ))
		FROM public.organizations o
		WHERE o.id = $1::uuid
	`
	var s Stats
	err := r.pool.QueryRow(ctx, q, orgID, orgAdminUserID).Scan(
		&s.OrganizationName,
		&s.EmployeesTotal,
		&s.EmployeesActive,
		&s.AssignedCourses,
		&s.EmployeeAssignments,
	)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	return &s, nil
}

func (r *postgresRepository) ListMembers(ctx context.Context, orgID string) ([]Member, error) {
	const q = `
		SELECT id::text, COALESCE(email, ''), full_name, role, is_active, created_at
		FROM public.profiles
		WHERE organization_id = $1::uuid AND role = 'user'
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	out := make([]Member, 0)
	for rows.Next() {
		var m Member
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

func (r *postgresRepository) UpsertMemberProfile(
	ctx context.Context,
	userID, email, orgID string,
	fullName *string,
	isActive bool,
) error {
	const q = `
		INSERT INTO public.profiles (id, email, full_name, role, is_active, organization_id)
		VALUES ($1::uuid, $2, $3, 'user', $4, $5::uuid)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			full_name = EXCLUDED.full_name,
			is_active = EXCLUDED.is_active,
			organization_id = EXCLUDED.organization_id,
			updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, q, userID, email, fullName, isActive, orgID)
	if err != nil {
		return fmt.Errorf("upsert member profile: %w", err)
	}
	return nil
}

func (r *postgresRepository) MemberInOrg(ctx context.Context, orgID, userID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.profiles
			WHERE id = $1::uuid AND organization_id = $2::uuid AND role = 'user'
		)
	`
	var ok bool
	err := r.pool.QueryRow(ctx, q, userID, orgID).Scan(&ok)
	return ok, err
}

func (r *postgresRepository) OrgAdminHasCourse(ctx context.Context, orgAdminUserID, courseID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.course_assignments ca
			JOIN public.courses c ON c.id = ca.course_id
			WHERE ca.user_id = $1::uuid
			  AND ca.course_id = $2::uuid
			  AND ca.status = 'active'
			  AND c.status = 'published'
		)
	`
	var ok bool
	err := r.pool.QueryRow(ctx, q, orgAdminUserID, courseID).Scan(&ok)
	return ok, err
}

func (r *postgresRepository) ListMyCourses(ctx context.Context, orgAdminUserID, orgID string) ([]Course, error) {
	const q = `
		SELECT ca.id::text, c.id::text, c.title, c.description, c.cover_url,
		       ca.assigned_at, ca.expires_at,
		       (ca.expires_at IS NOT NULL AND ca.expires_at < NOW()),
		       COALESCE((
		         SELECT COUNT(*)::int
		         FROM public.course_assignments eca
		         JOIN public.profiles p ON p.id = eca.user_id
		         WHERE p.organization_id = $2::uuid
		           AND p.role = 'user'
		           AND eca.course_id = c.id
		           AND eca.status = 'active'
		       ), 0)
		FROM public.course_assignments ca
		JOIN public.courses c ON c.id = ca.course_id
		WHERE ca.user_id = $1::uuid
		  AND ca.status = 'active'
		  AND c.status = 'published'
		ORDER BY ca.assigned_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgAdminUserID, orgID)
	if err != nil {
		return nil, fmt.Errorf("list my courses: %w", err)
	}
	defer rows.Close()

	out := make([]Course, 0)
	for rows.Next() {
		var c Course
		if err := rows.Scan(
			&c.AssignmentID, &c.CourseID, &c.Title, &c.Description, &c.CoverURL,
			&c.AssignedAt, &c.ExpiresAt, &c.IsExpired, &c.AssignedCount,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetCourseDetail(ctx context.Context, orgAdminUserID, orgID, courseID string) (*CourseDetail, error) {
	const q = `
		SELECT ca.id::text, c.id::text, c.title, c.description, c.cover_url,
		       ca.assigned_at, ca.expires_at,
		       (ca.expires_at IS NOT NULL AND ca.expires_at < NOW()),
		       COALESCE((
		         SELECT COUNT(*)::int
		         FROM public.course_assignments eca
		         JOIN public.profiles p ON p.id = eca.user_id
		         WHERE p.organization_id = $3::uuid
		           AND p.role = 'user'
		           AND eca.course_id = c.id
		           AND eca.status = 'active'
		       ), 0),
		       COALESCE((SELECT COUNT(*)::int FROM public.lessons l WHERE l.course_id = c.id), 0)
		FROM public.course_assignments ca
		JOIN public.courses c ON c.id = ca.course_id
		WHERE ca.user_id = $1::uuid
		  AND ca.course_id = $2::uuid
		  AND ca.status = 'active'
		  AND c.status = 'published'
		LIMIT 1
	`
	var d CourseDetail
	err := r.pool.QueryRow(ctx, q, orgAdminUserID, courseID, orgID).Scan(
		&d.AssignmentID, &d.CourseID, &d.Title, &d.Description, &d.CoverURL,
		&d.AssignedAt, &d.ExpiresAt, &d.IsExpired, &d.AssignedCount, &d.LessonCount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get course detail: %w", err)
	}
	return &d, nil
}

func (r *postgresRepository) ListAssignments(ctx context.Context, orgID, orgAdminUserID, courseID string) ([]AssignmentRow, error) {
	const q = `
		SELECT ca.id::text, ca.user_id::text, ca.course_id::text,
		       ca.assigned_at, ca.expires_at, ca.status,
		       COALESCE(p.email, ''), p.full_name
		FROM public.course_assignments ca
		JOIN public.profiles p ON p.id = ca.user_id
		WHERE p.organization_id = $1::uuid
		  AND p.role = 'user'
		  AND ca.status = 'active'
		  AND ca.course_id IN (
		    SELECT course_id FROM public.course_assignments
		    WHERE user_id = $2::uuid AND status = 'active'
		  )
		  AND ($3 = '' OR ca.course_id = $3::uuid)
		ORDER BY ca.assigned_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID, orgAdminUserID, courseID)
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}
	defer rows.Close()

	out := make([]AssignmentRow, 0)
	for rows.Next() {
		var a AssignmentRow
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.CourseID, &a.AssignedAt, &a.ExpiresAt, &a.Status,
			&a.Email, &a.FullName,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *postgresRepository) UpsertAssignment(
	ctx context.Context,
	userID, courseID, assignedBy string,
	expiresAt *string,
) (*AssignmentRow, error) {
	const q = `
		INSERT INTO public.course_assignments (
			user_id, course_id, assigned_by, assigned_at, expires_at, status, revoked_at
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, NOW(), $4::timestamptz, 'active', NULL
		)
		ON CONFLICT (user_id, course_id) DO UPDATE SET
			assigned_by = EXCLUDED.assigned_by,
			assigned_at = NOW(),
			expires_at = EXCLUDED.expires_at,
			status = 'active',
			revoked_at = NULL,
			updated_at = NOW()
		RETURNING id::text, user_id::text, course_id::text, assigned_at, expires_at, status
	`
	var a AssignmentRow
	err := r.pool.QueryRow(ctx, q, userID, courseID, assignedBy, expiresAt).Scan(
		&a.ID, &a.UserID, &a.CourseID, &a.AssignedAt, &a.ExpiresAt, &a.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert assignment: %w", err)
	}

	const profileQ = `SELECT COALESCE(email, ''), full_name FROM public.profiles WHERE id = $1::uuid`
	_ = r.pool.QueryRow(ctx, profileQ, userID).Scan(&a.Email, &a.FullName)
	return &a, nil
}

func (r *postgresRepository) RevokeAssignment(ctx context.Context, orgID, orgAdminUserID, assignmentID string) error {
	const q = `
		UPDATE public.course_assignments ca
		SET status = 'revoked', revoked_at = NOW(), updated_at = NOW()
		FROM public.profiles p
		WHERE ca.id = $1::uuid
		  AND ca.status = 'active'
		  AND p.id = ca.user_id
		  AND p.organization_id = $2::uuid
		  AND p.role = 'user'
		  AND ca.course_id IN (
		    SELECT course_id FROM public.course_assignments
		    WHERE user_id = $3::uuid AND status = 'active'
		  )
	`
	tag, err := r.pool.Exec(ctx, q, assignmentID, orgID, orgAdminUserID)
	if err != nil {
		return fmt.Errorf("revoke assignment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
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
