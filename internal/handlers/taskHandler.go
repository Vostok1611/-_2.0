package handlers

import (
	"context"
	"errors"

	"GOBAR/internal/taskService"
	"GOBAR/internal/web/tasks"
)

type TaskHandler struct {
	service taskService.TaskService
}

func NewTaskHandler(s taskService.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// из модели БД в модель API
func toAPI(t taskService.Task) tasks.Task {
	return tasks.Task{
		Id:     t.ID,
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

func (h *TaskHandler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	all, err := h.service.GetAllTasks()
	if err != nil {
		return tasks.GetTasks500JSONResponse{Error: "Internal error"}, nil
	}

	response := tasks.GetTasks200JSONResponse{}
	for _, t := range all {
		response = append(response, toAPI(t))
	}
	return response, nil
}

func (h *TaskHandler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	isDone := false
	if request.Body.IsDone != nil {
		isDone = *request.Body.IsDone
	}

	created, err := h.service.CreateTask(request.Body.Task, isDone)
	if err != nil {
		if errors.Is(err, taskService.ErrEmptyTask) {
			return tasks.PostTasks400JSONResponse{Error: err.Error()}, nil
		}
		return tasks.PostTasks500JSONResponse{Error: "Internal error"}, nil
	}

	return tasks.PostTasks201JSONResponse(toAPI(created)), nil
}

func (h *TaskHandler) GetTasksById(_ context.Context, request tasks.GetTasksByIdRequestObject) (tasks.GetTasksByIdResponseObject, error) {
	task, err := h.service.GetTaskByID(request.Id)
	if err != nil {
		if errors.Is(err, taskService.ErrTaskNotFound) {
			return tasks.GetTasksById404JSONResponse{Error: err.Error()}, nil
		}
		return tasks.GetTasksById500JSONResponse{Error: "Internal error"}, nil
	}

	return tasks.GetTasksById200JSONResponse(toAPI(task)), nil
}

func (h *TaskHandler) PatchTasksById(_ context.Context, request tasks.PatchTasksByIdRequestObject) (tasks.PatchTasksByIdResponseObject, error) {
	isDone := false
	if request.Body.IsDone != nil {
		isDone = *request.Body.IsDone
	}

	updated, err := h.service.UpdateTask(request.Id, request.Body.Task, isDone)
	if err != nil {
		switch {
		case errors.Is(err, taskService.ErrEmptyTask):
			return tasks.PatchTasksById400JSONResponse{Error: err.Error()}, nil
		case errors.Is(err, taskService.ErrTaskNotFound):
			return tasks.PatchTasksById404JSONResponse{Error: err.Error()}, nil
		default:
			return tasks.PatchTasksById500JSONResponse{Error: "Internal error"}, nil
		}
	}

	return tasks.PatchTasksById200JSONResponse(toAPI(updated)), nil
}

func (h *TaskHandler) DeleteTasksById(_ context.Context, request tasks.DeleteTasksByIdRequestObject) (tasks.DeleteTasksByIdResponseObject, error) {
	if err := h.service.DeleteTask(request.Id); err != nil {
		if errors.Is(err, taskService.ErrTaskNotFound) {
			return tasks.DeleteTasksById404JSONResponse{Error: err.Error()}, nil
		}
		return tasks.DeleteTasksById500JSONResponse{Error: "Internal error"}, nil
	}

	return tasks.DeleteTasksById204Response{}, nil
}
