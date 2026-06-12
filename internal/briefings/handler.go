package briefings

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// ── Admin: briefing videos ────────────────────────────────────────────────────

// ListBriefingVideos GET /functions/v1/courses/:courseId/briefing-videos
func (h *Handler) ListBriefingVideos(c echo.Context) error {
	out, err := h.service.ListBriefingVideos(c.Request().Context(), c.Param("courseId"))
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// UploadBriefingVideo POST /functions/v1/courses/:courseId/briefing-videos (multipart)
func (h *Handler) UploadBriefingVideo(c echo.Context) error {
	kind := c.FormValue("briefing_kind")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errJSON(c, http.StatusBadRequest, "file is required")
	}
	src, err := fileHeader.Open()
	if err != nil {
		return errJSON(c, http.StatusBadRequest, "cannot open uploaded file")
	}
	defer src.Close()
	if err := h.service.UploadBriefingVideo(c.Request().Context(), c.Param("courseId"), kind, fileHeader.Filename, src); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusCreated)
}

// DeleteBriefingVideo DELETE /functions/v1/courses/:courseId/briefing-videos/:kind
func (h *Handler) DeleteBriefingVideo(c echo.Context) error {
	if err := h.service.DeleteBriefingVideo(c.Request().Context(), c.Param("courseId"), c.Param("kind")); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListBriefingCourses GET /functions/v1/org-admin/briefing-courses
func (h *Handler) ListBriefingCourses(c echo.Context) error {
	out, err := h.service.ListBriefingCourses(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// Handler exposes briefing HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a briefings handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ── Org-Admin endpoints ──────────────────────────────────────────────────────

// ListOrgEvents GET /functions/v1/org-admin/events
func (h *Handler) ListOrgEvents(c echo.Context) error {
	out, err := h.service.ListOrgAdminEvents(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// CreateEvent POST /functions/v1/org-admin/events
func (h *Handler) CreateEvent(c echo.Context) error {
	var req CreateBriefingEventRequest
	if err := c.Bind(&req); err != nil {
		return errJSON(c, http.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.CreateManualBriefingEvent(c.Request().Context(), req)
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusCreated, out)
}

// UpdateEvent PATCH /functions/v1/org-admin/events/:id
func (h *Handler) UpdateEvent(c echo.Context) error {
	var req UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return errJSON(c, http.StatusBadRequest, "invalid request body")
	}
	if err := h.service.UpdateEventTime(c.Request().Context(), c.Param("id"), req.StartsAt); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListOrgRecords GET /functions/v1/org-admin/briefing-records
func (h *Handler) ListOrgRecords(c echo.Context) error {
	out, err := h.service.ListOrgAdminRecords(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// SignRecord PATCH /functions/v1/org-admin/briefing-records/:id/sign
func (h *Handler) SignRecord(c echo.Context) error {
	if err := h.service.InstructorSign(c.Request().Context(), c.Param("id")); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteRecord DELETE /functions/v1/org-admin/briefing-records/:id
func (h *Handler) DeleteRecord(c echo.Context) error {
	if err := h.service.DeleteOrgRecord(c.Request().Context(), c.Param("id")); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ── Employee endpoints ───────────────────────────────────────────────────────

// ListEmployeeEvents GET /functions/v1/employee/events
func (h *Handler) ListEmployeeEvents(c echo.Context) error {
	out, err := h.service.ListEmployeeEvents(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// ListEmployeeBriefings GET /functions/v1/employee/briefings
func (h *Handler) ListEmployeeBriefings(c echo.Context) error {
	out, err := h.service.ListEmployeeBriefings(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// GetEmployeeBriefing GET /functions/v1/employee/briefings/:eventId
func (h *Handler) GetEmployeeBriefing(c echo.Context) error {
	out, err := h.service.GetEmployeeBriefingDetail(c.Request().Context(), c.Param("eventId"))
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// CompleteBriefing POST /functions/v1/employee/briefings/:eventId/complete
func (h *Handler) CompleteBriefing(c echo.Context) error {
	if err := h.service.CompleteBriefing(c.Request().Context(), c.Param("eventId")); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListEmployeeJournalRecords GET /functions/v1/employee/journal-records
func (h *Handler) ListEmployeeJournalRecords(c echo.Context) error {
	out, err := h.service.ListEmployeeJournalRecords(c.Request().Context())
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusOK, out)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func errJSON(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}

func mapErr(c echo.Context, err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errJSON(c, http.StatusNotFound, "not found")
	default:
		msg := err.Error()
		switch msg {
		case "unauthorized":
			return errJSON(c, http.StatusUnauthorized, msg)
		case "forbidden":
			return errJSON(c, http.StatusForbidden, msg)
		case "briefing already confirmed",
			"event has no briefing kind",
			"invalid starts_at, use RFC3339 format",
			"invalid ends_at, use RFC3339 format",
			"ends_at must be after starts_at",
			"record not found or already signed",
			"employee_id is required",
			"course_id is required",
			"invalid briefing_kind",
			"course not available",
			"course is not a briefing course",
			"no video for this briefing kind",
			"briefing window not started",
			"briefing window expired",
			"employee not in organization":
			return errJSON(c, http.StatusBadRequest, msg)
		case "event not found",
			"employee profile not found",
			"org-admin profile not found":
			return errJSON(c, http.StatusNotFound, msg)
		default:
			return errJSON(c, http.StatusInternalServerError, msg)
		}
	}
}
