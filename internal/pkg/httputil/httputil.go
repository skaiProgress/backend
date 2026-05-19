package httputil

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorJSON writes a JSON error response compatible with legacy edge functions.
func ErrorJSON(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

// RequireAdminRole checks JWT claims role from context (set by auth middleware).
func RequireAdminRole(c echo.Context, role string) bool {
	return role == "admin" || role == "super_admin"
}

// AbortForbidden returns 403 if not admin.
func AbortForbidden(c echo.Context, role string) error {
	if RequireAdminRole(c, role) {
		return nil
	}
	return ErrorJSON(c, http.StatusForbidden, "Forbidden: requires admin or super_admin role")
}
