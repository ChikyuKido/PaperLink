package auth

import (
	"net/http"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh"`
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Issues a new access token using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest  true  "Refresh payload"
// @Success      200      {object}  routes.Response
// @Failure      400      {object}  routes.Response
// @Failure      401      {object}  routes.Response
// @Failure      500      {object}  routes.Response
// @Router       /auth/refresh [post]
func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "missing refresh token"))
		return
	}
	access, err := util.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, routes.NewError(401, "invalid refresh token"))
		return
	}
	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{
		"access": access,
	}))
}
