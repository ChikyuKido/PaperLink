package document

import (
	"paperlink/server/middleware"

	"github.com/gin-gonic/gin"
)

func InitDocumentRouter(r *gin.Engine) {
	group := r.Group("/api/v1/document")
	group.Use(middleware.Auth)
	group.POST("/create", Create)
	group.POST("/upload", Upload)
}
