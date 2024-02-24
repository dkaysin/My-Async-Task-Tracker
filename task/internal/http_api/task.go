package http_handler

import (
	"async_course/task"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// get tasks

type getTasksReq struct {
	UserID string `json:"userID"`
}

type getTasksRes struct {
	Tasks []task.Task `json:"tasks"`
}

func (h *HttpAPI) getTasks(c echo.Context) error {
	payload, err := validatePayload[getTasksReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	tasks, err := h.s.GetTasks(context.Background(), payload.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(getTasksRes{tasks}))
}

// create task

type createTaskReq struct {
	Description string `json:"description"`
}

type createTaskRes struct {
	TaskID string `json:"task_id"`
}

func (h *HttpAPI) createTask(c echo.Context) error {
	payload, err := validatePayload[createTaskReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	taskID, err := h.s.CreateTask(context.Background(), payload.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(createTaskRes{taskID}))
}

// complete task

type completeTaskReq struct {
	TaskID string `json:"task_id"`
	UserID string `json:"user_id"`
}

func (h *HttpAPI) completeTask(c echo.Context) error {
	payload, err := validatePayload[completeTaskReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	err = h.s.CompleteTask(context.Background(), payload.TaskID, payload.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(nil))
}

// assign tasks

func (h *HttpAPI) assignTasks(c echo.Context) error {
	err := h.s.AssignTasks(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
