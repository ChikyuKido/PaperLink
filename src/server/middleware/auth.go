package middleware

import (
	"net/http"
	"strings"

	"paperlink/db/repo"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "no token found or invalid format")
		return
	}

	token := auth[7:]

	claims, err := util.ParseJWT(token)
	if err != nil || claims == nil {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "token invalid")
		return
	}

	if claims.Type != "access" {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "invalid token type")
		return
	}

	user, err := repo.User.Get(claims.UserID)
	if err != nil || user == nil {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "user not found")
		return
	}

	c.Set("userId", claims.UserID)
	c.Next()
}
