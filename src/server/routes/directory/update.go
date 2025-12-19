package directory

import (
	"github.com/gin-gonic/gin"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"strconv"
)

type UpdateDirectoryRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parentId"`
}

func Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, routes.NewError(400, "invalid directory id"))
		return
	}

	userID := c.GetInt("userId")

	dir, err := repo.Directory.Get(id)
	if err != nil {
		c.JSON(404, routes.NewError(404, "directory not found"))
		return
	}
	if dir.UserID != userID {
		c.JSON(403, routes.NewError(403, "you are not authorized to update this directory"))
	}

	var req UpdateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, routes.NewError(400, "failed to parse json: "+err.Error()))
		return
	}

	if req.Name != nil {
		dir.Name = *req.Name
	}
	if req.ParentID != nil {
		dir.ParentID = req.ParentID
	}

	if err := repo.Directory.Save(dir); err != nil {
		c.JSON(409, routes.NewError(409, "failed to update directory: "+err.Error()))
		return
	}

	c.JSON(200, routes.NewSuccess(gin.H{"msg": "ok"}))
}
