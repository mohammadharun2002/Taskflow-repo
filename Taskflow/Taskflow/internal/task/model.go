package task

import "time"

type TaskFlow string

const (
	StatusTodo       TaskFlow = "todo"
	StatusInProgress TaskFlow = "in_progress"
	StatusCompleted  TaskFlow = "completed"
)

type Task struct {
	ID          int64      `json:"id"`
	ProjectID   int64      `json:"project_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      TaskFlow   `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      TaskFlow   `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

type ProjectSummary struct {
	ProjectID            int64   `json:"project_id"`
	TotalTasks           int     `json:"total_tasks"`
	Todo                 int     `json:"todo"`
	InProgress           int     `json:"in_progress"`
	Completed            int     `json:"completed"`
	CompletionPercentage float64 `json:"completion_percentage"`
}
