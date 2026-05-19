package materials

import "time"

// Material is a course material row.
type Material struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Name      string    `json:"name"`
	FileURL   string    `json:"file_url"`
	FileType  *string   `json:"file_type"`
	FileSize  *int64    `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}
