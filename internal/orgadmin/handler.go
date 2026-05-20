package orgadmin

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes org-admin HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an org-admin handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetStats(c echo.Context) error {
	out, err := h.service.GetStats(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) GetProfile(c echo.Context) error {
	out, err := h.service.GetProfile(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) ListMembers(c echo.Context) error {
	out, err := h.service.ListMembers(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) CreateMember(c echo.Context) error {
	var req CreateMemberRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.CreateMember(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) ListCourses(c echo.Context) error {
	out, err := h.service.ListCourses(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) GetCourse(c echo.Context) error {
	out, err := h.service.GetCourse(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) ListAssignments(c echo.Context) error {
	out, err := h.service.ListAssignments(c.Request().Context(), c.QueryParam("course_id"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) CreateAssignment(c echo.Context) error {
	var req CreateAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return errorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.CreateAssignment(c.Request().Context(), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) RevokeAssignment(c echo.Context) error {
	if err := h.service.RevokeAssignment(c.Request().Context(), c.Param("id")); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func errorJSON(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrUnauthorized):
		return errorJSON(c, http.StatusUnauthorized, "Invalid or expired token")
	case errors.Is(err, ErrForbidden):
		return errorJSON(c, http.StatusForbidden, "Forbidden: org-admin access required")
	case errors.Is(err, ErrEmailExists):
		return errorJSON(c, http.StatusBadRequest, "User already registered")
	case errors.Is(err, ErrInvalidInput):
		return errorJSON(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, pgx.ErrNoRows):
		return errorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		switch msg {
		case "email and password are required",
			"password must be at least 6 characters",
			"user_id and course_id are required",
			"user is not an employee of your organization",
			"course is not assigned to you",
			"invalid expires_at":
			return errorJSON(c, http.StatusBadRequest, msg)
		default:
			return errorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
