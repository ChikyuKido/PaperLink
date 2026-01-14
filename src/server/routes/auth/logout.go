package auth

import (
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type LogoutResponse struct {
	Ok bool `json:"ok"`
}

// Logout godoc
// @Summary      Logout
// @Description  Clears the refresh token cookie.
// @Tags         auth
// @Produce      json
// @Success      200 {object} LogoutResponse
// @Router       /api/v1/auth/logout [post]
func Logout(c *gin.Context) {
	// Browsers delete cookies by matching name + path (+ domain).
	// We currently set Path=/ at login, but clear a few legacy paths too.
	paths := []string{
		"/",
		"/api/v1/auth/refresh",
		"/api/v1/auth",
		"/api/v1",
	}

	for _, p := range paths {
		// Try both maxAge=0 and maxAge=-1 for broader compatibility
		c.SetCookie("refresh", "", 0, p, "", false, true)
		c.SetCookie("refresh", "", -1, p, "", false, true)
	}

	routes.JSONSuccessOK(c, LogoutResponse{Ok: true})
}
