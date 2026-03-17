package pdfws

import (
	"paperlink/server/middleware"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("PDFWS")

func InitPDFWSRouter(r *gin.Engine) {
	group := r.Group("/api/v1/pdfws")
	group.GET("/connect/:id", Connect)

	authGroup := group.Group("")
	authGroup.Use(middleware.Auth)
	authGroup.GET("/create/:id", Create)
}
