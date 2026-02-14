package task

import (
	"paperlink/server/middleware"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("INVITE")

func InitTasksTasks(r *gin.Engine) {
	group := r.Group("/api/v1/task")
	group.Use(middleware.Auth, middleware.Admin)
	group.GET("/list", List)
	group.GET("/view/:id", View)
	group.POST("/stop/:id", Stop)
}
