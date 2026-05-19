package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aiqadam-backend/internal/adminprofile"
	"aiqadam-backend/internal/assignments"
	"aiqadam-backend/internal/auth"
	"aiqadam-backend/internal/config"
	"aiqadam-backend/internal/courses"
	"aiqadam-backend/internal/database"
	"aiqadam-backend/internal/employee"
	apphttp "aiqadam-backend/internal/http"
	"aiqadam-backend/internal/http/routes"
	"aiqadam-backend/internal/lessons"
	"aiqadam-backend/internal/materials"
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

	coursesRepo := courses.NewRepository(pool)
	coursesService := courses.NewService(coursesRepo)
	coursesHandler := courses.NewHandler(coursesService)

	lessonsRepo := lessons.NewRepository(pool)
	lessonsService := lessons.NewService(lessonsRepo)
	lessonsHandler := lessons.NewHandler(lessonsService)

	fileStorage, err := storage.NewLocal(cfg.StoragePath, cfg.PublicAPIURL)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	materialsRepo := materials.NewRepository(pool)
	materialsService := materials.NewService(materialsRepo, fileStorage)
	materialsHandler := materials.NewHandler(materialsService)

	assignmentsRepo := assignments.NewRepository(pool)
	assignmentsService := assignments.NewService(assignmentsRepo)
	assignmentsHandler := assignments.NewHandler(assignmentsService)

	employeeRepo := employee.NewRepository(pool)
	employeeService := employee.NewService(employeeRepo)
	employeeHandler := employee.NewHandler(employeeService)

	adminProfileRepo := adminprofile.NewRepository(pool)
	adminProfileService := adminprofile.NewService(adminProfileRepo, fileStorage)
	adminProfileHandler := adminprofile.NewHandler(adminProfileService)

	server := apphttp.NewServer(cfg, routes.Deps{
		Health:      db,
		Auth:        authHandler,
		AuthService: authService,
		AdminUsers:  adminUsersHandler,
		Courses:     coursesHandler,
		Lessons:     lessonsHandler,
		Materials:   materialsHandler,
		Assignments: assignmentsHandler,
		Employee:     employeeHandler,
		AdminProfile: adminProfileHandler,
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
