package briefings

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/storage"
)

// Service handles briefing automation logic.
type Service struct {
	repo    Repository
	storage *storage.Local
}

// NewService creates a briefings service.
func NewService(repo Repository, fileStorage *storage.Local) *Service {
	return &Service{repo: repo, storage: fileStorage}
}

// isValidKind reports whether s is one of the five briefing kinds.
func isValidKind(k BriefingKind) bool {
	switch k {
	case KindIntroductory, KindPrimary, KindRepeat, KindUnscheduled, KindTargeted:
		return true
	default:
		return false
	}
}

// ── Briefing videos (super admin) ─────────────────────────────────────────────

func (s *Service) requireAdmin(ctx context.Context) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}
	if claims.Role != "admin" && claims.Role != "super_admin" {
		return errors.New("forbidden")
	}
	return nil
}

// ListBriefingVideos returns the videos uploaded for a course (admin).
func (s *Service) ListBriefingVideos(ctx context.Context, courseID string) ([]BriefingVideo, error) {
	if err := s.requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.repo.ListBriefingVideosByCourse(ctx, courseID)
}

// UploadBriefingVideo stores a video for a course/kind (admin) and marks the course as a briefing course.
func (s *Service) UploadBriefingVideo(ctx context.Context, courseID, kind, filename string, r io.Reader) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	if !isValidKind(BriefingKind(kind)) {
		return errors.New("invalid briefing_kind")
	}
	url, path, err := s.storage.SaveBriefingVideo(courseID, kind, filename, r)
	if err != nil {
		return fmt.Errorf("save briefing video: %w", err)
	}
	if err := s.repo.UpsertBriefingVideo(ctx, courseID, kind, url, path); err != nil {
		_ = s.storage.Delete(path)
		return err
	}
	// A course that has briefing videos is implicitly a briefing course.
	_ = s.repo.SetCourseBriefingFlag(ctx, courseID, true)
	return nil
}

// DeleteBriefingVideo removes a video for a course/kind (admin).
func (s *Service) DeleteBriefingVideo(ctx context.Context, courseID, kind string) error {
	if err := s.requireAdmin(ctx); err != nil {
		return err
	}
	path, err := s.repo.DeleteBriefingVideo(ctx, courseID, kind)
	if err != nil {
		return err
	}
	if path != "" {
		_ = s.storage.Delete(path)
	}
	return nil
}

// ListBriefingCourses returns courses the org-admin can use to create briefing links.
func (s *Service) ListBriefingCourses(ctx context.Context) ([]BriefingCourse, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	if _, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject); err != nil {
		return nil, err
	}
	return s.repo.ListBriefingCoursesForOrgAdmin(ctx, claims.Subject)
}

// ListOrgAdminEvents returns all calendar events for the org-admin's organization.
func (s *Service) ListOrgAdminEvents(ctx context.Context) ([]OrgEvent, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	return s.repo.ListOrgEvents(ctx, orgID)
}

