package auth

import (
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.Engine) {
	group := r.Group("/api/v1/auth")

	group.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"test": "test"})
	})

	// ✅ Add register + login routes here
	group.POST("/register", Register)
	group.POST("/login", Login)
}
