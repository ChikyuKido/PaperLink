package digi4school

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
	"paperlink/server/routes/digi4school/account"
	"paperlink/util"
)

var log = util.GroupLog("DIGI4SCHOOL")

func InitDigi4SchoolRouter(r *gin.Engine) {
	group := r.Group("/api/v1/d4s")
	group.Use(middleware.Auth)
	account.InitDigi4SchoolAccountRouter(group)
}
