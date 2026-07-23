package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"taskflow/internal/auth"
	"taskflow/internal/database"
	"taskflow/internal/health"
	"taskflow/internal/project"
	"taskflow/internal/task"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	db, err := database.Open(
		context.Background(),
		databaseURL,
	)
	if err != nil {
		log.Fatalf(
			"failed to connect to database: %v",
			err,
		)
	}
	defer db.Close()

	log.Println("connected to postgresql")

	authRepository := auth.NewPostgresRepository(db)
	authService := auth.NewService(
		authRepository,
		jwtSecret,
	)
	authHandler := auth.NewHandler(authService)

	projectRepository := project.NewPostgresRepository(db)
	projectService := project.NewService(projectRepository)
	projectHandler := project.NewHandler(projectService)

	taskRepository := task.NewPostgresRepository(db)
	taskService := task.NewService(
		taskRepository,
		projectService,
	)
	taskHandler := task.NewHandler(taskService)

	healthHandler := health.NewHandler(db)

	router := http.NewServeMux()

	// Public health endpoints.
	router.HandleFunc(
		"GET /healthz",
		healthHandler.HealthCheck,
	)
	router.HandleFunc(
		"GET /readyz",
		healthHandler.ReadinessCheck,
	)

	// Public authentication endpoints.
	router.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)
	router.HandleFunc(
		"POST /api/v1/auth/login",
		authHandler.Login,
	)

	// Protected project endpoints.
	router.Handle(
		"POST /api/v1/projects",
		authService.RequireAuth(
			http.HandlerFunc(projectHandler.Create),
		),
	)
	router.Handle(
		"GET /api/v1/projects",
		authService.RequireAuth(
			http.HandlerFunc(projectHandler.FindAll),
		),
	)
	router.Handle(
		"GET /api/v1/projects/{id}",
		authService.RequireAuth(
			http.HandlerFunc(projectHandler.FindByID),
		),
	)
	router.Handle(
		"DELETE /api/v1/projects/{id}",
		authService.RequireAuth(
			http.HandlerFunc(projectHandler.Delete),
		),
	)

	// Protected task endpoints.
	router.Handle(
		"POST /api/v1/projects/{projectID}/tasks",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.Create),
		),
	)
	router.Handle(
		"GET /api/v1/projects/{projectID}/tasks",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.FindByProjectID),
		),
	)
	router.Handle(
		"GET /api/v1/tasks/{id}",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.FindByID),
		),
	)
	router.Handle(
		"PUT /api/v1/tasks/{id}",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.Update),
		),
	)
	router.Handle(
		"PATCH /api/v1/tasks/{id}/status",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.AdvanceStatus),
		),
	)
	router.Handle(
		"DELETE /api/v1/tasks/{id}",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.Delete),
		),
	)

	// Protected project summary endpoint.
	router.Handle(
		"GET /api/v1/projects/{id}/summary",
		authService.RequireAuth(
			http.HandlerFunc(taskHandler.GetProjectSummary),
		),
	)

	log.Println("TaskFlow API is running on :8080")

	if err := http.ListenAndServe(
		":8080",
		router,
	); err != nil {
		log.Fatal(err)
	}
}
