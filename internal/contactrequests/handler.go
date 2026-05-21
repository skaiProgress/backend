package contactrequests

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes contact request HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a contact requests handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /contact-requests (public).
func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, CreateResponse{
		ID:      out.ID,
		Message: "Заявка успешно отправлена",
	})
}

// List handles GET /functions/v1/admin/contact-requests.
func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context(), c.QueryParam("status"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Get handles GET /functions/v1/admin/contact-requests/:id.
func (h *Handler) Get(c echo.Context) error {
	item, err := h.service.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// Update handles PATCH /functions/v1/admin/contact-requests/:id.
func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.UpdateStatus(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// CountNew handles GET /functions/v1/admin/contact-requests/count/new.
func (h *Handler) CountNew(c echo.Context) error {
	count, err := h.service.CountNew(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]int{"count": count})
}

func errorJSON(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return errorJSON(c, http.StatusUnauthorized, "Invalid or expired token")
	case errors.Is(err, ErrForbidden):
		return errorJSON(c, http.StatusForbidden, "Forbidden: requires super_admin role")
	case errors.Is(err, ErrInvalidInput):
		return errorJSON(c, http.StatusBadRequest, "invalid input")
	case errors.Is(err, pgx.ErrNoRows):
		return errorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		switch msg {
		case "name is required", "email is required", "phone is required", "message is required",
			"invalid email", "name is too long", "phone is too long", "message is too long",
			"company is too long", "invalid status":
			return errorJSON(c, http.StatusBadRequest, msg)
		default:
			return errorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
