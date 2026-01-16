package task

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
	task_service "paperlink/service/task"
)

type GetTaskResponse struct {
	ID        string            `gorm:"primary_key" json:"id"`
	Status    entity.TaskStatus `json:"status"`
	Name      string            `gorm:"not null" json:"name"`
	StartTime int64             `json:"startTime"`
	EndTime   int64             `json:"endTime"`
	Content   []string          `json:"content"`
}

// View godoc
// @Summary      Get task
// @Description  Returns a single task by ID.
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  GetTaskResponse
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403  {object}  routes.ErrorResponse "Forbidden"
// @Failure      404  {object}  routes.ErrorResponse "Not found"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/tasks/{id} [get]
// @Security     BearerAuth
func View(c *gin.Context) {
	id := c.Param("id")
	task, err := repo.Task.Get(id)
	if err != nil {
		log.Errorf("failed to fetch task %s: %v", id, err)
		routes.JSONError(c, http.StatusNotFound, "task not found")
		return
	}

	t, err := task_service.GetTaskInfo(id)
	if err != nil {
		log.Errorf("failed to fetch content task %s: %v", id, err)
		routes.JSONError(c, http.StatusNotFound, "task content not found")
		return
	}
	routes.JSONSuccessOK(c, GetTaskResponse{
		ID:        task.ID,
		Status:    task.Status,
		Name:      task.Name,
		StartTime: task.StartTime,
		EndTime:   task.EndTime,
		Content:   t.Lines,
	})
}
