package task

import (
	"errors"
	"net/http"

	"paperlink/server/routes"
	task_service "paperlink/service/task"

	"github.com/gin-gonic/gin"
)

// Stop godoc
// @Summary      Stop running task
// @Description  Stops a running task by ID.
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  routes.SuccessResponse
// @Failure      400  {object}  routes.ErrorResponse "Task is not running"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403  {object}  routes.ErrorResponse "Forbidden"
// @Failure      404  {object}  routes.ErrorResponse "Not found"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/tasks/stop/{id} [post]
// @Security     BearerAuth
func Stop(c *gin.Context) {
	id := c.Param("id")
	err := task_service.StopTask(id)
	if err == nil {
		routes.JSONSuccessOK(c, gin.H{"id": id})
		return
	}

	switch {
	case errors.Is(err, task_service.ErrTaskNotFound):
		routes.JSONError(c, http.StatusNotFound, "task not found")
	case errors.Is(err, task_service.ErrTaskNotRunning):
		routes.JSONError(c, http.StatusBadRequest, "task is not running")
	default:
		log.Errorf("failed to stop task %s: %v", id, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to stop task")
	}
}
