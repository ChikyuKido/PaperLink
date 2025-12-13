package auth

import (
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("AUTH")

func InitAuthRouter(r *gin.Engine) {
	group := r.Group("/api/v1/auth")
	group.POST("/register", Register)
	group.POST("/login", Login)
	group.POST("/refresh", Refresh)
}
