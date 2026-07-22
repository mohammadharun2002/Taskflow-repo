package project

import "sync"

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

func (r *MemoryRepository) Create(project Project) (Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	project.ID = r.nextID
	r.nextID++

	r.projects = append(r.projects, project)

	return project, nil
}

func (r *MemoryRepository) FindAll() []Project {
	r.mu.RLock()
	defer r.mu.RUnlock()
	projects := make([]Project, 0, len(r.projects))
	for _, project := range r.projects {
		projects = append(projects, project)
	}
	return projects
}

func (r *MemoryRepository) FindByID(id int64) (Project, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, project := range r.projects {
		if project.ID == id {
			return project, true
		}
	}
	return Project{}, false
}

func (r *MemoryRepository) Delete(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, p := range r.projects {
		if p.ID == id {
			r.projects = append(r.projects[:index], r.projects[index+1:]...)
			return true
		}
	}
	return false
}
