package account

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/service/d4s"
	"strconv"
	"strings"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type SyncD4SResponse struct {
	ID string `json:"id"`
}

// Sync godoc
// @Summary      Sync Digi4School accounts
// @Description  Syncs the books from the accounts to disk
// @Tags         digi4school
// @Produce      json
// @Param        ids  query      string  true  "Comma-separated IDs or 'all'"
// @Success      200  {array}   SyncD4SResponse
// @Failure      400  {object}  routes.ErrorResponse "Invalid IDs"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/digi4school/accounts/sync/{ids} [get]
// @Security     BearerAuth
func Sync(c *gin.Context) {
	idsParam := c.Query("ids")
	var accounts []entity.Digi4SchoolAccount
	var err error

	if idsParam == "all" {
		accounts, err = repo.Digi4SchoolAccount.GetList()
		if err != nil {
			log.Errorf("failed to fetch all digi4school accounts: %v", err)
			routes.JSONError(c, http.StatusInternalServerError, "failed to fetch accounts")
			return
		}
	} else {
		idStrs := strings.Split(idsParam, ",")
		ids := make([]any, 0, len(idStrs))
		for _, s := range idStrs {
			id, convErr := strconv.Atoi(strings.TrimSpace(s))
			if convErr != nil {
				routes.JSONError(c, http.StatusBadRequest, "invalid account IDs")
				return
			}
			ids = append(ids, id)
		}

		accounts, err = repo.Digi4SchoolAccount.GetByIDs(ids)
		if err != nil {
			log.Errorf("failed to fetch digi4school accounts: %v", err)
			routes.JSONError(c, http.StatusInternalServerError, "failed to fetch accounts")
			return
		}
	}

	id, err := d4s.StartSyncTask(accounts)
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to start sync task")
		return
	}

	routes.JSONSuccess(c, 200, SyncD4SResponse{
		ID: id,
	})
}
