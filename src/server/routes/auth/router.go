package auth

import (
	"github.com/gin-gonic/gin"
)

func InitAuthRouter(r *gin.Engine) {
	group := r.Group("/api/v1/auth")
	group.POST("/register", Register)
	group.POST("/login", Login)
}
