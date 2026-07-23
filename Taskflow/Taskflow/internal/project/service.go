package project

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Repository interface {
	Create(ctx context.Context, project Project) (Project, error)
	FindAll(ctx context.Context) ([]Project, error)
	FindByID(ctx context.Context, id int64) (Project, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo Repository
}

var ErrNameRequired = errors.New("project name is required")
var ErrNotFound = errors.New("project not found")
var ErrHasTasks = errors.New("project has existing tasks")

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Project, error) {
	name := strings.TrimSpace(request.Name)

	if name == "" {
		return Project{}, ErrNameRequired
	}

	now := time.Now()

	project := Project{
		Name:        name,
		Description: request.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, project)
}

func (s *Service) FindAll(ctx context.Context) ([]Project, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id int64) (Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
