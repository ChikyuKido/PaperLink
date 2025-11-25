package pdf

import "github.com/gin-gonic/gin"

func InitPDFRouter(r *gin.Engine) {
	group := r.Group("/api/v1/pdf")
	group.GET("/:id/:page", GetPage())
}
