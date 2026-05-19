package auth

import "time"

// AuthUser is a row from auth.users used for credential verification.
type AuthUser struct {
	ID                string
	Email             string
	EncryptedPassword string
	BannedUntil       *time.Time
	DeletedAt         *time.Time
}

// Profile is a row from public.profiles.
type Profile struct {
	ID        string
	Email     string
	Role      string
	FullName  *string
	IsActive  bool
	AvatarURL *string
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned after successful login.
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
	User        UserPublic `json:"user"`
}

// UserPublic is safe user data exposed to clients.
type UserPublic struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	FullName *string `json:"full_name,omitempty"`
}

// ChangePasswordRequest is the JSON body for POST /functions/v1/auth/change-password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// MeResponse is returned by GET /functions/v1/auth/me.
type MeResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	FullName  *string `json:"full_name,omitempty"`
	IsActive  bool    `json:"is_active"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
