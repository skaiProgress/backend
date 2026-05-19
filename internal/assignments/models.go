package assignments

import "time"

// Assignment row returned to clients.
type Assignment struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	CourseID   string     `json:"course_id"`
	AssignedBy *string    `json:"assigned_by"`
	AssignedAt time.Time  `json:"assigned_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	Status     string     `json:"status"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

// CreateRequest is POST /functions/v1/assignments.
type CreateRequest struct {
	UserID     string  `json:"user_id"`
	CourseID   string  `json:"course_id"`
	ExpiresAt  *string `json:"expires_at"`
	AssignedBy *string `json:"assigned_by"`
}

// BulkRequest is POST /functions/v1/assignments/bulk.
type BulkRequest struct {
	UserIDs    []string `json:"user_ids"`
	CourseIDs  []string `json:"course_ids"`
	ExpiresAt  *string  `json:"expires_at"`
	AssignedBy *string  `json:"assigned_by"`
}

// BulkResponse is bulk assign result.
type BulkResponse struct {
	Count int `json:"count"`
}

// UserProfileSnippet is nested user info on assignment list rows.
type UserProfileSnippet struct {
	Email    string  `json:"email"`
	FullName *string `json:"full_name"`
}

// CourseSnippet is nested course info on assignment list rows.
type CourseSnippet struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// ListItem is a row for GET /functions/v1/assignments.
type ListItem struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	CourseID    string              `json:"course_id"`
	AssignedAt  time.Time           `json:"assigned_at"`
	ExpiresAt   *time.Time          `json:"expires_at"`
	Status      string              `json:"status"`
	RevokedAt   *time.Time          `json:"revoked_at"`
	UserProfile *UserProfileSnippet `json:"user_profile"`
	Courses     *CourseSnippet      `json:"courses"`
}
