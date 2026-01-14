package d4s

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type ListBooksResponse struct {
	Books []entity.Digi4SchoolBook `json:"books"`
}

// ListBooks godoc
// @Summary      List Digi4School books
// @Description  Lists all synced Digi4School books.
// @Tags         digi4school
// @Produce      json
// @Success      200 {object} ListBooksResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/d4s/list [get]
// @Security     BearerAuth
func ListBooks(c *gin.Context) {
	books, err := repo.Digi4SchoolBook.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list books")
		return
	}
	routes.JSONSuccessOK(c, ListBooksResponse{Books: books})
}
