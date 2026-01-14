package d4s

import (
	"net/http"
	"strconv"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	userID := c.GetInt("userId")

	book, err := repo.Digi4SchoolBook.Get(id)
	if err != nil || book == nil {
		routes.JSONError(c, http.StatusNotFound, "book not found")
		return
	}

	if book.FileUUID == "" {
		routes.JSONError(c, http.StatusBadRequest, "book has no file")
		return
	}

	// Create a normal Document that points to the already-downloaded FileDocument.
	doc := entity.Document{
		UUID:        uuid.NewString(),
		Name:        book.BookName,
		Description: "Digi4School book",
		UserID:      userID,
		FileUUID:    book.FileUUID,
	}

	if err := repo.Document.Save(&doc); err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to create document")
		return
	}

	// Also add a DocumentUser link (future-proof for permissions).
	_ = repo.DocumentUser.Save(&entity.DocumentUser{
		UserID:     userID,
		DocumentID: doc.ID,
		Role:       entity.Editor,
	})

	routes.JSONSuccessOK(c, TakeBookResponse{ID: doc.ID})
}
