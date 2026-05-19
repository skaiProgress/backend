package employee

import "time"

// MyCourse is an assigned published course row.
type MyCourse struct {
	AssignmentID string     `json:"assignment_id"`
	CourseID     string     `json:"course_id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	CoverURL     *string    `json:"cover_url"`
	AssignedAt   time.Time  `json:"assigned_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	IsExpired    bool       `json:"is_expired"`
}

// MyCourseDetail is course detail for an assigned course.
type MyCourseDetail struct {
	CourseID    string     `json:"course_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	CoverURL    *string    `json:"cover_url"`
	AssignedAt  time.Time  `json:"assigned_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IsExpired   bool       `json:"is_expired"`
}

// MyLesson is a lesson visible to an assigned employee.
type MyLesson struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	YoutubeURL     string  `json:"youtube_url"`
	YoutubeVideoID string  `json:"youtube_video_id"`
	OrderIndex     int     `json:"order_index"`
	IsFree         bool    `json:"is_free"`
}

// Profile is the employee's own profile.
type Profile struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	FullName *string `json:"full_name"`
	Role     string  `json:"role"`
	IsActive bool    `json:"is_active"`
}

// UpdateProfileRequest is PATCH /employee/profile body.
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
}
