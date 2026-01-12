package admin

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/util"
)

var log = util.GroupLog("ADMIN")

func InitAdminRouter(r *gin.Engine) {
	group := r.Group("/api/v1/admin")
	group.Use(middleware.Auth, middleware.Admin)

	group.GET("/stats", Stats)
}

