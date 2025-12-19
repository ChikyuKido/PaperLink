package directory

import (
	"github.com/gin-gonic/gin"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"
)

type CreateDirectoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int   `json:"parentId"`
}

func Create(c *gin.Context) {
	var req CreateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, routes.NewError(400, "failed to parse json: "+err.Error()))
		return
	}

	userID := c.GetInt("userId")

	dir := entity.Directory{
		UserID:   userID,
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if err := repo.Directory.Save(&dir); err != nil {
		c.JSON(409, routes.NewError(409, "failed to save directory: "+err.Error()))
		return
	}

	c.JSON(201, routes.NewSuccess(gin.H{"id": dir.ID}))
}
