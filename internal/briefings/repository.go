package briefings

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines persistence for events and briefing records.
type Repository interface {
	// Events
	CreateEvent(ctx context.Context, in CreateEventInput) (*OrgEvent, error)
	ListOrgEvents(ctx context.Context, orgID string) ([]OrgEvent, error)
	ListEmployeeEvents(ctx context.Context, employeeID, orgID string) ([]OrgEvent, error)
	GetEvent(ctx context.Context, eventID string) (*OrgEvent, error)
	UpdateEventTime(ctx context.Context, eventID, orgID string, startsAt time.Time) error

	// Briefing records
	CreateRecord(ctx context.Context, rec BriefingRecord) (*BriefingRecord, error)
	ListOrgRecords(ctx context.Context, orgID string) ([]BriefingRecord, error)
	ListEmployeeRecords(ctx context.Context, employeeID string) ([]BriefingRecord, error)
	GetRecordByEvent(ctx context.Context, eventID, employeeID string) (*BriefingRecord, error)
	EmployeeSignRecord(ctx context.Context, recordID, employeeID string) error
	InstructorSignRecord(ctx context.Context, recordID, orgID string) error
	DeleteOrgRecord(ctx context.Context, recordID, orgID string) error

	// Briefing videos (super admin)
	UpsertBriefingVideo(ctx context.Context, courseID, kind, videoURL, videoPath string) error
	ListBriefingVideosByCourse(ctx context.Context, courseID string) ([]BriefingVideo, error)
	GetBriefingVideo(ctx context.Context, courseID, kind string) (string, error)
	DeleteBriefingVideo(ctx context.Context, courseID, kind string) (videoPath string, err error)
	SetCourseBriefingFlag(ctx context.Context, courseID string, value bool) error

	// Helpers
	GetOrgIDByEmployee(ctx context.Context, employeeID string) (string, error)
	GetEmployeeProfile(ctx context.Context, employeeID string) (email string, fullName *string, err error)
	GetEmployeePosition(ctx context.Context, employeeID string) (string, error)
	GetOrgAdminProfile(ctx context.Context, orgAdminID string) (email string, fullName *string, err error)
	HasRecordForKind(ctx context.Context, employeeID string, kind BriefingKind) (bool, error)
	GetOrgIDByAdmin(ctx context.Context, orgAdminID string) (string, error)
	EmployeeInOrg(ctx context.Context, orgID, employeeID string) (bool, error)
	OrgAdminHasCourse(ctx context.Context, orgAdminID, courseID string) (bool, error)
	CourseIsBriefing(ctx context.Context, courseID string) (bool, error)
	ListBriefingCoursesForOrgAdmin(ctx context.Context, orgAdminID string) ([]BriefingCourse, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a briefings repository backed by PostgreSQL.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CreateEvent(ctx context.Context, in CreateEventInput) (*OrgEvent, error) {
	var bk *string
	if in.BriefingKind != nil {
		s := string(*in.BriefingKind)
		bk = &s
	}
	const q = `
		INSERT INTO public.org_events
			(organization_id, employee_id, course_id, title, event_type, briefing_kind, starts_at, ends_at, location, participants, created_by)
		VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11::uuid)
		RETURNING id::text, organization_id::text, employee_id::text, course_id::text, title, event_type, briefing_kind,
		          starts_at, ends_at, location, participants, created_at
	`
	var ev OrgEvent
	var empID *string
	var courseID *string
	err := r.pool.QueryRow(ctx, q,
		in.OrganizationID, in.EmployeeID, in.CourseID, in.Title, in.EventType, bk,
		in.StartsAt, in.EndsAt, in.Location, in.Participants, in.CreatedBy,
	).Scan(
		&ev.ID, &ev.OrganizationID, &empID, &courseID, &ev.Title, &ev.EventType, &ev.BriefingKind,
		&ev.StartsAt, &ev.EndsAt, &ev.Location, &ev.Participants, &ev.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create event: %w", err)
	}
	ev.EmployeeID = empID
	ev.CourseID = courseID
	ev.Time = ev.StartsAt.Format("15:04")
	return &ev, nil
}

func (r *postgresRepository) ListOrgEvents(ctx context.Context, orgID string) ([]OrgEvent, error) {
	const q = `
		SELECT id::text, organization_id::text, employee_id::text, course_id::text, title, event_type,
		       briefing_kind, starts_at, ends_at, location, participants, created_at
		FROM public.org_events
		WHERE organization_id = $1::uuid
		ORDER BY starts_at ASC
	`
	return r.scanEvents(ctx, q, orgID)
}

func (r *postgresRepository) ListEmployeeEvents(ctx context.Context, employeeID, orgID string) ([]OrgEvent, error) {
	const q = `
		SELECT id::text, organization_id::text, employee_id::text, course_id::text, title, event_type,
		       briefing_kind, starts_at, ends_at, location, participants, created_at
		FROM public.org_events
		WHERE organization_id = $1::uuid
		  AND (employee_id = $2::uuid OR employee_id IS NULL)
		ORDER BY starts_at ASC
	`
	rows, err := r.pool.Query(ctx, q, orgID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list employee events: %w", err)
	}
	defer rows.Close()
	return scanEventRows(rows)
}

func (r *postgresRepository) GetEvent(ctx context.Context, eventID string) (*OrgEvent, error) {
	const q = `
		SELECT id::text, organization_id::text, employee_id::text, course_id::text, title, event_type,
		       briefing_kind, starts_at, ends_at, location, participants, created_at
		FROM public.org_events
		WHERE id = $1::uuid
	`
	rows, err := r.pool.Query(ctx, q, eventID)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	defer rows.Close()
	evs, err := scanEventRows(rows)
	if err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return nil, pgx.ErrNoRows
	}
	return &evs[0], nil
}

func (r *postgresRepository) UpdateEventTime(ctx context.Context, eventID, orgID string, startsAt time.Time) error {
	const q = `
		UPDATE public.org_events
		SET starts_at = $1, updated_at = NOW()
		WHERE id = $2::uuid AND organization_id = $3::uuid
	`
	tag, err := r.pool.Exec(ctx, q, startsAt, eventID, orgID)
	if err != nil {
		return fmt.Errorf("update event time: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *postgresRepository) CreateRecord(ctx context.Context, rec BriefingRecord) (*BriefingRecord, error) {
	const q = `
		INSERT INTO public.briefing_records
			(organization_id, event_id, employee_id, employee_name, position,
			 briefing_kind, instructor_name, instructor_id, date_conducted)
		VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, $6, $7, $8::uuid, $9::date)
		RETURNING id::text, organization_id::text, event_id::text, employee_id::text,
		          employee_name, position, briefing_kind, instructor_name, instructor_id::text,
		          date_conducted::text, employee_signed, employee_signed_at,
		          instructor_signed, instructor_signed_at, created_at
	`
	var out BriefingRecord
	err := r.pool.QueryRow(ctx, q,
		rec.OrganizationID, rec.EventID, rec.EmployeeID, rec.EmployeeName, rec.Position,
		rec.BriefingKind, rec.InstructorName, rec.InstructorID, rec.DateConducted,
	).Scan(
		&out.ID, &out.OrganizationID, &out.EventID, &out.EmployeeID,
		&out.EmployeeName, &out.Position, &out.BriefingKind, &out.InstructorName, &out.InstructorID,
		&out.DateConducted, &out.EmployeeSigned, &out.EmployeeSignedAt,
		&out.InstructorSigned, &out.InstructorSignedAt, &out.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create briefing record: %w", err)
	}
	return &out, nil
}

func (r *postgresRepository) ListOrgRecords(ctx context.Context, orgID string) ([]BriefingRecord, error) {
	const q = `
		SELECT id::text, organization_id::text, event_id::text, employee_id::text,
		       employee_name, position, briefing_kind, instructor_name, instructor_id::text,
		       date_conducted::text, employee_signed, employee_signed_at,
		       instructor_signed, instructor_signed_at, created_at,
		       ROW_NUMBER() OVER (ORDER BY date_conducted, created_at) AS row_number
		FROM public.briefing_records
		WHERE organization_id = $1::uuid
		ORDER BY date_conducted DESC, created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, orgID)
	if err != nil {
		return nil, fmt.Errorf("list org records: %w", err)
	}
	defer rows.Close()
	return scanRecordRows(rows, true)
}

func (r *postgresRepository) ListEmployeeRecords(ctx context.Context, employeeID string) ([]BriefingRecord, error) {
	const q = `
		SELECT id::text, organization_id::text, event_id::text, employee_id::text,
		       employee_name, position, briefing_kind, instructor_name, instructor_id::text,
		       date_conducted::text, employee_signed, employee_signed_at,
		       instructor_signed, instructor_signed_at, created_at,
		       ROW_NUMBER() OVER (ORDER BY date_conducted, created_at) AS row_number
		FROM public.briefing_records
		WHERE employee_id = $1::uuid
		ORDER BY date_conducted DESC, created_at DESC
	`
	rows, err := r.pool.Query(ctx, q, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list employee records: %w", err)
	}
	defer rows.Close()
	return scanRecordRows(rows, true)
}

func (r *postgresRepository) GetRecordByEvent(ctx context.Context, eventID, employeeID string) (*BriefingRecord, error) {
	const q = `
		SELECT id::text, organization_id::text, event_id::text, employee_id::text,
		       employee_name, position, briefing_kind, instructor_name, instructor_id::text,
		       date_conducted::text, employee_signed, employee_signed_at,
		       instructor_signed, instructor_signed_at, created_at, 1
		FROM public.briefing_records
		WHERE event_id = $1::uuid AND employee_id = $2::uuid
		LIMIT 1
	`
	rows, err := r.pool.Query(ctx, q, eventID, employeeID)
	if err != nil {
		return nil, fmt.Errorf("get record by event: %w", err)
	}
	defer rows.Close()
	recs, err := scanRecordRows(rows, false)
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, nil
	}
	return &recs[0], nil
}

func (r *postgresRepository) EmployeeSignRecord(ctx context.Context, recordID, employeeID string) error {
	const q = `
		UPDATE public.briefing_records
		SET employee_signed = TRUE, employee_signed_at = NOW(), updated_at = NOW()
		WHERE id = $1::uuid AND employee_id = $2::uuid AND employee_signed = FALSE
	`
	tag, err := r.pool.Exec(ctx, q, recordID, employeeID)
	if err != nil {
		return fmt.Errorf("employee sign record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("record not found or already signed")
	}
	return nil
}

func (r *postgresRepository) InstructorSignRecord(ctx context.Context, recordID, orgID string) error {
	const q = `
		UPDATE public.briefing_records
		SET instructor_signed = TRUE, instructor_signed_at = NOW(), updated_at = NOW()
		WHERE id = $1::uuid AND organization_id = $2::uuid AND instructor_signed = FALSE
	`
	tag, err := r.pool.Exec(ctx, q, recordID, orgID)
	if err != nil {
		return fmt.Errorf("instructor sign record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("record not found or already signed")
	}
	return nil
}

func (r *postgresRepository) DeleteOrgRecord(ctx context.Context, recordID, orgID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var eventID *string
	const selQ = `
		SELECT event_id::text FROM public.briefing_records
		WHERE id = $1::uuid AND organization_id = $2::uuid
	`
	err = tx.QueryRow(ctx, selQ, recordID, orgID).Scan(&eventID)
	if errors.Is(err, pgx.ErrNoRows) {
		return pgx.ErrNoRows
	}
	if err != nil {
		return fmt.Errorf("find briefing record: %w", err)
	}

	const delRecQ = `
		DELETE FROM public.briefing_records
		WHERE id = $1::uuid AND organization_id = $2::uuid
	`
	if _, err = tx.Exec(ctx, delRecQ, recordID, orgID); err != nil {
		return fmt.Errorf("delete briefing record: %w", err)
	}

	if eventID != nil && *eventID != "" {
		const delEvQ = `
			DELETE FROM public.org_events
			WHERE id = $1::uuid AND organization_id = $2::uuid
		`
		if _, err = tx.Exec(ctx, delEvQ, *eventID, orgID); err != nil {
			return fmt.Errorf("delete linked event: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) GetOrgIDByEmployee(ctx context.Context, employeeID string) (string, error) {
	const q = `SELECT organization_id::text FROM public.profiles WHERE id = $1::uuid AND role = 'user'`
	var orgID string
	err := r.pool.QueryRow(ctx, q, employeeID).Scan(&orgID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("employee profile not found")
	}
	return orgID, err
}

func (r *postgresRepository) GetOrgIDByAdmin(ctx context.Context, orgAdminID string) (string, error) {
	const q = `SELECT organization_id::text FROM public.profiles WHERE id = $1::uuid AND role = 'org_admin'`
	var orgID string
	err := r.pool.QueryRow(ctx, q, orgAdminID).Scan(&orgID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("org-admin profile not found")
	}
	return orgID, err
}

func (r *postgresRepository) GetEmployeeProfile(ctx context.Context, employeeID string) (string, *string, error) {
	const q = `SELECT COALESCE(email,''), full_name FROM public.profiles WHERE id = $1::uuid`
	var email string
	var fullName *string
	err := r.pool.QueryRow(ctx, q, employeeID).Scan(&email, &fullName)
	return email, fullName, err
}

func (r *postgresRepository) GetEmployeePosition(ctx context.Context, employeeID string) (string, error) {
	const q = `SELECT COALESCE(position,'') FROM public.profiles WHERE id = $1::uuid`
	var position string
	err := r.pool.QueryRow(ctx, q, employeeID).Scan(&position)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return position, err
}

func (r *postgresRepository) GetOrgAdminProfile(ctx context.Context, orgAdminID string) (string, *string, error) {
	const q = `SELECT COALESCE(email,''), full_name FROM public.profiles WHERE id = $1::uuid AND role = 'org_admin'`
	var email string
	var fullName *string
	err := r.pool.QueryRow(ctx, q, orgAdminID).Scan(&email, &fullName)
	return email, fullName, err
}

func (r *postgresRepository) EmployeeInOrg(ctx context.Context, orgID, employeeID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.profiles
			WHERE id = $1::uuid AND organization_id = $2::uuid
			  AND role = 'user' AND is_active = TRUE
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, q, employeeID, orgID).Scan(&exists)
	return exists, err
}

func (r *postgresRepository) HasRecordForKind(ctx context.Context, employeeID string, kind BriefingKind) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.briefing_records
			WHERE employee_id = $1::uuid AND briefing_kind = $2
		)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, q, employeeID, string(kind)).Scan(&exists)
	return exists, err
}

// ── scanners ────────────────────────────────────────────────────────────────

func (r *postgresRepository) scanEvents(ctx context.Context, q string, args ...any) ([]OrgEvent, error) {
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()
	return scanEventRows(rows)
}

func scanEventRows(rows pgx.Rows) ([]OrgEvent, error) {
	out := make([]OrgEvent, 0)
	for rows.Next() {
		var ev OrgEvent
		var empID *string
		var courseID *string
		if err := rows.Scan(
			&ev.ID, &ev.OrganizationID, &empID, &courseID, &ev.Title, &ev.EventType,
			&ev.BriefingKind, &ev.StartsAt, &ev.EndsAt, &ev.Location, &ev.Participants, &ev.CreatedAt,
		); err != nil {
			return nil, err
		}
		ev.EmployeeID = empID
		ev.CourseID = courseID
		ev.Time = ev.StartsAt.Format("15:04")
		out = append(out, ev)
	}
	return out, rows.Err()
}

// ── briefing videos & course helpers ──────────────────────────────────────────

func (r *postgresRepository) UpsertBriefingVideo(ctx context.Context, courseID, kind, videoURL, videoPath string) error {
	const q = `
		INSERT INTO public.briefing_videos (course_id, briefing_kind, video_url, video_path)
		VALUES ($1::uuid, $2, $3, $4)
		ON CONFLICT (course_id, briefing_kind)
		DO UPDATE SET video_url = EXCLUDED.video_url, video_path = EXCLUDED.video_path, updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, q, courseID, kind, videoURL, videoPath)
	if err != nil {
		return fmt.Errorf("upsert briefing video: %w", err)
	}
	return nil
}

func (r *postgresRepository) ListBriefingVideosByCourse(ctx context.Context, courseID string) ([]BriefingVideo, error) {
	const q = `
		SELECT briefing_kind, video_url FROM public.briefing_videos
		WHERE course_id = $1::uuid
		ORDER BY briefing_kind
	`
	rows, err := r.pool.Query(ctx, q, courseID)
	if err != nil {
		return nil, fmt.Errorf("list briefing videos: %w", err)
	}
	defer rows.Close()
	out := make([]BriefingVideo, 0)
	for rows.Next() {
		var v BriefingVideo
		if err := rows.Scan(&v.BriefingKind, &v.VideoURL); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *postgresRepository) GetBriefingVideo(ctx context.Context, courseID, kind string) (string, error) {
	const q = `SELECT video_url FROM public.briefing_videos WHERE course_id = $1::uuid AND briefing_kind = $2`
	var url string
	err := r.pool.QueryRow(ctx, q, courseID, kind).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return url, err
}

func (r *postgresRepository) DeleteBriefingVideo(ctx context.Context, courseID, kind string) (string, error) {
	const q = `
		DELETE FROM public.briefing_videos
		WHERE course_id = $1::uuid AND briefing_kind = $2
		RETURNING video_path
	`
	var path string
	err := r.pool.QueryRow(ctx, q, courseID, kind).Scan(&path)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", pgx.ErrNoRows
	}
	return path, err
}

func (r *postgresRepository) SetCourseBriefingFlag(ctx context.Context, courseID string, value bool) error {
	const q = `UPDATE public.courses SET is_briefing_course = $2, updated_at = NOW() WHERE id = $1::uuid`
	_, err := r.pool.Exec(ctx, q, courseID, value)
	if err != nil {
		return fmt.Errorf("set course briefing flag: %w", err)
	}
	return nil
}

func (r *postgresRepository) CourseIsBriefing(ctx context.Context, courseID string) (bool, error) {
	const q = `SELECT is_briefing_course FROM public.courses WHERE id = $1::uuid`
	var ok bool
	err := r.pool.QueryRow(ctx, q, courseID).Scan(&ok)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return ok, err
}

func (r *postgresRepository) OrgAdminHasCourse(ctx context.Context, orgAdminID, courseID string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM public.course_assignments
			WHERE user_id = $1::uuid AND course_id = $2::uuid
		)
	`
	var ok bool
	err := r.pool.QueryRow(ctx, q, orgAdminID, courseID).Scan(&ok)
	return ok, err
}

func (r *postgresRepository) ListBriefingCoursesForOrgAdmin(ctx context.Context, orgAdminID string) ([]BriefingCourse, error) {
	const q = `
		SELECT c.id::text, c.title,
		       COALESCE(ARRAY_AGG(bv.briefing_kind ORDER BY bv.briefing_kind)
		                FILTER (WHERE bv.briefing_kind IS NOT NULL), '{}') AS kinds
		FROM public.course_assignments ca
		JOIN public.courses c ON c.id = ca.course_id
		LEFT JOIN public.briefing_videos bv ON bv.course_id = c.id
		WHERE ca.user_id = $1::uuid AND c.is_briefing_course = TRUE
		GROUP BY c.id, c.title
		ORDER BY c.title
	`
	rows, err := r.pool.Query(ctx, q, orgAdminID)
	if err != nil {
		return nil, fmt.Errorf("list briefing courses: %w", err)
	}
	defer rows.Close()
	out := make([]BriefingCourse, 0)
	for rows.Next() {
		var bc BriefingCourse
		if err := rows.Scan(&bc.CourseID, &bc.Title, &bc.Kinds); err != nil {
			return nil, err
		}
		out = append(out, bc)
	}
	return out, rows.Err()
}

func scanRecordRows(rows pgx.Rows, withRowNum bool) ([]BriefingRecord, error) {
	out := make([]BriefingRecord, 0)
	for rows.Next() {
		var rec BriefingRecord
		if withRowNum {
			if err := rows.Scan(
				&rec.ID, &rec.OrganizationID, &rec.EventID, &rec.EmployeeID,
				&rec.EmployeeName, &rec.Position, &rec.BriefingKind, &rec.InstructorName, &rec.InstructorID,
				&rec.DateConducted, &rec.EmployeeSigned, &rec.EmployeeSignedAt,
				&rec.InstructorSigned, &rec.InstructorSignedAt, &rec.CreatedAt, &rec.RowNumber,
			); err != nil {
				return nil, err
			}
		} else {
			if err := rows.Scan(
				&rec.ID, &rec.OrganizationID, &rec.EventID, &rec.EmployeeID,
				&rec.EmployeeName, &rec.Position, &rec.BriefingKind, &rec.InstructorName, &rec.InstructorID,
				&rec.DateConducted, &rec.EmployeeSigned, &rec.EmployeeSignedAt,
				&rec.InstructorSigned, &rec.InstructorSignedAt, &rec.CreatedAt, &rec.RowNumber,
			); err != nil {
				return nil, err
			}
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}
