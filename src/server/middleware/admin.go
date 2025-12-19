package middleware

import (
	"net/http"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Admin ensures the authenticated user has admin privileges.
// Requires Auth middleware to have set "userId" in context.
func Admin(c *gin.Context) {
	userID := c.GetInt("userId")
	if userID == 0 {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	user, err := repo.User.Get(userID)
	if err != nil || user == nil {
		c.Abort()
		routes.JSONError(c, http.StatusUnauthorized, "user not found")
		return
	}

	if !user.IsAdmin {
		c.Abort()
		routes.JSONError(c, http.StatusForbidden, "admin permission required")
		return
	}

	c.Next()
}
