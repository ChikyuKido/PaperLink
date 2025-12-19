package middleware

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Admin stellt sicher, dass der aktuelle User Admin ist.
// Vorausgesetzt: Auth-Middleware lief vorher und hat "userId" gesetzt.
func Admin(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, routes.NewError(http.StatusUnauthorized, "user not in context"))
		return
	}

	userIDInt, ok := userIDVal.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, routes.NewError(http.StatusUnauthorized, "invalid user id type"))
		return
	}

	user, err := repo.User.Get(userIDInt)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, routes.NewError(http.StatusUnauthorized, "user not found"))
		return
	}

	if !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, routes.NewError(http.StatusForbidden, "admin permission required"))
		return
	}

	c.Next()
}
