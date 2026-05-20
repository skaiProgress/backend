package employee

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes employee HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates an employee handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListCourses handles GET /functions/v1/employee/courses.
func (h *Handler) ListCourses(c echo.Context) error {
	items, err := h.service.ListCourses(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// GetCourse handles GET /functions/v1/employee/courses/:courseId.
func (h *Handler) GetCourse(c echo.Context) error {
	item, err := h.service.GetCourseDetail(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	if item == nil {
		return httputil.ErrorJSON(c, http.StatusNotFound, "course not found or access denied")
	}
	return c.JSON(http.StatusOK, item)
}

// ListLessons handles GET /functions/v1/employee/courses/:courseId/lessons.
func (h *Handler) ListLessons(c echo.Context) error {
	items, err := h.service.ListLessons(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// ListMaterials handles GET /functions/v1/employee/courses/:courseId/materials.
func (h *Handler) ListMaterials(c echo.Context) error {
	items, err := h.service.ListMaterials(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// GetProfile handles GET /functions/v1/employee/profile.
func (h *Handler) GetProfile(c echo.Context) error {
	item, err := h.service.GetProfile(c.Request().Context())
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, item)
}

// GetCourseProgress handles GET /functions/v1/employee/courses/:courseId/progress
func (h *Handler) GetCourseProgress(c echo.Context) error {
	item, err := h.service.GetCourseProgress(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	if item == nil {
		return httputil.ErrorJSON(c, http.StatusNotFound, "course not found or access denied")
	}
	return c.JSON(http.StatusOK, item)
}

// CompleteTraining handles POST /functions/v1/employee/courses/:courseId/complete-training
func (h *Handler) CompleteTraining(c echo.Context) error {
	out, err := h.service.CompleteTraining(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// UpdateProfile handles PATCH /functions/v1/employee/profile.
func (h *Handler) UpdateProfile(c echo.Context) error {
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.service.UpdateProfile(c.Request().Context(), req.FullName); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return httputil.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
	case errors.Is(err, auth.ErrForbidden):
		return httputil.ErrorJSON(c, http.StatusForbidden, "Account is disabled")
	case errors.Is(err, pgx.ErrNoRows):
		return httputil.ErrorJSON(c, http.StatusNotFound, "not found")
	default:
		switch err.Error() {
		case "course_id is required",
			"сдайте все тесты по урокам, чтобы завершить обучение",
			"AI module is not configured":
			return httputil.ErrorJSON(c, http.StatusBadRequest, err.Error())
		}
		return httputil.ErrorJSON(c, http.StatusInternalServerError, err.Error())
	}
}
