package organizations

import "time"

// Organization is a tenant company on the platform.
type Organization struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	BIN           *string   `json:"bin"`
	Phone         *string   `json:"phone"`
	Email         *string   `json:"email"`
	Address       *string   `json:"address"`
	ContactPerson *string   `json:"contact_person"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserCount     int       `json:"user_count,omitempty"`
}

// OrganizationWithUsers includes members for detail view.
type OrganizationWithUsers struct {
	Organization
	Users []OrgMember `json:"users"`
}

// OrgMember is a user belonging to an organization.
type OrgMember struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  *string   `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateRequest body for POST /admin/organizations.
type CreateRequest struct {
	Name          string  `json:"name"`
	BIN           *string `json:"bin"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
	ContactPerson *string `json:"contact_person"`
	IsActive      *bool   `json:"is_active"`
}

// UpdateRequest body for PATCH /admin/organizations/:id.
type UpdateRequest struct {
	Name          *string `json:"name"`
	BIN           *string `json:"bin"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	Address       *string `json:"address"`
	ContactPerson *string `json:"contact_person"`
	IsActive      *bool   `json:"is_active"`
}

// AddMemberRequest body for POST /admin/organizations/:id/users.
type AddMemberRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName *string `json:"full_name"`
	Role     string  `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// AddMemberResponse after creating org user.
type AddMemberResponse struct {
	UserID string `json:"user_id"`
}
