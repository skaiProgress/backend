package organizations

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes organization HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an organizations handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /functions/v1/admin/organizations.
func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context(), c.QueryParam("search"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Get handles GET /functions/v1/admin/organizations/:id.
func (h *Handler) Get(c echo.Context) error {
	item, err := h.service.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// Create handles POST /functions/v1/admin/organizations.
func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /functions/v1/admin/organizations/:id.
func (h *Handler) Update(c echo.Context) error {
	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Update(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /functions/v1/admin/organizations/:id.
func (h *Handler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// AddMember handles POST /functions/v1/admin/organizations/:id/users.
func (h *Handler) AddMember(c echo.Context) error {
	var req AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	res, err := h.service.AddMember(c.Request().Context(), c.Param("id"), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, res)
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
	case errors.Is(err, ErrEmailExists):
		return errorJSON(c, http.StatusBadRequest, "User already registered")
	case errors.Is(err, ErrInvalidInput):
		return errorJSON(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, pgx.ErrNoRows):
		return errorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		switch msg {
		case "name is required", "name cannot be empty",
			"email and password are required",
			"role must be user or org_admin":
			return errorJSON(c, http.StatusBadRequest, msg)
		case "organization with this BIN already exists":
			return errorJSON(c, http.StatusBadRequest, msg)
		default:
			return errorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
