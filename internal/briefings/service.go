package briefings

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"aiqadam-backend/internal/auth"
)

// Service handles briefing automation logic.
type Service struct {
	repo Repository
}

// NewService creates a briefings service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// defaultTime is when a briefing event is scheduled by default (14:00 local).
const defaultHour = 14

// ScheduleIntroductoryBriefing creates a вводный-инструктаж calendar event
// immediately when org-admin registers a new employee.
func (s *Service) ScheduleIntroductoryBriefing(
	ctx context.Context,
	orgID, employeeID, employeeName, orgAdminID string,
) error {
	// Default time: today at 14:00 UTC
	now := time.Now().UTC()
	startsAt := time.Date(now.Year(), now.Month(), now.Day(), defaultHour, 0, 0, 0, time.UTC)
	kind := KindIntroductory
	loc := "Переговорная / Учебный класс"
	title := fmt.Sprintf("Вводный инструктаж — %s", employeeName)

	_, err := s.repo.CreateEvent(ctx, CreateEventInput{
		OrganizationID: orgID,
		EmployeeID:     &employeeID,
		Title:          title,
		EventType:      "training",
		BriefingKind:   &kind,
		StartsAt:       startsAt,
		Location:       loc,
		CreatedBy:      &orgAdminID,
	})
	return err
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

// CreateManualBriefingEvent lets org-admin schedule a targeted or unscheduled briefing.
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

	kind := BriefingKind(strings.TrimSpace(req.BriefingKind))
	if kind != KindTargeted && kind != KindUnscheduled {
		return nil, errors.New("briefing_kind must be targeted or unscheduled")
	}

	inOrg, err := s.repo.EmployeeInOrg(ctx, orgID, employeeID)
	if err != nil {
		return nil, err
	}
	if !inOrg {
		return nil, errors.New("employee not in organization")
	}

	t, err := time.Parse(time.RFC3339, strings.TrimSpace(req.StartsAt))
	if err != nil {
		return nil, errors.New("invalid starts_at, use RFC3339 format")
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
		loc = "Переговорная / Учебный класс"
	}

	title := fmt.Sprintf("%s инструктаж — %s", kind.Label(), empName)
	adminID := claims.Subject

	return s.repo.CreateEvent(ctx, CreateEventInput{
		OrganizationID: orgID,
		EmployeeID:     &employeeID,
		Title:          title,
		EventType:      "training",
		BriefingKind:   &kind,
		StartsAt:       t,
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

// ConfirmBriefing is called when an employee clicks "Подтвердить прохождение".
// It creates a briefing_record with employee_signed=true, then schedules the
// next briefing in the chain.
func (s *Service) ConfirmBriefing(ctx context.Context, eventID string, req ConfirmBriefingRequest) error {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return errors.New("unauthorized")
	}

	orgID, err := s.repo.GetOrgIDByEmployee(ctx, claims.Subject)
	if err != nil {
		return err
	}

	// Prevent double-confirm
	existing, err := s.repo.GetRecordByEvent(ctx, eventID, claims.Subject)
	if err != nil {
		return err
	}
	if existing != nil && existing.EmployeeSigned {
		return errors.New("briefing already confirmed")
	}

	// Find the event to get its kind
	orgEvents, err := s.repo.ListOrgEvents(ctx, orgID)
	if err != nil {
		return err
	}
	var ev *OrgEvent
	for i := range orgEvents {
		if orgEvents[i].ID == eventID {
			ev = &orgEvents[i]
			break
		}
	}
	if ev == nil {
		return errors.New("event not found")
	}
	if ev.BriefingKind == nil {
		return errors.New("event has no briefing kind")
	}
	if ev.EmployeeID == nil || *ev.EmployeeID != claims.Subject {
		return errors.New("event not found")
	}

	_, empFullName, err := s.repo.GetEmployeeProfile(ctx, claims.Subject)
	if err != nil {
		return err
	}
	name := claims.Subject
	if empFullName != nil && *empFullName != "" {
		name = *empFullName
	}

	position := req.Position
	if position == "" {
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

	// Employee immediately signs
	if err := s.repo.EmployeeSignRecord(ctx, saved.ID, claims.Subject); err != nil {
		return fmt.Errorf("sign record: %w", err)
	}

	// Schedule next briefing in the chain
	kind := BriefingKind(*ev.BriefingKind)
	switch kind {
	case KindIntroductory:
		// Next: первичный — завтра в 14:00
		s.scheduleNext(ctx, orgID, claims.Subject, name, KindPrimary, 1)
	case KindPrimary:
		// Next: повторный — через 6 месяцев
		s.scheduleNext(ctx, orgID, claims.Subject, name, KindRepeat, 183)
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

// scheduleNext creates the next briefing event after daysFromNow days.
func (s *Service) scheduleNext(
	ctx context.Context,
	orgID, employeeID, employeeName string,
	kind BriefingKind,
	daysFromNow int,
) {
	future := time.Now().UTC().AddDate(0, 0, daysFromNow)
	startsAt := time.Date(future.Year(), future.Month(), future.Day(), defaultHour, 0, 0, 0, time.UTC)
	title := fmt.Sprintf("%s инструктаж — %s", kind.Label(), employeeName)

	_, _ = s.repo.CreateEvent(ctx, CreateEventInput{
		OrganizationID: orgID,
		EmployeeID:     &employeeID,
		Title:          title,
		EventType:      "training",
		BriefingKind:   &kind,
		StartsAt:       startsAt,
		Location:       "Переговорная / Учебный класс",
	})
}
