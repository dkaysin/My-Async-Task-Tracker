package http_handler

import (
	"async_course/task"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) RegisterPublic(g *echo.Group) {
	g.GET("/status", h.status)
}

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.GET("/tasks/:user_id", h.requireRoles(h.getTasks, []string{task.RoleAdmin, task.RoleManager}))
	g.GET("/tasks/", h.requireRoles(h.getTasks, []string{task.RoleDeveloper, task.RoleAdmin, task.RoleManager}))
	g.POST("/create-task", h.requireRoles(h.createTask, []string{task.RoleDeveloper, task.RoleAdmin, task.RoleManager, task.RoleAccountant}))
	g.POST("/complete-task", h.requireRoles(h.completeTask, []string{task.RoleDeveloper, task.RoleAdmin, task.RoleManager}))
	g.POST("/assign-tasks", h.requireRoles(h.assignTasks, []string{task.RoleAdmin, task.RoleManager}))
}

func (h *HttpAPI) status(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
