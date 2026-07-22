package main

import (
	"log"
	"net/http"

	"taskflow/internal/health"
	"taskflow/internal/project"
	"taskflow/internal/task"
)

func main() {
	projectRepository := project.NewMemoryRepository()
	projectService := project.NewService(projectRepository)
	projectHandler := project.NewHandler(projectService)

	taskRepository := task.NewMemoryRepository()
	taskService := task.NewService(taskRepository, projectService)
	taskHandler := task.NewHandler(taskService)

	router := http.NewServeMux()

	router.HandleFunc("GET /healthz", health.HealthCheck)
	router.HandleFunc("GET /readyz", health.ReadinessCheck)

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

	log.Println("TaskFlow API is running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
