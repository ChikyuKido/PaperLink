package auth

import "github.com/gin-gonic/gin"

// HasAdmin godoc
// @Summary      Check admin access
// @Description  Endpoint to verify if the authenticated user has admin privileges.
//
//	Only accessible if the auth and admin middleware succeed.
//
// @Tags         admin
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Access granted"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403  {object}  routes.ErrorResponse "Forbidden"
// @Router       /api/v1/admin/has [get]
// @Security     BearerAuth
func HasAdmin(c *gin.Context) {
	c.JSON(200, gin.H{})
}
