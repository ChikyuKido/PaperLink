package middleware

import (
	"net/http"
	"paperlink/server/routes"
	"paperlink/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, routes.NewError(http.StatusUnauthorized, "no token found or invalid format"))
		return
	}
	token = token[7:]
	claims, err := util.ParseJWT(token)
	if err != nil || claims == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, routes.NewError(http.StatusUnauthorized, "token invalid"))
		return
	}

	c.Set("userId", claims.UserID)
}
