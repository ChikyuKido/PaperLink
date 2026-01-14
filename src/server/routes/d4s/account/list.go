package account

import (
	"net/http"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type ListD4SAccountsResponse struct {
	Accounts []entity.Digi4SchoolAccount `json:"accounts"`
}

// List godoc
// @Summary      List Digi4School accounts
// @Description  Lists all stored Digi4School accounts (admin only). Passwords are not returned.
// @Tags         digi4school
// @Produce      json
// @Success      200 {object} ListD4SAccountsResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/d4s/account/list [get]
// @Security     BearerAuth
func List(c *gin.Context) {
	accounts, err := repo.Digi4SchoolAccount.GetList()
	if err != nil {
		log.Errorf("failed to fetch digi4school accounts: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to fetch accounts")
		return
	}

	// Never return passwords.
	for i := range accounts {
		accounts[i].Password = ""
	}

	routes.JSONSuccessOK(c, ListD4SAccountsResponse{Accounts: accounts})
}

