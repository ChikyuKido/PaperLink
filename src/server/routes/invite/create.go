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

// POST /api/v1/invite/create
// nur Admin (wegen Middleware)
func Create(c *gin.Context) {
	invite, err := repo.RegistrationInvite.Create(3) // 3 Tage gültig
	if err != nil {
		log.Errorf("failed to create registration invite: %v", err)
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to create invite"))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(CreateInviteResponse{
		Code:      invite.Code,
		ExpiresAt: invite.ExpiresAt,
	}))
}
