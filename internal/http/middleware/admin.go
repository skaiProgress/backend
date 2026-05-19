package middleware

import (
	"net/http"

	"aiqadam-backend/internal/auth"

	"github.com/labstack/echo/v4"
)

// RequireAdmin ensures the JWT bearer has admin or super_admin role.
// Must run after JWT middleware.
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := auth.ClaimsFromContext(c.Request().Context())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token. Please sign in again.",
				})
			}
			if claims.Role != "admin" && claims.Role != "super_admin" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Forbidden: requires admin or super_admin role",
				})
			}
			return next(c)
		}
	}
}
