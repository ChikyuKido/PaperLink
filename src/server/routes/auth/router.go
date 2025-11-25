package auth

import "github.com/gin-gonic/gin"

func InitAuthRouter(r *gin.Engine) {
	group := r.Group("/api/v1/auth")
	group.GET("/test", Ping)
}

// @Summary Ping
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.String(200, "pong")
}
