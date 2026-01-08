package account

import (
	"net/http"
	"strconv"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary      Delete Digi4School account
// @Description  Deletes a Digi4School account.
// @Tags         digi4school
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      204  "No Content"
// @Failure      400  {object}  routes.ErrorResponse "Invalid account ID"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      404  {object}  routes.ErrorResponse "Account not found"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/digi4school/accounts/{id} [delete]
// @Security     BearerAuth
func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		routes.JSONError(c, http.StatusBadRequest, "invalid account id")
		return
	}

	if err := repo.Digi4SchoolAccount.Delete(id); err != nil {
		log.Errorf("failed to delete digi4school account %d: %v", id, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to delete account")
		return
	}

	c.Status(http.StatusNoContent)
}
