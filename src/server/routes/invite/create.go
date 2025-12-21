package invite

import (
	"net/http"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type CreateInviteResponse struct {
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expiresAt"`
}

// Create godoc
// @Summary      Create registration invite
// @Description  Creates a new registration invite (admin only).
// @Tags         invite
// @Produce      json
// @Success      200 {object} CreateInviteResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/invite/create [post]
// @Security     BearerAuth
func Create(c *gin.Context) {
	invite, err := repo.RegistrationInvite.Create(3, 1)
	if err != nil {
		log.Errorf("failed to create registration invite: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to create invite")
		return
	}

	routes.JSONSuccess(c, http.StatusOK, CreateInviteResponse{
		Code:      invite.Code,
		ExpiresAt: invite.ExpiresAt,
	})
}
