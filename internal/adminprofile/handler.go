package adminprofile

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes admin profile HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an admin profile handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get handles GET /functions/v1/admin/profile.
func (h *Handler) Get(c echo.Context) error {
	p, err := h.service.Get(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, p)
}

// Update handles PATCH /functions/v1/admin/profile.
func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	p, err := h.service.Update(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, p)
}

// UploadAvatar handles POST /functions/v1/admin/profile/avatar.
func (h *Handler) UploadAvatar(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "file is required")
	}
	p, err := h.service.UploadAvatar(c.Request().Context(), file)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, p)
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return httputil.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
	case errors.Is(err, auth.ErrForbidden):
		return httputil.ErrorJSON(c, http.StatusForbidden, "Forbidden: requires admin or super_admin role")
	case errors.Is(err, pgx.ErrNoRows):
		return httputil.ErrorJSON(c, http.StatusNotFound, "profile not found")
	default:
		msg := err.Error()
		switch msg {
		case "full_name must be at least 2 characters", "full_name is too long",
			"file is required", "invalid file size", "file must be an image":
			return httputil.ErrorJSON(c, http.StatusBadRequest, msg)
		default:
			return httputil.ErrorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
