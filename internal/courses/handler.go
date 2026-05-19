package courses

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes course HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a courses handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /functions/v1/courses.
func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context(), c.QueryParam("search"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Create handles POST /functions/v1/courses.
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

// Update handles PATCH /functions/v1/courses/:id.
func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Update(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /functions/v1/courses/:id.
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
		if msg := err.Error(); msg == "title is required" || msg == "invalid status" || msg == "title cannot be empty" {
			return httputil.ErrorJSON(c, http.StatusBadRequest, msg)
		}
		return httputil.ErrorJSON(c, http.StatusInternalServerError, err.Error())
	}
}