// CreateManualBriefingEvent lets org-admin create a briefing link tied to a ПБ course video.
func (s *Service) CreateManualBriefingEvent(ctx context.Context, req CreateBriefingEventRequest) (*OrgEvent, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	employeeID := strings.TrimSpace(req.EmployeeID)
	if employeeID == "" {
		return nil, errors.New("employee_id is required")
	}

	courseID := strings.TrimSpace(req.CourseID)
	if courseID == "" {
		return nil, errors.New("course_id is required")
	}

	kind := BriefingKind(strings.TrimSpace(req.BriefingKind))
	if !isValidKind(kind) {
		return nil, errors.New("invalid briefing_kind")
	}

	inOrg, err := s.repo.EmployeeInOrg(ctx, orgID, employeeID)
	if err != nil {
		return nil, err
	}
	if !inOrg {
		return nil, errors.New("employee not in organization")
	}

	// Org-admin must have access to the course, and it must be a briefing course
	// with a video for the chosen kind.
	hasCourse, err := s.repo.OrgAdminHasCourse(ctx, claims.Subject, courseID)
	if err != nil {
		return nil, err
	}
	if !hasCourse {
		return nil, errors.New("course not available")
	}
	isBriefing, err := s.repo.CourseIsBriefing(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !isBriefing {
		return nil, errors.New("course is not a briefing course")
	}
	videoURL, err := s.repo.GetBriefingVideo(ctx, courseID, string(kind))
	if err != nil {
		return nil, err
	}
	if videoURL == "" {
		return nil, errors.New("no video for this briefing kind")
	}

	startsAt, err := time.Parse(time.RFC3339, strings.TrimSpace(req.StartsAt))
	if err != nil {
		return nil, errors.New("invalid starts_at, use RFC3339 format")
	}
	endsAt, err := time.Parse(time.RFC3339, strings.TrimSpace(req.EndsAt))
	if err != nil {
		return nil, errors.New("invalid ends_at, use RFC3339 format")
	}
	if !endsAt.After(startsAt) {
		return nil, errors.New("ends_at must be after starts_at")
	}

	_, empFullName, err := s.repo.GetEmployeeProfile(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	empName := employeeID
	if empFullName != nil && *empFullName != "" {
		empName = *empFullName
	}

	loc := strings.TrimSpace(req.Location)
	if loc == "" {
		loc = "Онлайн (видео-инструктаж)"
	}

	title := fmt.Sprintf("%s инструктаж — %s", kind.Label(), empName)
	adminID := claims.Subject

	return s.repo.CreateEvent(ctx, CreateEventInput{
		OrganizationID: orgID,
		EmployeeID:     &employeeID,
		CourseID:       &courseID,
		Title:          title,
		EventType:      "training",
		BriefingKind:   &kind,
		StartsAt:       startsAt,
		EndsAt:         &endsAt,
		Location:       loc,
		CreatedBy:      &adminID,
	})
}

// UpdateEventTime lets org-admin reschedule an event.
func (s *Service) UpdateEventTime(ctx context.Context, eventID, rawStartsAt string) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, rawStartsAt)
	if err != nil {
		return errors.New("invalid starts_at, use RFC3339 format")
	}
	return s.repo.UpdateEventTime(ctx, eventID, orgID, t)
}

// ListOrgAdminRecords returns all briefing records for the org-admin journal.
func (s *Service) ListOrgAdminRecords(ctx context.Context) ([]BriefingRecord, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	return s.repo.ListOrgRecords(ctx, orgID)
}

// InstructorSign lets org-admin sign a briefing record.
func (s *Service) InstructorSign(ctx context.Context, recordID string) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return err
	}
	return s.repo.InstructorSignRecord(ctx, recordID, orgID)
}

// DeleteOrgRecord removes a briefing journal entry and its linked calendar event.
func (s *Service) DeleteOrgRecord(ctx context.Context, recordID string) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByAdmin(ctx, claims.Subject)
	if err != nil {
		return err
	}
	return s.repo.DeleteOrgRecord(ctx, recordID, orgID)
}

