package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aiqadam-backend/internal/adminprofile"
	"aiqadam-backend/internal/ai"
	"aiqadam-backend/internal/assignments"
	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/briefings"
	"aiqadam-backend/internal/config"
	"aiqadam-backend/internal/contactrequests"
	"aiqadam-backend/internal/courses"
	"aiqadam-backend/internal/database"
	"aiqadam-backend/internal/employee"
	apphttp "aiqadam-backend/internal/http"
	"aiqadam-backend/internal/http/routes"
	"aiqadam-backend/internal/lessons"
	"aiqadam-backend/internal/materials"
	"aiqadam-backend/internal/organizations"
	"aiqadam-backend/internal/orgadmin"
	"aiqadam-backend/internal/quizzes"
	"aiqadam-backend/internal/storage"
	"aiqadam-backend/internal/users"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	db, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	tokens, err := auth.NewTokenManager(cfg.AuthJWTSecret, cfg.AuthJWTIssuer, cfg.AuthJWTExpiry)
	if err != nil {
		log.Fatalf("jwt config: %v", err)
	}

	pool := db.Pool()

	authRepo := auth.NewRepository(pool)
	authService := auth.NewService(authRepo, tokens)
	authHandler := auth.NewHandler(authService)

	usersRepo := users.NewRepository(pool)
	usersService := users.NewService(usersRepo)
	adminUsersHandler := users.NewHandler(usersService)

	orgsRepo := organizations.NewRepository(pool)
	orgsService := organizations.NewService(orgsRepo)
	orgsHandler := organizations.NewHandler(orgsService)

	coursesRepo := courses.NewRepository(pool)
	coursesService := courses.NewService(coursesRepo)
	coursesHandler := courses.NewHandler(coursesService)

	fileStorage, err := storage.NewLocal(cfg.StoragePath, cfg.PublicAPIURL)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	lessonsRepo := lessons.NewRepository(pool)
	lessonsService := lessons.NewService(lessonsRepo, fileStorage)
	lessonsHandler := lessons.NewHandler(lessonsService)

	materialsRepo := materials.NewRepository(pool)
	materialsService := materials.NewService(materialsRepo, fileStorage)
	materialsHandler := materials.NewHandler(materialsService)

	assignmentsRepo := assignments.NewRepository(pool)
	assignmentsService := assignments.NewService(assignmentsRepo)
	assignmentsHandler := assignments.NewHandler(assignmentsService)

	employeeRepo := employee.NewRepository(pool)
	employeeService := employee.NewService(employeeRepo)
	employeeHandler := employee.NewHandler(employeeService)

	aiRepo := ai.NewRepository(pool)
	aiHandler := ai.NewHandler(aiRepo)
	employeeService.SetAIRepository(aiRepo)

	orgAdminRepo := orgadmin.NewRepository(pool)
	orgAdminService := orgadmin.NewService(orgAdminRepo)
	orgAdminHandler := orgadmin.NewHandler(orgAdminService)

	briefingsRepo := briefings.NewRepository(pool)
	briefingsService := briefings.NewService(briefingsRepo, fileStorage)
	briefingsHandler := briefings.NewHandler(briefingsService)

	contactRequestsRepo := contactrequests.NewRepository(pool)
	contactRequestsService := contactrequests.NewService(contactRequestsRepo)
	contactRequestsHandler := contactrequests.NewHandler(contactRequestsService)

	adminProfileRepo := adminprofile.NewRepository(pool)
	adminProfileService := adminprofile.NewService(adminProfileRepo, fileStorage)
	adminProfileHandler := adminprofile.NewHandler(adminProfileService)

	quizzesRepo := quizzes.NewRepository(pool)
	quizzesService := quizzes.NewService(quizzesRepo)
	quizzesHandler := quizzes.NewHandler(quizzesService)

	server := apphttp.NewServer(cfg, routes.Deps{
		Health:        db,
		Auth:          authHandler,
		AuthService:   authService,
		AdminUsers:    adminUsersHandler,
		Organizations: orgsHandler,
		Courses:       coursesHandler,
		Lessons:       lessonsHandler,
		Materials:     materialsHandler,
		Assignments:   assignmentsHandler,
		Employee:      employeeHandler,
		OrgAdmin:      orgAdminHandler,
		AdminProfile:  adminProfileHandler,
		Quizzes:       quizzesHandler,
		Briefings:       briefingsHandler,
		ContactRequests: contactRequestsHandler,
		AI:              aiHandler,
	})

	go func() {
		log.Printf("api server listening on http://localhost:%s", cfg.Port)
		if err := server.Start(); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
