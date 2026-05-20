package users

import "time"

// AddUserRequest is the body for POST /functions/v1/admin-add-user.
type AddUserRequest struct {
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	FullName       *string `json:"full_name"`
	Role           string  `json:"role"`
	IsActive       *bool   `json:"is_active"`
	OrganizationID *string `json:"organization_id"`
}

// AddUserResponse is returned after creating a user.
type AddUserResponse struct {
	UserID string `json:"user_id"`
}

// UpdateUserRequest is the body for POST /functions/v1/admin-update-user.
type UpdateUserRequest struct {
	UserID      string  `json:"user_id"`
	FullName    *string `json:"full_name"`
	Role        *string `json:"role"`
	IsActive    *bool   `json:"is_active"`
	NewPassword *string `json:"new_password"`
}

// DeleteUserRequest is the body for POST /functions/v1/admin-delete-user.
type DeleteUserRequest struct {
	UserIDs []string `json:"user_ids"`
}

// DeleteUserResponse is returned after bulk delete.
type DeleteUserResponse struct {
	OK      bool `json:"ok"`
	Deleted int  `json:"deleted"`
}

// OkResponse is a generic success payload.
type OkResponse struct {
	OK bool `json:"ok"`
}

// Caller holds authenticated admin context.
type Caller struct {
	ID   string
	Role string
}

// Profile is an admin-visible user row from public.profiles.
type Profile struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	FullName         *string   `json:"full_name"`
	Role             string    `json:"role"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	AssignmentCount  int       `json:"assignment_count"`
	OrganizationID   *string   `json:"organization_id,omitempty"`
	OrganizationName *string   `json:"organization_name,omitempty"`
}
