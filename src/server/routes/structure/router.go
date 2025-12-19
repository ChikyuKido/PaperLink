package structure

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/util"
)

var log = util.GroupLog("STRUCTURE")

func InitStructureRoutes(r *gin.Engine) {
	group := r.Group("/api/v1/structure")
	group.Use(middleware.Auth)
	group.GET("/tree", Tree)
}
