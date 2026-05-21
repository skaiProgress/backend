package contactrequests

import "time"

// ContactRequest is a public contact form submission.
type ContactRequest struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Company   *string   `json:"company"`
	Message   *string   `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRequest body for POST /contact-requests.
type CreateRequest struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Phone   string  `json:"phone"`
	Company *string `json:"company"`
	Message *string `json:"message"`
}

// UpdateRequest body for PATCH /admin/contact-requests/:id.
type UpdateRequest struct {
	Status *string `json:"status"`
}

// CreateResponse after submitting a contact form.
type CreateResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
