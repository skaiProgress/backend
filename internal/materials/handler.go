package materials

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes materials HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a materials handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /functions/v1/materials?course_id=
func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context(), c.QueryParam("course_id"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Upload handles POST /functions/v1/materials (multipart).
func (h *Handler) Upload(c echo.Context) error {
	courseID := c.FormValue("course_id")
	file, err := c.FormFile("file")
	if err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "file is required")
	}
	out, err := h.service.UploadFromMultipart(c.Request().Context(), courseID, file)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// Delete handles DELETE /functions/v1/materials/:id.
func (h *Handler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return httputil.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
	case errors.Is(err, auth.ErrForbidden):
		return httputil.ErrorJSON(c, http.StatusForbidden, "Forbidden: requires admin or super_admin role")
	case errors.Is(err, pgx.ErrNoRows):
		return httputil.ErrorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		if msg == "course_id and file are required" || msg == "invalid file size" || msg == "file is required" {
			return httputil.ErrorJSON(c, http.StatusBadRequest, msg)
		}
		return httputil.ErrorJSON(c, http.StatusInternalServerError, msg)
	}
}
