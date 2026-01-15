package task

import (
	"net/http"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type ListTasksResponse struct {
	Tasks []entity.Task `json:"tasks"`
}

// List godoc
// @Summary      List tasks
// @Description  Lists all stored tasks.
// @Tags         tasks
// @Produce      json
// @Success      200 {object} ListTasksResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/tasks/lists [get]
// @Security     BearerAuth
func List(c *gin.Context) {
	tasks, err := repo.Task.GetList()
	if err != nil {
		log.Errorf("failed to fetch tasks: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	routes.JSONSuccessOK(c, ListTasksResponse{Tasks: tasks})
}
