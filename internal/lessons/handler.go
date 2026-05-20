package lessons

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/pkg/httputil"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes lesson HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a lessons handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /functions/v1/lessons?course_id=
func (h *Handler) List(c echo.Context) error {
	items, err := h.service.List(c.Request().Context(), c.QueryParam("course_id"))
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, items)
}

// Reorder handles PATCH /functions/v1/lessons/reorder.
func (h *Handler) Reorder(c echo.Context) error {
	var req ReorderRequest
	if err := c.Bind(&req); err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.service.Reorder(c.Request().Context(), req); err != nil {
		return mapError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Create handles POST /functions/v1/lessons (JSON or multipart).
func (h *Handler) Create(c echo.Context) error {
	if strings.Contains(c.Request().Header.Get("Content-Type"), "multipart/form-data") {
		return h.createUpload(c)
	}

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

func (h *Handler) createUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return httputil.ErrorJSON(c, http.StatusBadRequest, "video file is required")
	}

	req := CreateUploadRequest{
		CourseID: c.FormValue("course_id"),
		Title:    c.FormValue("title"),
	}
	if desc := strings.TrimSpace(c.FormValue("description")); desc != "" {
		req.Description = &desc
	}
	if orderStr := c.FormValue("order_index"); orderStr != "" {
		if n, err := strconv.Atoi(orderStr); err == nil {
			req.OrderIndex = n
		}
	}
	if isFreeStr := c.FormValue("is_free"); isFreeStr == "true" || isFreeStr == "1" {
		req.IsFree = true
	}

	out, err := h.service.CreateFromMultipart(c.Request().Context(), req, file)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// Update handles PATCH /functions/v1/lessons/:id (JSON or multipart).
func (h *Handler) Update(c echo.Context) error {
	if strings.Contains(c.Request().Header.Get("Content-Type"), "multipart/form-data") {
		return h.updateUpload(c)
	}

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

func (h *Handler) updateUpload(c echo.Context) error {
	var req UpdateUploadRequest
	if title := strings.TrimSpace(c.FormValue("title")); title != "" {
		req.Title = &title
	}
	if desc := c.FormValue("description"); desc != "" {
		req.Description = &desc
	}
	if orderStr := c.FormValue("order_index"); orderStr != "" {
		if n, err := strconv.Atoi(orderStr); err == nil {
			req.OrderIndex = &n
		}
	}
	if isFreeStr := c.FormValue("is_free"); isFreeStr != "" {
		v := isFreeStr == "true" || isFreeStr == "1"
		req.IsFree = &v
	}

	var fh *multipart.FileHeader
	if file, err := c.FormFile("file"); err == nil {
		fh = file
	}

	out, err := h.service.UpdateFromMultipart(c.Request().Context(), c.Param("id"), req, fh)
	if err != nil {
		return mapError(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// Delete handles DELETE /functions/v1/lessons/:id.
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
		msg := err.Error()
		switch msg {
		case "invalid youtube url",
			"course_id, title and youtube_url are required",
			"course_id, title and video file are required",
			"video file is required",
			"title cannot be empty",
			"course_id is required",
			"course_id and ordered_ids are required",
			"unsupported video format (allowed: mp4, webm, mov, m4v)",
			"invalid video file size":
			return httputil.ErrorJSON(c, http.StatusBadRequest, msg)
		default:
			return httputil.ErrorJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
