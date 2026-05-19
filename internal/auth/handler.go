package auth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler exposes HTTP handlers for authentication.
type Handler struct {
	service *Service
}

// NewHandler creates a new auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Login handles POST /auth/login.
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	result, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid login credentials")
		case errors.Is(err, ErrUserBanned):
			return echo.NewHTTPError(http.StatusForbidden, "Account is disabled")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, result)
}

// Me handles GET /functions/v1/auth/me.
func (h *Handler) Me(c echo.Context) error {
	result, err := h.service.Me(c.Request().Context())
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

// ChangePassword handles POST /functions/v1/auth/change-password.
func (h *Handler) ChangePassword(c echo.Context) error {
	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err := h.service.ChangePassword(c.Request().Context(), req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		case errors.Is(err, ErrWrongPassword):
			return echo.NewHTTPError(http.StatusBadRequest, "Неверный текущий пароль")
		default:
			msg := err.Error()
			if msg == "current_password and new_password are required" ||
				msg == "password must be at least 8 characters" ||
				msg == "password must contain at least one digit" {
				return echo.NewHTTPError(http.StatusBadRequest, msg)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}
	}

	return c.NoContent(http.StatusNoContent)
}
