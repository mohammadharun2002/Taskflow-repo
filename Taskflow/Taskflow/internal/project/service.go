package project

import (
	"errors"
	"strings"
	"time"
)

type Repository interface {
	Create(project Project) (Project, error)
	FindAll() []Project
	FindByID(id int64) (Project, bool)
	Delete(id int64) bool
}

type Service struct {
	repo Repository
}

var ErrNameRequired = errors.New("project name is required")
var ErrNotFound = errors.New("project not found")

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(request CreateRequest) (Project, error) {
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

	return s.repo.Create(project)
}

func (s *Service) FindAll() []Project {
	return s.repo.FindAll()
}

func (s *Service) FindByID(id int64) (Project, error) {
	project, found := s.repo.FindByID(id)
	if !found {
		return Project{}, ErrNotFound
	}
	return project, nil
}

func (s *Service) Delete(id int64) error {
	deleted := s.repo.Delete(id)
	if !deleted {
		return ErrNotFound
	}
	return nil
}
