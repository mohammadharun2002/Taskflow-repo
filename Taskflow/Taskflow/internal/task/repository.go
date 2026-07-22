package task

import "sync"

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

func (r *MemoryRepository) Create(task Task) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *MemoryRepository) FindByProjectID(projectID int64) []Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projectTasks := make([]Task, 0)
	for _, task := range r.tasks {
		if task.ProjectID == projectID {
			projectTasks = append(projectTasks, task)
		}
	}

	return projectTasks
}

func (r *MemoryRepository) FindByID(id int64) (Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.ID == id {
			return task, true
		}
	}

	return Task{}, false
}

func (r *MemoryRepository) Update(updatedTask Task) (Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, task := range r.tasks {
		if task.ID == updatedTask.ID {
			r.tasks[index] = updatedTask
			return updatedTask, true
		}
	}

	return Task{}, false
}

func (r *MemoryRepository) Delete(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:index], r.tasks[index+1:]...)
			return true
		}
	}

	return false
}
