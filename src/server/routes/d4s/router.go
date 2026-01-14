package d4s

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/server/routes/d4s/account"
	"paperlink/util"
)

var log = util.GroupLog("DIGI4SCHOOL")

func InitDigi4SchoolRouter(r *gin.Engine) {
	group := r.Group("/api/v1/d4s")
	group.Use(middleware.Auth)

	account.InitDigi4SchoolAccountRouter(group)

	group.GET("/list", ListBooks)
	group.POST("/takeBook/:id", TakeBook)
}
