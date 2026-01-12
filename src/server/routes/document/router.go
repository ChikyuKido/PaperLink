package document

import (
	"paperlink/server/middleware"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("DOCUMENT")

func InitDocumentRouter(r *gin.Engine) {
	group := r.Group("/api/v1/document")
	group.Use(middleware.Auth)
	group.GET("/filter", Filter)
	group.POST("/update", Update)
	group.POST("/create", Create)
	group.POST("/upload", Upload)
	group.DELETE("/delete/:id", Delete)
}
