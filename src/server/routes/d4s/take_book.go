package d4s

import (
	"net/http"
	"strconv"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type TakeBookResponse struct {
	ID int `json:"id"`
}

// TakeBook godoc
// @Summary      Take Digi4School book
// @Description  Marks a Digi4School book as taken (currently a no-op placeholder) and returns the book ID.
// @Tags         digi4school
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  TakeBookResponse
// @Failure      400  {object}  routes.ErrorResponse "Invalid book ID"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      404  {object}  routes.ErrorResponse "Book not found"
// @Router       /api/v1/d4s/takeBook/{id} [post]
// @Security     BearerAuth
func TakeBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		routes.JSONError(c, http.StatusBadRequest, "invalid book id")
		return
	}

	// Ensure the book exists. (Actual "take" behavior can be implemented later.)
	_, err = repo.Digi4SchoolBook.Get(id)
	if err != nil {
		routes.JSONError(c, http.StatusNotFound, "book not found")
		return
	}

	routes.JSONSuccessOK(c, TakeBookResponse{ID: id})
}

