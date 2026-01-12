package invite

import (
	"net/http"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type CreateInviteRequest struct {
	// ValidDays is how many days the invite should remain valid. Defaults to 3.
	ValidDays int `json:"validDays"`
	// Uses is how many times the invite can be used. Defaults to 1.
	Uses int `json:"uses"`
}

type CreateInviteResponse struct {
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expiresAt"`
	Uses      int    `json:"uses"`
}

// Create godoc
// @Summary      Create registration invite
// @Description  Creates a new registration invite (admin only).
// @Tags         invite
// @Accept       json
// @Produce      json
// @Param        request  body      CreateInviteRequest  false  "Invite creation options"
// @Success      200 {object} CreateInviteResponse
// @Failure      400 {object} routes.ErrorResponse "Invalid request body"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/invite/create [post]
// @Security     BearerAuth
func Create(c *gin.Context) {
	// Optional body
	var req CreateInviteRequest
	_ = c.ShouldBindJSON(&req)

	validDays := req.ValidDays
	uses := req.Uses
	if uses <= 0 {
		uses = 1
	}

	invite, err := repo.RegistrationInvite.Create(validDays, uses)
	if err != nil {
		log.Errorf("failed to create registration invite: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to create invite")
		return
	}

	routes.JSONSuccess(c, http.StatusOK, CreateInviteResponse{
		Code:      invite.Code,
		ExpiresAt: invite.ExpiresAt,
		Uses:      invite.Uses,
	})
}
