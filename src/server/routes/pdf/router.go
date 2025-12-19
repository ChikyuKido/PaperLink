package pdf

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
)

func InitPDFRouter(r *gin.Engine) {
	group := r.Group("/api/v1/pdf")
	group.Use(middleware.Auth)
	group.GET("/:id/:page", GetPage)
}
