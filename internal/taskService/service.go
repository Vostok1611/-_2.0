package taskService

import "errors"

var ErrEmptyTask = errors.New("task cannot be empty")

type TaskService interface {
	CreateTask(task string, is_done bool) (Task, error)
	GetAllTasks() ([]Task, error)
	GetTaskByID(id uint) (Task, error)
	UpdateTask(id uint, task string, is_done bool) (Task, error)
	DeleteTask(id uint) error
}

type taskService struct {
	repo TaskRepository
}

func NewTaskService(r TaskRepository) TaskService {
	return &taskService{repo: r}
}

// CreateTask implements [TaskService].
func (s *taskService) CreateTask(task string, is_done bool) (Task, error) {
	if task == "" {
		return Task{}, ErrEmptyTask
	}
	newTask := Task{
		Task:   task,
		IsDone: is_done,
	}
	if err := s.repo.CreateTask(&newTask); err != nil {
		return Task{}, err
	}
	return newTask, nil
}

// GetAllTasks implements [TaskService].
func (s *taskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

// GetTaskByID implements [TaskService].
func (s *taskService) GetTaskByID(id uint) (Task, error) {
	return s.repo.GetTaskByID(id)
}

// UpdateTask implements [TaskService].
func (s *taskService) UpdateTask(id uint, task string, is_done bool) (Task, error) {
	if task == "" {
		return Task{}, ErrEmptyTask
	}
	existing, err := s.repo.GetTaskByID(id)
	if err != nil {
		return Task{}, err
	}

	existing.Task = task
	existing.IsDone = is_done

	if err := s.repo.UpdateTask(&existing); err != nil {
		return Task{}, err
	}
	return existing, nil
}

// DeleteTask implements [TaskService].
func (s *taskService) DeleteTask(id uint) error {
	if _, err := s.repo.GetTaskByID(id); err != nil {
		return err
	}
	return s.repo.DeleteTask(id)
}
