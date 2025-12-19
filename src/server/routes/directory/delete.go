package directory

import (
	"github.com/gin-gonic/gin"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"strconv"
)

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, routes.NewError(400, "invalid directory id"))
		return
	}

	if err := repo.Directory.Delete(id); err != nil {
		c.JSON(404, routes.NewError(404, "failed to delete directory: "+err.Error()))
		return
	}

	c.JSON(200, routes.NewSuccess(gin.H{"msg": "ok"}))
}
