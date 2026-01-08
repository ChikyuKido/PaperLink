package account

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"paperlink/service/d4s"

	"github.com/gin-gonic/gin"
)

type CreateDigi4SchoolAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateDigi4SchoolAccountResponse struct {
	ID int `json:"id"`
}

// Create godoc
// @Summary      Create Digi4School account
// @Description  Creates a new Digi4School account.
// @Tags         digi4school
// @Accept       json
// @Produce      json
// @Param        request body CreateDigi4SchoolAccountRequest true "Create Digi4School account payload"
// @Success      201 {object} CreateDigi4SchoolAccountResponse
// @Failure      400 {object} routes.ErrorResponse "Invalid request body"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      409 {object} routes.ErrorResponse "Failed to create account"
// @Router       /api/v1/digi4school/accounts [post]
// @Security     BearerAuth
func Create(c *gin.Context) {
	var req CreateDigi4SchoolAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid create digi4school account body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	account := entity.Digi4SchoolAccount{
		Username: req.Username,
		Password: req.Password,
	}

	if !d4s.TestLogin(&account) {
		log.Warnf("invalid credentials for account")
		routes.JSONError(c, http.StatusBadRequest, "invalid credentials for account")
		return
	}

	if err := repo.Digi4SchoolAccount.Save(&account); err != nil {
		log.Errorf("failed to save digi4school account: %v", err)
		routes.JSONError(c, http.StatusConflict, "failed to create account")
		return
	}

	routes.JSONSuccess(c, http.StatusCreated, CreateDigi4SchoolAccountResponse{
		ID: account.ID,
	})
}
