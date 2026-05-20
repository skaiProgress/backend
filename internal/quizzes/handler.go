package quizzes

import (
	"errors"
	"net/http"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes quiz HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a quizzes handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAdmin handles GET /functions/v1/lessons/:id/quiz
func (h *Handler) GetAdmin(c echo.Context) error {
	item, err := h.service.GetAdmin(c.Request().Context(), c.Param("id"))
	if err != nil {
		return mapError(c, err)
	}
	if item == nil {
		return c.JSON(http.StatusOK, nil)
	}
	return c.JSON(http.StatusOK, item)
}

// Upload handles POST /functions/v1/lessons/:id/quiz
func (h *Handler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "file is required")
	}
	src, err := file.Open()
	if err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "cannot read file")
	}
	defer src.Close()

	out, err := h.service.Upload(c.Request().Context(), c.Param("id"), file.Filename, src)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// Delete handles DELETE /functions/v1/lessons/:id/quiz
func (h *Handler) Delete(c echo.Context) error {
	if err := h.service.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetEmployee handles GET /functions/v1/employee/lessons/:lessonId/quiz
func (h *Handler) GetEmployee(c echo.Context) error {
	item, err := h.service.GetEmployee(c.Request().Context(), c.Param("lessonId"))
	if err != nil {
		return mapError(c, err)
	}
	if item == nil {
		return httputil.ErrorJSON(c, http.StatusNotFound, "quiz not found")
	}
	return c.JSON(http.StatusOK, item)
}

// SubmitEmployee handles POST /functions/v1/employee/lessons/:lessonId/quiz/submit
func (h *Handler) SubmitEmployee(c echo.Context) error {
	var req SubmitRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.Submit(c.Request().Context(), c.Param("lessonId"), req)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return httputil.ErrorJSON(c, http.StatusUnauthorized, "Unauthorized")
	case errors.Is(err, auth.ErrForbidden):
		return httputil.ErrorJSON(c, http.StatusForbidden, "Forbidden")
	case errors.Is(err, pgx.ErrNoRows):
		return httputil.ErrorJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		if strings.Contains(msg, "нужно") ||
			strings.Contains(msg, "вопрос") ||
			strings.Contains(msg, "файл") ||
			strings.Contains(msg, "lesson_id") ||
			strings.Contains(msg, "question_id") ||
			strings.Contains(msg, "ANSWER") ||
			strings.Contains(msg, "вариант") ||
			strings.Contains(msg, "повторная") {
			return httputil.ErrorJSON(c, http.StatusBadRequest, msg)
		}
		return httputil.ErrorJSON(c, http.StatusInternalServerError, msg)
	}
}
