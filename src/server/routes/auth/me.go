package auth

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type MeResponse struct {
	Username string `json:"username"`
}

// Me godoc
// @Summary      Current user
// @Description  Returns basic information about the currently authenticated user.
// @Tags         auth
// @Produce      json
// @Success      200 {object} MeResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Router       /api/v1/auth/me [get]
// @Security     BearerAuth
func Me(c *gin.Context) {
	userID := c.GetInt("userId")
	if userID == 0 {
		routes.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	user, err := repo.User.Get(userID)
	if err != nil || user == nil {
		routes.JSONError(c, http.StatusUnauthorized, "user not found")
		return
	}

	routes.JSONSuccessOK(c, MeResponse{Username: user.Username})
}

