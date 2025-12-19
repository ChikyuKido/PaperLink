package invite

import (
	"paperlink/server/middleware"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("INVITE")

func InitInviteRouter(r *gin.Engine) {
	group := r.Group("/api/v1/invite")
	group.Use(middleware.Auth, middleware.Admin)
	group.POST("/create", Create)
}
