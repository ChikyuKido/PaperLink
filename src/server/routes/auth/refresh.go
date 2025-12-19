package auth

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

// Refresh godoc
// @Summary      Refresh access token
// @Description  Issues a new access token using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200      {object}  LoginResponse
// @Failure      401      {object}  routes.ErrorResponse "Invalid or missing refresh token"
// @Failure      500      {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/auth/refresh [post]
func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh")
	if err != nil || refreshToken == "" {
		routes.JSONError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	claims, err := util.ParseJWT(refreshToken)
	if err != nil {
		routes.JSONError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	user, err := repo.User.Get(claims.UserID)
	if err != nil || user == nil {
		routes.JSONError(c, http.StatusUnauthorized, "user no longer exists")
		return
	}

	access, err := util.RefreshAccessToken(refreshToken)
	if err != nil {
		routes.JSONError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	routes.JSONSuccess(c, http.StatusOK, LoginResponse{
		AccessToken: access,
	})

}
