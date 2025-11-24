package auth

import "github.com/gin-gonic/gin"

func InitAuthRouter(r *gin.Engine) {
	group := r.Group("/api/v1/auth")
	group.GET("/test", func(context *gin.Context) {
		context.JSON(200, gin.H{"test": "test"})
	})
}
