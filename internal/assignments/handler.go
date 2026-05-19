package assignments

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes assignment HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an assignments handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /functions/v1/assignments.
func (h *Handler) List(c echo.Context) error {
	activeOnly := c.QueryParam("status") != "all"
	items, err := h.service.List(
		c.Request().Context(),
		c.QueryParam("user_id"),
		c.QueryParam("course_id"),
		activeOnly,
	)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Create handles POST /functions/v1/assignments.
func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// Bulk handles POST /functions/v1/assignments/bulk.
func (h *Handler) Bulk(c echo.Context) error {
	var req BulkRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Bulk(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// Revoke handles DELETE /functions/v1/assignments/:id.
func (h *Handler) Revoke(c echo.Context) error {
	if err := h.service.Revoke(c.Request().Context(), c.Param("id")); err != nil {
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
		if msg == "user_id and course_id are required" ||
			msg == "Выберите хотя бы одного пользователя и один курс" ||
			msg == "invalid expires_at" ||
			msg == "Этот курс уже назначен пользователю" {
			status := http.StatusBadRequest
			return httputil.ErrorJSON(c, status, msg)
		}
		return httputil.ErrorJSON(c, http.StatusInternalServerError, msg)
	}
}
