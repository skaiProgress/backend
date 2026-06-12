package orgadmin

import "time"

// OrgContext is the authenticated org-admin scope.
type OrgContext struct {
	UserID           string
	Email            string
	FullName         *string
	OrganizationID   string
	OrganizationName string
}

// Stats is dashboard metrics for the organization.
type Stats struct {
	OrganizationName   string `json:"organization_name"`
	EmployeesTotal     int    `json:"employees_total"`
	EmployeesActive    int    `json:"employees_active"`
	AssignedCourses    int    `json:"assigned_courses"`
	EmployeeAssignments int   `json:"employee_assignments"`
}

// Member is an employee in the organization.
type Member struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  *string   `json:"full_name"`
	Position  *string   `json:"position"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMemberRequest body for POST /org-admin/members.
type CreateMemberRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName *string `json:"full_name"`
	Position *string `json:"position"`
	IsActive *bool   `json:"is_active"`
}

// CreateMemberResponse after creating employee.
type CreateMemberResponse struct {
	UserID string `json:"user_id"`
}

// Course is a course assigned to the org-admin by super admin.
type Course struct {
	AssignmentID string     `json:"assignment_id"`
	CourseID     string     `json:"course_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	CoverURL     *string    `json:"cover_url"`
	AssignedAt   time.Time  `json:"assigned_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	IsExpired    bool       `json:"is_expired"`
	AssignedCount int       `json:"assigned_count"`
}

// CourseDetail extends Course with org assignment info.
type CourseDetail struct {
	Course
	LessonCount int `json:"lesson_count"`
}

// AssignmentRow is an assignment to an org employee.
type AssignmentRow struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	CourseID   string     `json:"course_id"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	Status     string     `json:"status"`
	Email      string     `json:"email"`
	FullName   *string    `json:"full_name"`
}

// CreateAssignmentRequest assigns a course to an org employee.
type CreateAssignmentRequest struct {
	UserID    string  `json:"user_id"`
	CourseID  string  `json:"course_id"`
	ExpiresAt *string `json:"expires_at"`
}

// ProfileResponse for org-admin profile.
type ProfileResponse struct {
	ID               string  `json:"id"`
	Email            string  `json:"email"`
	FullName         *string `json:"full_name"`
	Role             string  `json:"role"`
	OrganizationID   string  `json:"organization_id"`
	OrganizationName string  `json:"organization_name"`
}
