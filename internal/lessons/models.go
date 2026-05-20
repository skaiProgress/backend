package lessons

import "time"

const (
	VideoSourceYouTube = "youtube"
	VideoSourceUpload  = "upload"
)

// Lesson is a course lesson row.
type Lesson struct {
	ID             string    `json:"id"`
	CourseID       string    `json:"course_id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	VideoSource    string    `json:"video_source"`
	YoutubeURL     *string   `json:"youtube_url"`
	YoutubeVideoID *string   `json:"youtube_video_id"`
	VideoURL       *string   `json:"video_url"`
	OrderIndex     int       `json:"order_index"`
	IsFree         bool      `json:"is_free"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateRequest is POST /functions/v1/lessons body (YouTube).
type CreateRequest struct {
	CourseID       string  `json:"course_id"`
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	YoutubeURL     string  `json:"youtube_url"`
	YoutubeVideoID string  `json:"youtube_video_id"`
	OrderIndex     int     `json:"order_index"`
	IsFree         bool    `json:"is_free"`
}

// CreateUploadRequest is POST /functions/v1/lessons multipart body.
type CreateUploadRequest struct {
	CourseID    string
	Title       string
	Description *string
	OrderIndex  int
	IsFree      bool
}

// ReorderRequest is PATCH /functions/v1/lessons/reorder body.
type ReorderRequest struct {
	CourseID   string   `json:"course_id"`
	OrderedIDs []string `json:"ordered_ids"`
}

// UpdateRequest is PATCH /functions/v1/lessons/:id body (YouTube).
type UpdateRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	YoutubeURL     *string `json:"youtube_url"`
	YoutubeVideoID *string `json:"youtube_video_id"`
	OrderIndex     *int    `json:"order_index"`
	IsFree         *bool   `json:"is_free"`
}

// UpdateUploadRequest is PATCH /functions/v1/lessons/:id multipart metadata.
type UpdateUploadRequest struct {
	Title       *string
	Description *string
	OrderIndex  *int
	IsFree      *bool
}
