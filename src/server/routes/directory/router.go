package directory

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/util"
)

var log = util.GroupLog("DIRECTORY")

func InitDirectoryRouter(r *gin.Engine) {
	group := r.Group("/api/v1/directory")
	group.Use(middleware.Auth)
	group.POST("/create", Create)
	group.DELETE("/delete/:id", Delete)
	group.PATCH("/update/:id", Update)
}
