package adminprofile

import "time"

// Profile is the admin settings profile row.
type Profile struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FullName   *string    `json:"full_name"`
	Phone      *string    `json:"phone"`
	Position   *string    `json:"position"`
	Department *string    `json:"department"`
	Bio        *string    `json:"bio"`
	AvatarURL  *string    `json:"avatar_url"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// UpdateRequest is PATCH /admin/profile body.
type UpdateRequest struct {
	FullName   *string `json:"full_name"`
	Phone      *string `json:"phone"`
	Position   *string `json:"position"`
	Department *string `json:"department"`
	Bio        *string `json:"bio"`
}
