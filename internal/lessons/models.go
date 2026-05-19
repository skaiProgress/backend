package lessons

import "time"

// Lesson is a course lesson row.
type Lesson struct {
	ID             string    `json:"id"`
	CourseID       string    `json:"course_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	YoutubeURL     string    `json:"youtube_url"`
	YoutubeVideoID string    `json:"youtube_video_id"`
	OrderIndex     int       `json:"order_index"`
	IsFree         bool      `json:"is_free"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateRequest is POST /functions/v1/lessons body.
type CreateRequest struct {
	CourseID       string  `json:"course_id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	YoutubeURL     string  `json:"youtube_url"`
	YoutubeVideoID string  `json:"youtube_video_id"`
	OrderIndex     int     `json:"order_index"`
	IsFree         bool    `json:"is_free"`
}

// ReorderRequest is PATCH /functions/v1/lessons/reorder body.
type ReorderRequest struct {
	CourseID   string   `json:"course_id"`
	OrderedIDs []string `json:"ordered_ids"`
}

// UpdateRequest is PATCH /functions/v1/lessons/:id body.
type UpdateRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	YoutubeURL     *string `json:"youtube_url"`
	YoutubeVideoID *string `json:"youtube_video_id"`
	OrderIndex     *int    `json:"order_index"`
	IsFree         *bool   `json:"is_free"`
}
