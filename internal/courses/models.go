package courses

import "time"

// Course is a training course row.
type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CoverURL    *string   `json:"cover_url"`
	CreatedBy   *string   `json:"created_by"`
	IsBriefingCourse bool  `json:"is_briefing_course"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CourseWithCount includes lesson count for list views.
type CourseWithCount struct {
	Course
	LessonCount int `json:"lesson_count"`
}

// CreateRequest is POST /functions/v1/courses body.
type CreateRequest struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Status           string  `json:"status"`
	CoverURL         *string `json:"cover_url"`
	IsBriefingCourse *bool   `json:"is_briefing_course"`
}

// UpdateRequest is PATCH /functions/v1/courses/:id body.
type UpdateRequest struct {
	Title            *string `json:"title"`
	Description      *string `json:"description"`
	Status           *string `json:"status"`
	CoverURL         *string `json:"cover_url"`
	IsBriefingCourse *bool   `json:"is_briefing_course"`
}
