package orgadmin

import (
	"context"
	"errors"
	"strings"
	"time"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// BriefingScheduler is a minimal interface that orgadmin uses to trigger
// briefing automation without importing the full briefings package (avoids cycles).
type BriefingScheduler interface {
	ScheduleIntroductoryBriefing(ctx context.Context, orgID, employeeID, employeeName, orgAdminID string) error
}

// Service handles org-admin cabinet logic.
type Service struct {
	repo      Repository
	briefings BriefingScheduler
}

// NewService creates an org-admin service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// SetBriefingScheduler injects the briefing scheduler after construction
// to avoid import cycles.
func (s *Service) SetBriefingScheduler(bs BriefingScheduler) {
	s.briefings = bs
}

func (s *Service) orgContext(ctx context.Context) (*OrgContext, error) {
	claims, err := auth.ClaimsFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}
	if claims.Role != "org_admin" {
		return nil, ErrForbidden
	}
	oc, err := s.repo.LoadOrgContext(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}
	if oc == nil {
		return nil, ErrForbidden
	}
	return oc, nil
}

// GetStats returns dashboard metrics.
func (s *Service) GetStats(ctx context.Context) (*Stats, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetStats(ctx, oc.OrganizationID, oc.UserID)
}

// GetProfile returns org-admin profile with organization.
func (s *Service) GetProfile(ctx context.Context) (*ProfileResponse, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	return &ProfileResponse{
		ID:               oc.UserID,
		Email:            oc.Email,
		FullName:         oc.FullName,
		OrganizationID:   oc.OrganizationID,
		OrganizationName: oc.OrganizationName,
		Role:             "org_admin",
	}, nil
}

// ListMembers returns employees in the organization.
func (s *Service) ListMembers(ctx context.Context) ([]Member, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListMembers(ctx, oc.OrganizationID)
}

// CreateMember adds an employee (role user) to the organization.
func (s *Service) CreateMember(ctx context.Context, req CreateMemberRequest) (*CreateMemberResponse, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(req.Email)
	password := req.Password
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := s.repo.CreateAuthUser(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpsertMemberProfile(ctx, userID, email, oc.OrganizationID, req.FullName, isActive); err != nil {
		return nil, err
	}

	// Auto-schedule introductory briefing for the new employee.
	if s.briefings != nil {
		empName := email
		if req.FullName != nil && *req.FullName != "" {
			empName = *req.FullName
		}
		// Run in background so member creation succeeds even if scheduling fails.
		go func() {
			_ = s.briefings.ScheduleIntroductoryBriefing(
				context.Background(), oc.OrganizationID, userID, empName, oc.UserID,
			)
		}()
	}

	return &CreateMemberResponse{UserID: userID}, nil
}

// ListCourses returns courses assigned to the org-admin.
func (s *Service) ListCourses(ctx context.Context) ([]Course, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListMyCourses(ctx, oc.UserID, oc.OrganizationID)
}

// GetCourse returns one assigned course.
func (s *Service) GetCourse(ctx context.Context, courseID string) (*CourseDetail, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	if courseID == "" {
		return nil, ErrInvalidInput
	}
	d, err := s.repo.GetCourseDetail(ctx, oc.UserID, oc.OrganizationID, courseID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, pgx.ErrNoRows
	}
	return d, nil
}

// ListAssignments lists employee assignments for org-admin's courses.
func (s *Service) ListAssignments(ctx context.Context, courseID string) ([]AssignmentRow, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListAssignments(ctx, oc.OrganizationID, oc.UserID, strings.TrimSpace(courseID))
}

// CreateAssignment assigns a course to an org employee.
func (s *Service) CreateAssignment(ctx context.Context, req CreateAssignmentRequest) (*AssignmentRow, error) {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserID == "" || req.CourseID == "" {
		return nil, errors.New("user_id and course_id are required")
	}

	ok, err := s.repo.MemberInOrg(ctx, oc.OrganizationID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("user is not an employee of your organization")
	}

	hasCourse, err := s.repo.OrgAdminHasCourse(ctx, oc.UserID, req.CourseID)
	if err != nil {
		return nil, err
	}
	if !hasCourse {
		return nil, errors.New("course is not assigned to you")
	}

	var expires *string
	if req.ExpiresAt != nil && strings.TrimSpace(*req.ExpiresAt) != "" {
		raw := strings.TrimSpace(*req.ExpiresAt)
		if _, err := time.Parse(time.RFC3339, raw); err != nil {
			return nil, errors.New("invalid expires_at")
		}
		expires = &raw
	}

	return s.repo.UpsertAssignment(ctx, req.UserID, req.CourseID, oc.UserID, expires)
}

// RevokeAssignment removes an employee's course assignment.
func (s *Service) RevokeAssignment(ctx context.Context, assignmentID string) error {
	oc, err := s.orgContext(ctx)
	if err != nil {
		return err
	}
	if assignmentID == "" {
		return ErrInvalidInput
	}
	return s.repo.RevokeAssignment(ctx, oc.OrganizationID, oc.UserID, assignmentID)
}