// ListEmployeeEvents returns calendar events visible to the employee.
func (s *Service) ListEmployeeEvents(ctx context.Context) ([]OrgEvent, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByEmployee(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	return s.repo.ListEmployeeEvents(ctx, claims.Subject, orgID)
}

// ListEmployeeBriefings returns pending and past briefings for the employee.
func (s *Service) ListEmployeeBriefings(ctx context.Context) ([]EmployeeBriefing, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	orgID, err := s.repo.GetOrgIDByEmployee(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	events, err := s.repo.ListEmployeeEvents(ctx, claims.Subject, orgID)
	if err != nil {
		return nil, err
	}

	out := make([]EmployeeBriefing, 0, len(events))
	for _, ev := range events {
		if ev.BriefingKind == nil || ev.EmployeeID == nil || *ev.EmployeeID != claims.Subject {
			continue
		}
		rec, _ := s.repo.GetRecordByEvent(ctx, ev.ID, claims.Subject)
		eb := EmployeeBriefing{
			EventID:      ev.ID,
			Title:        ev.Title,
			BriefingKind: *ev.BriefingKind,
			StartsAt:     ev.StartsAt,
			EndsAt:       ev.EndsAt,
			Location:     ev.Location,
		}
		if rec != nil {
			eb.RecordID = &rec.ID
			eb.Confirmed = rec.EmployeeSigned
		}
		out = append(out, eb)
	}
	return out, nil
}

// getOwnedBriefingEvent loads an event and verifies it is a briefing assigned to the employee.
func (s *Service) getOwnedBriefingEvent(ctx context.Context, eventID, employeeID string) (*OrgEvent, error) {
	ev, err := s.repo.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if ev.BriefingKind == nil {
		return nil, errors.New("event has no briefing kind")
	}
	if ev.EmployeeID == nil || *ev.EmployeeID != employeeID {
		return nil, errors.New("event not found")
	}
	return ev, nil
}

// GetEmployeeBriefingDetail returns the briefing page payload for one event.
func (s *Service) GetEmployeeBriefingDetail(ctx context.Context, eventID string) (*EmployeeBriefingDetail, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	if _, err := s.repo.GetOrgIDByEmployee(ctx, claims.Subject); err != nil {
		return nil, err
	}

	ev, err := s.getOwnedBriefingEvent(ctx, eventID, claims.Subject)
	if err != nil {
		return nil, err
	}

	videoURL := ""
	if ev.CourseID != nil && *ev.CourseID != "" {
		videoURL, err = s.repo.GetBriefingVideo(ctx, *ev.CourseID, *ev.BriefingKind)
		if err != nil {
			return nil, err
		}
	}

	rec, _ := s.repo.GetRecordByEvent(ctx, ev.ID, claims.Subject)
	confirmed := rec != nil && rec.EmployeeSigned

	now := time.Now()
	notStarted := now.Before(ev.StartsAt)
	expired := ev.EndsAt != nil && now.After(*ev.EndsAt)

	return &EmployeeBriefingDetail{
		EventID:      ev.ID,
		Title:        ev.Title,
		BriefingKind: *ev.BriefingKind,
		Location:     ev.Location,
		StartsAt:     ev.StartsAt,
		EndsAt:       ev.EndsAt,
		VideoURL:     videoURL,
		Confirmed:    confirmed,
		WindowActive: !notStarted && !expired,
		NotStarted:   notStarted,
		Expired:      expired,
	}, nil
}

// CompleteBriefing records that the employee watched the briefing video and signs the journal.
func (s *Service) CompleteBriefing(ctx context.Context, eventID string) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}

	orgID, err := s.repo.GetOrgIDByEmployee(ctx, claims.Subject)
	if err != nil {
		return err
	}

	ev, err := s.getOwnedBriefingEvent(ctx, eventID, claims.Subject)
	if err != nil {
		return err
	}

	// Must be inside the validity window.
	now := time.Now()
	if now.Before(ev.StartsAt) {
		return errors.New("briefing window not started")
	}
	if ev.EndsAt != nil && now.After(*ev.EndsAt) {
		return errors.New("briefing window expired")
	}

	// A video must exist for this briefing.
	if ev.CourseID == nil || *ev.CourseID == "" {
		return errors.New("no video for this briefing kind")
	}
	videoURL, err := s.repo.GetBriefingVideo(ctx, *ev.CourseID, *ev.BriefingKind)
	if err != nil {
		return err
	}
	if videoURL == "" {
		return errors.New("no video for this briefing kind")
	}

	// Prevent double-completion.
	existing, err := s.repo.GetRecordByEvent(ctx, eventID, claims.Subject)
	if err != nil {
		return err
	}
	if existing != nil && existing.EmployeeSigned {
		return errors.New("briefing already confirmed")
	}

	_, empFullName, err := s.repo.GetEmployeeProfile(ctx, claims.Subject)
	if err != nil {
		return err
	}
	name := claims.Subject
	if empFullName != nil && *empFullName != "" {
		name = *empFullName
	}

	position, _ := s.repo.GetEmployeePosition(ctx, claims.Subject)
	if strings.TrimSpace(position) == "" {
		position = "Не указана"
	}

	today := time.Now().UTC().Format("2006-01-02")

	rec := BriefingRecord{
		OrganizationID: orgID,
		EventID:        &eventID,
		EmployeeID:     claims.Subject,
		EmployeeName:   name,
		Position:       position,
		BriefingKind:   *ev.BriefingKind,
		InstructorName: "Инструктор по ПБ",
		DateConducted:  today,
	}

	saved, err := s.repo.CreateRecord(ctx, rec)
	if err != nil {
		return fmt.Errorf("create record: %w", err)
	}

	// Watching the video counts as the employee's signature.
	if err := s.repo.EmployeeSignRecord(ctx, saved.ID, claims.Subject); err != nil {
		return fmt.Errorf("sign record: %w", err)
	}

	return nil
}

// ListEmployeeJournalRecords returns the employee's own briefing records.
func (s *Service) ListEmployeeJournalRecords(ctx context.Context) ([]BriefingRecord, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.repo.ListEmployeeRecords(ctx, claims.Subject)
}
