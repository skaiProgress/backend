package users

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes admin user HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an admin users handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListUsers handles GET /functions/v1/admin/users.
func (h *Handler) ListUsers(c echo.Context) error {
	items, err := h.service.ListUsers(c.Request().Context(), c.QueryParam("search"))
	if err != nil {
		return mapServiceError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// GetUser handles GET /functions/v1/admin/users/:id.
func (h *Handler) GetUser(c echo.Context) error {
	item, err := h.service.GetUser(c.Request().Context(), c.Param("id"))
	if err != nil {
		return mapServiceError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// AddUser handles POST /functions/v1/admin-add-user.
func (h *Handler) AddUser(c echo.Context) error {
	var req AddUserRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}

	res, err := h.service.AddUser(c.Request().Context(), req)
	if err != nil {
		return mapServiceError(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateUser handles POST /functions/v1/admin-update-user.
func (h *Handler) UpdateUser(c echo.Context) error {
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}

	if err := h.service.UpdateUser(c.Request().Context(), req); err != nil {
		return mapServiceError(c, err)
	}
	return c.JSON(http.StatusOK, OkResponse{OK: true})
}

// DeleteUser handles POST /functions/v1/admin-delete-user.
func (h *Handler) DeleteUser(c echo.Context) error {
	var req DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}

	res, err := h.service.DeleteUsers(c.Request().Context(), req)
	if err != nil {
		return mapServiceError(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func errorJSON(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

func mapServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return errorJSON(c, http.StatusUnauthorized, "Invalid or expired token. Please sign in again.")
	case errors.Is(err, ErrForbidden):
		return errorJSON(c, http.StatusForbidden, "Forbidden: requires admin or super_admin role")
	case errors.Is(err, ErrCannotDeleteSelf):
		return errorJSON(c, http.StatusForbidden, "Нельзя удалить собственный аккаунт")
	case errors.Is(err, ErrEmailExists):
		return errorJSON(c, http.StatusBadRequest, "User already registered")
	case errors.Is(err, ErrInvalidInput):
		return errorJSON(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, pgx.ErrNoRows):
		return errorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		switch msg {
		case "only super_admin can create admin users",
			"only super_admin can grant super_admin role",
			"only super_admin can create organization users",
			"organization users must have role user or org_admin",
			"organization not found":
			return errorJSON(c, http.StatusForbidden, msg)
		case "email and password are required",
			"user_id is required",
			"user_ids must be a non-empty array",
			"password must be at least 6 characters":
			return errorJSON(c, http.StatusBadRequest, msg)
		default:
			if errors.Is(err, ErrInvalidInput) {
				return errorJSON(c, http.StatusBadRequest, msg)
			}
			return errorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
