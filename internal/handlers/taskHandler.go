package handlers

import (
	"GOBAR/internal/taskService"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type requestBody struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

type TaskHandler struct {
	service taskService.TaskService
}

func NewTaskHandler(s taskService.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func handError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, taskService.ErrEmptyTask):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, taskService.ErrTaskNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal error"})
	}
}

func parseID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *TaskHandler) PostTask(c echo.Context) error {
	var body requestBody

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	task, err := h.service.CreateTask(body.Task, body.IsDone)
	if err != nil {
		return handError(c, err)
	}
	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		return handError(c, err)
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		return handError(c, err)
	}
	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) PatchTask(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var body requestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	task, err := h.service.UpdateTask(id, body.Task, body.IsDone)
	if err != nil {
		return handError(c, err)
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	if err := h.service.DeleteTask(id); err != nil {
		return handError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
