package task

import (
	"context"
	"errors"
	"strings"
	"taskflow/internal/project"
	"time"
)

type ProjectFinder interface {
	FindByID(ctx context.Context, id int64) (project.Project, error)
}

type Repository interface {
	Create(ctx context.Context, task Task) (Task, error)
	FindByProjectID(ctx context.Context, projectID int64) ([]Task, error)
	FindByID(ctx context.Context, id int64) (Task, error)
	Update(ctx context.Context, task Task) (Task, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo     Repository
	projects ProjectFinder
}

func NewService(repo Repository, projects ProjectFinder) *Service {
	return &Service{
		repo:     repo,
		projects: projects,
	}
}

var (
	ErrNameRequired     = errors.New("task name is required")
	ErrInvalidStatus    = errors.New("invalid task status")
	ErrNotFound         = errors.New("task not found")
	ErrProjectNotFound  = errors.New("project not found")
	ErrAlreadyCompleted = errors.New("task is already completed")
)

func (s *Service) Create(
	ctx context.Context,
	projectID int64,
	request CreateRequest,
) (Task, error) {
	_, err := s.projects.FindByID(ctx, projectID)
	if errors.Is(err, project.ErrNotFound) {
		return Task{}, ErrProjectNotFound
	}
	if err != nil {
		return Task{}, err
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		return Task{}, ErrNameRequired
	}

	status := request.Status
	if status == "" {
		status = StatusTodo
	}

	if status != StatusTodo &&
		status != StatusInProgress &&
		status != StatusCompleted {
		return Task{}, ErrInvalidStatus
	}

	now := time.Now()

	task := Task{
		ProjectID:   projectID,
		Name:        name,
		Description: request.Description,
		Status:      status,
		Priority:    request.Priority,
		DueDate:     request.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, task)
}

func (s *Service) FindByProjectID(ctx context.Context, projectID int64) ([]Task, error) {
	_, err := s.projects.FindByID(ctx, projectID)
	if errors.Is(err, project.ErrNotFound) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.repo.FindByProjectID(ctx, projectID)
}

func (s *Service) FindByID(ctx context.Context, id int64) (Task, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(
	ctx context.Context,
	id int64,
	request CreateRequest,
) (Task, error) {
	existingTask, err := s.FindByID(ctx, id)
	if err != nil {
		return Task{}, err
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		return Task{}, ErrNameRequired
	}

	if request.Status != StatusTodo &&
		request.Status != StatusInProgress &&
		request.Status != StatusCompleted {
		return Task{}, ErrInvalidStatus
	}

	existingTask.Name = name
	existingTask.Description = request.Description
	existingTask.Status = request.Status
	existingTask.Priority = request.Priority
	existingTask.DueDate = request.DueDate
	existingTask.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existingTask)

}

func (s *Service) AdvanceStatus(ctx context.Context, id int64) (Task, error) {
	existingTask, err := s.FindByID(ctx, id)
	if err != nil {
		return Task{}, err
	}

	switch existingTask.Status {
	case StatusTodo:
		existingTask.Status = StatusInProgress

	case StatusInProgress:
		existingTask.Status = StatusCompleted

	case StatusCompleted:
		return Task{}, ErrAlreadyCompleted

	default:
		return Task{}, ErrInvalidStatus
	}

	existingTask.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existingTask)

}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetProjectSummary(
	ctx context.Context,
	projectID int64,
) (ProjectSummary, error) {
	projectTasks, err := s.FindByProjectID(ctx, projectID)
	if err != nil {
		return ProjectSummary{}, err
	}

	summary := ProjectSummary{
		ProjectID:  projectID,
		TotalTasks: len(projectTasks),
	}

	for _, task := range projectTasks {
		switch task.Status {
		case StatusTodo:
			summary.Todo++

		case StatusInProgress:
			summary.InProgress++

		case StatusCompleted:
			summary.Completed++
		}
	}

	if summary.TotalTasks > 0 {
		summary.CompletionPercentage =
			float64(summary.Completed) /
				float64(summary.TotalTasks) *
				100
	}

	return summary, nil
}
