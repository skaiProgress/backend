package routes

import (
	"context"
	"net/http"
	"time"

	"aiqadam-backend/internal/adminprofile"
	"aiqadam-backend/internal/ai"
	"aiqadam-backend/internal/assignments"
	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/briefings"
	"aiqadam-backend/internal/courses"
	"aiqadam-backend/internal/employee"
	appmiddleware "aiqadam-backend/internal/http/middleware"
	"aiqadam-backend/internal/lessons"
	"aiqadam-backend/internal/materials"
	"aiqadam-backend/internal/organizations"
	"aiqadam-backend/internal/orgadmin"
	"aiqadam-backend/internal/quizzes"
	"aiqadam-backend/internal/users"

	"github.com/labstack/echo/v4"
)

// HealthChecker verifies database connectivity.
type HealthChecker interface {
	Ping(ctx context.Context) error
}

// Deps holds dependencies required to register HTTP routes.
type Deps struct {
	Health        HealthChecker
	Auth          *auth.Handler
	AuthService   *auth.Service
	AdminUsers    *users.Handler
	Organizations *organizations.Handler
	Courses       *courses.Handler
	Lessons       *lessons.Handler
	Materials     *materials.Handler
	Assignments   *assignments.Handler
	Employee      *employee.Handler
	OrgAdmin      *orgadmin.Handler
	AdminProfile  *adminprofile.Handler
	Quizzes       *quizzes.Handler
	Briefings     *briefings.Handler
	AI            *ai.Handler
}

