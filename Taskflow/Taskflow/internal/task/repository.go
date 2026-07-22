package task

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	tasks  []Task
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (r *MemoryRepository) Create(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *MemoryRepository) FindByProjectID(ctx context.Context, projectID int64) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	projectTasks := make([]Task, 0)
	for _, task := range r.tasks {
		if task.ProjectID == projectID {
			projectTasks = append(projectTasks, task)
		}
	}

	return projectTasks, nil
}

func (r *MemoryRepository) FindByID(ctx context.Context, id int64) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, ErrNotFound
}

func (r *MemoryRepository) Update(ctx context.Context, updatedTask Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for index, task := range r.tasks {
		if task.ID == updatedTask.ID {
			r.tasks[index] = updatedTask
			return updatedTask, nil
		}
	}

	return Task{}, ErrNotFound
}

func (r *MemoryRepository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for index, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:index], r.tasks[index+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
