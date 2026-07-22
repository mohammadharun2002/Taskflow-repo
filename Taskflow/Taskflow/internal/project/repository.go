package project

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	projects []Project
	nextID   int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		projects: make([]Project, 0),
		nextID:   1,
	}
}

func (r *MemoryRepository) Create(ctx context.Context, project Project) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	project.ID = r.nextID
	r.nextID++

	r.projects = append(r.projects, project)

	return project, nil
}

func (r *MemoryRepository) FindAll(ctx context.Context) ([]Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	projects := make([]Project, 0, len(r.projects))
	for _, project := range r.projects {
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *MemoryRepository) FindByID(ctx context.Context, id int64) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, project := range r.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return Project{}, ErrNotFound
}

func (r *MemoryRepository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for index, p := range r.projects {
		if p.ID == id {
			r.projects = append(r.projects[:index], r.projects[index+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