// Register mounts all application routes on the Echo instance.
func Register(e *echo.Echo, deps Deps) {
	e.GET("/", rootHandler)
	e.GET("/healthz", healthHandler(deps.Health))

	e.POST("/auth/login", deps.Auth.Login)

	authGroup := e.Group("/functions/v1/auth")
	authGroup.GET("/me", deps.Auth.Me, appmiddleware.JWT(deps.AuthService))
	authGroup.POST("/change-password", deps.Auth.ChangePassword, appmiddleware.JWT(deps.AuthService))

	adminFn := e.Group("/functions/v1")
	adminFn.Use(appmiddleware.JWT(deps.AuthService))
	adminFn.Use(appmiddleware.RequireAdmin())

	adminFn.GET("/admin/profile", deps.AdminProfile.Get)
	adminFn.PATCH("/admin/profile", deps.AdminProfile.Update)
	adminFn.POST("/admin/profile/avatar", deps.AdminProfile.UploadAvatar)

	adminFn.GET("/admin/users", deps.AdminUsers.ListUsers)
	adminFn.GET("/admin/users/:id", deps.AdminUsers.GetUser)
	adminFn.POST("/admin-add-user", deps.AdminUsers.AddUser)
	adminFn.POST("/admin-update-user", deps.AdminUsers.UpdateUser)
	adminFn.POST("/admin-delete-user", deps.AdminUsers.DeleteUser)

	adminFn.GET("/admin/organizations", deps.Organizations.List)
	adminFn.GET("/admin/organizations/:id", deps.Organizations.Get)
	adminFn.POST("/admin/organizations", deps.Organizations.Create)
	adminFn.PATCH("/admin/organizations/:id", deps.Organizations.Update)
	adminFn.DELETE("/admin/organizations/:id", deps.Organizations.Delete)
	adminFn.POST("/admin/organizations/:id/users", deps.Organizations.AddMember)

	adminFn.GET("/courses", deps.Courses.List)
	adminFn.POST("/courses", deps.Courses.Create)
	adminFn.PATCH("/courses/:id", deps.Courses.Update)
	adminFn.DELETE("/courses/:id", deps.Courses.Delete)

	adminFn.GET("/lessons", deps.Lessons.List)
	adminFn.PATCH("/lessons/reorder", deps.Lessons.Reorder)
	adminFn.POST("/lessons", deps.Lessons.Create)
	adminFn.PATCH("/lessons/:id", deps.Lessons.Update)
	adminFn.DELETE("/lessons/:id", deps.Lessons.Delete)

	adminFn.GET("/lessons/:id/quiz", deps.Quizzes.GetAdmin)
	adminFn.POST("/lessons/:id/quiz", deps.Quizzes.Upload)
	adminFn.DELETE("/lessons/:id/quiz", deps.Quizzes.Delete)

	adminFn.GET("/materials", deps.Materials.List)
	adminFn.POST("/materials", deps.Materials.Upload)
	adminFn.DELETE("/materials/:id", deps.Materials.Delete)

	adminFn.GET("/assignments", deps.Assignments.List)
	adminFn.POST("/assignments", deps.Assignments.Create)
	adminFn.POST("/assignments/bulk", deps.Assignments.Bulk)
	adminFn.DELETE("/assignments/:id", deps.Assignments.Revoke)

	employeeFn := e.Group("/functions/v1/employee")
	employeeFn.Use(appmiddleware.JWT(deps.AuthService))
	employeeFn.GET("/courses", deps.Employee.ListCourses)
	employeeFn.GET("/courses/:courseId", deps.Employee.GetCourse)
	employeeFn.GET("/courses/:courseId/lessons", deps.Employee.ListLessons)
	employeeFn.GET("/courses/:courseId/progress", deps.Employee.GetCourseProgress)
	employeeFn.POST("/courses/:courseId/complete-training", deps.Employee.CompleteTraining)
	employeeFn.GET("/lessons/:lessonId/quiz", deps.Quizzes.GetEmployee)
	employeeFn.POST("/lessons/:lessonId/quiz/submit", deps.Quizzes.SubmitEmployee)
	employeeFn.GET("/courses/:courseId/materials", deps.Employee.ListMaterials)
	employeeFn.GET("/profile", deps.Employee.GetProfile)
	employeeFn.PATCH("/profile", deps.Employee.UpdateProfile)

	orgAdminFn := e.Group("/functions/v1/org-admin")
	orgAdminFn.Use(appmiddleware.JWT(deps.AuthService))
	orgAdminFn.GET("/stats", deps.OrgAdmin.GetStats)
	orgAdminFn.GET("/profile", deps.OrgAdmin.GetProfile)
	orgAdminFn.GET("/members", deps.OrgAdmin.ListMembers)
	orgAdminFn.POST("/members", deps.OrgAdmin.CreateMember)
	orgAdminFn.GET("/courses", deps.OrgAdmin.ListCourses)
	orgAdminFn.GET("/courses/:courseId", deps.OrgAdmin.GetCourse)
	orgAdminFn.GET("/assignments", deps.OrgAdmin.ListAssignments)
	orgAdminFn.POST("/assignments", deps.OrgAdmin.CreateAssignment)
	orgAdminFn.DELETE("/assignments/:id", deps.OrgAdmin.RevokeAssignment)
	// Briefing calendar & journal
	orgAdminFn.GET("/events", deps.Briefings.ListOrgEvents)
	orgAdminFn.POST("/events", deps.Briefings.CreateEvent)
	orgAdminFn.PATCH("/events/:id", deps.Briefings.UpdateEvent)
	orgAdminFn.GET("/briefing-records", deps.Briefings.ListOrgRecords)
	orgAdminFn.PATCH("/briefing-records/:id/sign", deps.Briefings.SignRecord)
	orgAdminFn.DELETE("/briefing-records/:id", deps.Briefings.DeleteRecord)

	// Employee briefings
	employeeFn.GET("/events", deps.Briefings.ListEmployeeEvents)
	employeeFn.GET("/briefings", deps.Briefings.ListEmployeeBriefings)
	employeeFn.POST("/briefings/:eventId/confirm", deps.Briefings.ConfirmBriefing)
	employeeFn.GET("/journal-records", deps.Briefings.ListEmployeeJournalRecords)

	// AI analytics (org-admin scope)
	orgAdminFn.POST("/ai/analyze", deps.AI.Analyze)
	orgAdminFn.GET("/ai/stats/:organization_id", deps.AI.GetOrgStats)
	orgAdminFn.POST("/ai/analyze/batch/:organization_id", deps.AI.AnalyzeBatch)
}

func rootHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"name":   "AIQADAM Backend",
		"status": "running",
	})
}

func healthHandler(health HealthChecker) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
		defer cancel()

		dbStatus := "ok"
		httpStatus := http.StatusOK
		overall := "ok"

		if err := health.Ping(ctx); err != nil {
			dbStatus = "error"
			overall = "error"
			httpStatus = http.StatusServiceUnavailable
		}

		return c.JSON(httpStatus, map[string]string{
			"status":   overall,
			"database": dbStatus,
		})
	}
}
