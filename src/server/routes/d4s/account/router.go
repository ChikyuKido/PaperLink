package account

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/util"
)

var log = util.GroupLog("DIGI4SCHOOL_ACCOUNT")

func InitDigi4SchoolAccountRouter(r *gin.RouterGroup) {
	group := r.Group("/account")
	group.Use(middleware.Admin)
	group.POST("/create", Create)
	group.DELETE("/delete/:id", Delete)
	group.GET("/sync", Sync)

}
