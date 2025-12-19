package directory

import (
	"github.com/gin-gonic/gin"
	"paperlink/server/middleware"
)

func InitDirectoryRouter(r *gin.Engine) {
	group := r.Group("/api/v1/directory")
	group.Use(middleware.Auth)
	group.POST("/create", Create)
	group.POST("/delete", Delete)
	group.POST("/update", Update)
}
