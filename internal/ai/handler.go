package ai

import (
	"errors"
	"net/http"

	"aiqadam-backend/internal/auth"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Handler exposes AI analytics HTTP endpoints.
type Handler struct {
	repo Repository
}

// NewHandler creates an AI handler.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

type analyzeRequest struct {
	EmployeeID     string `json:"employee_id"`
	QuizResultID   string `json:"quiz_result_id"`
	OrganizationID string `json:"organization_id"`
}

// Analyze handles POST /functions/v1/org-admin/ai/analyze
func (h *Handler) Analyze(c echo.Context) error {
	claims, err := auth.ClaimsFromContext(c.Request().Context())
	if err != nil {
		return aiError(c, http.StatusUnauthorized, "unauthorized")
	}
	if claims.Role != "org_admin" && claims.Role != "admin" && claims.Role != "super_admin" {
		return aiError(c, http.StatusForbidden, "forbidden")
	}

	var req analyzeRequest
	if err := c.Bind(&req); err != nil {
		return aiError(c, http.StatusBadRequest, "invalid request body")
	}
	if req.EmployeeID == "" || req.QuizResultID == "" || req.OrganizationID == "" {
		return aiError(c, http.StatusBadRequest, "employee_id, quiz_result_id and organization_id are required")
	}

	ctx := c.Request().Context()

	// Fetch real quiz data from DB
	analysis, err := h.repo.GetQuizResultForAnalysis(ctx, req.EmployeeID, req.QuizResultID)
	if errors.Is(err, pgx.ErrNoRows) {
		return aiError(c, http.StatusNotFound, "quiz result not found")
	}
	if err != nil {
		return aiError(c, http.StatusInternalServerError, err.Error())
	}

	// Call Gemini with dynamic data
	result, err := Analyze(*analysis)
	if err != nil {
		return aiError(c, http.StatusInternalServerError, "AI analysis failed: "+err.Error())
	}

	// Persist result
	if err := h.repo.SaveAnalysis(ctx,
		req.EmployeeID, req.OrganizationID, req.QuizResultID,
		analysis.CourseName, analysis.Score, result,
	); err != nil {
		return aiError(c, http.StatusInternalServerError, "save analysis: "+err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// GetOrgStats handles GET /functions/v1/org-admin/ai/stats
func (h *Handler) GetOrgStats(c echo.Context) error {
	claims, err := auth.ClaimsFromContext(c.Request().Context())
	if err != nil {
		return aiError(c, http.StatusUnauthorized, "unauthorized")
	}
	if claims.Role != "org_admin" && claims.Role != "admin" && claims.Role != "super_admin" {
		return aiError(c, http.StatusForbidden, "forbidden")
	}

	orgID := c.Param("organization_id")
	if orgID == "" {
		return aiError(c, http.StatusBadRequest, "organization_id is required")
	}

	stats, err := h.repo.GetOrgStats(c.Request().Context(), orgID)
	if err != nil {
		return aiError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

// AnalyzeBatch handles POST /functions/v1/org-admin/ai/analyze/batch/:organization_id
// Analyzes all employees in an org who have quiz results but no AI analysis yet.
func (h *Handler) AnalyzeBatch(c echo.Context) error {
	claims, err := auth.ClaimsFromContext(c.Request().Context())
	if err != nil {
		return aiError(c, http.StatusUnauthorized, "unauthorized")
	}
	if claims.Role != "org_admin" && claims.Role != "admin" && claims.Role != "super_admin" {
		return aiError(c, http.StatusForbidden, "forbidden")
	}

	orgID := c.Param("organization_id")
	if orgID == "" {
		return aiError(c, http.StatusBadRequest, "organization_id is required")
	}

	ctx := c.Request().Context()

	pending, err := h.repo.ListPendingAnalysis(ctx, orgID)
	if err != nil {
		return aiError(c, http.StatusInternalServerError, err.Error())
	}

	processed := 0
	failed := 0

	for _, item := range pending {
		analysis, err := h.repo.GetQuizResultForAnalysis(ctx, item.EmployeeID, item.QuizResultID)
		if err != nil {
			failed++
			continue
		}

		result, err := Analyze(*analysis)
		if err != nil {
			failed++
			continue
		}

		if err := h.repo.SaveAnalysis(ctx,
			item.EmployeeID, item.OrganizationID, item.QuizResultID,
			analysis.CourseName, analysis.Score, result,
		); err != nil {
			failed++
			continue
		}
		processed++
	}

	return c.JSON(http.StatusOK, map[string]any{
		"total":     len(pending),
		"processed": processed,
		"failed":    failed,
	})
}

func aiError(c echo.Context, status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}
