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
		log.Fatal("DATABASE URL IS REQUORED")
	}

	db, err := database.Open(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to postgresql")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	authRepository := auth.NewPostgresRepository(db)
	authService := auth.NewService(authRepository, jwtSecret)
	authHandler := auth.NewHandler(authService)

	projectRepository := project.NewPostgresRepository(db)
	projectService := project.NewService(projectRepository)
	projectHandler := project.NewHandler(projectService)

	taskRepository := task.NewPostgresRepository(db)
	taskService := task.NewService(taskRepository, projectService)
	taskHandler := task.NewHandler(taskService)

	router := http.NewServeMux()

	healthHandler := health.NewHandler(db)

	router.HandleFunc("GET /healthz", healthHandler.HealthCheck)
	router.HandleFunc("GET /readyz", healthHandler.ReadinessCheck)

	router.HandleFunc(
		"POST /api/v1/projects",
		projectHandler.Create,
	)
	router.HandleFunc(
		"GET /api/v1/projects",
		projectHandler.FindAll,
	)
	router.HandleFunc(
		"GET /api/v1/projects/{id}",
		projectHandler.FindByID,
	)
	router.HandleFunc(
		"DELETE /api/v1/projects/{id}",
		projectHandler.Delete,
	)

	router.HandleFunc(
		"POST /api/v1/projects/{projectID}/tasks",
		taskHandler.Create,
	)
	router.HandleFunc(
		"GET /api/v1/projects/{projectID}/tasks",
		taskHandler.FindByProjectID,
	)
	router.HandleFunc(
		"GET /api/v1/tasks/{id}",
		taskHandler.FindByID,
	)
	router.HandleFunc(
		"PUT /api/v1/tasks/{id}",
		taskHandler.Update,
	)
	router.HandleFunc(
		"PATCH /api/v1/tasks/{id}/status",
		taskHandler.AdvanceStatus,
	)
	router.HandleFunc(
		"DELETE /api/v1/tasks/{id}",
		taskHandler.Delete,
	)

	router.HandleFunc(
		"GET /api/v1/projects/{id}/summary",
		taskHandler.GetProjectSummary,
	)

	router.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)

	router.HandleFunc(
		"POST /api/v1/auth/login",
		authHandler.Login,
	)

	log.Println("TaskFlow API is running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
