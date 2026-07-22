package task

import (
	"errors"
	"strings"
	"taskflow/internal/project"
	"time"
)

type ProjectFinder interface {
	FindByID(id int64) (project.Project, error)
}

type Repository interface {
	Create(task Task) (Task, error)
	FindByProjectID(projectID int64) []Task
	FindByID(id int64) (Task, bool)
	Update(task Task) (Task, bool)
	Delete(id int64) bool
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
	projectID int64,
	request CreateRequest,
) (Task, error) {
	_, err := s.projects.FindByID(projectID)
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

	return s.repo.Create(task)
}

func (s *Service) FindByProjectID(projectID int64) ([]Task, error) {
	_, err := s.projects.FindByID(projectID)
	if errors.Is(err, project.ErrNotFound) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.repo.FindByProjectID(projectID), nil
}

func (s *Service) FindByID(id int64) (Task, error) {
	task, found := s.repo.FindByID(id)
	if !found {
		return Task{}, ErrNotFound
	}

	return task, nil
}

func (s *Service) Update(
	id int64,
	request CreateRequest,
) (Task, error) {
	existingTask, err := s.FindByID(id)
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

	updatedTask, updated := s.repo.Update(existingTask)
	if !updated {
		return Task{}, ErrNotFound
	}

	return updatedTask, nil
}

func (s *Service) AdvanceStatus(id int64) (Task, error) {
	existingTask, err := s.FindByID(id)
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

	updatedTask, updated := s.repo.Update(existingTask)
	if !updated {
		return Task{}, ErrNotFound
	}

	return updatedTask, nil
}

func (s *Service) Delete(id int64) error {
	deleted := s.repo.Delete(id)
	if !deleted {
		return ErrNotFound
	}

	return nil
}

func (s *Service) GetProjectSummary(
	projectID int64,
) (ProjectSummary, error) {
	projectTasks, err := s.FindByProjectID(projectID)
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
