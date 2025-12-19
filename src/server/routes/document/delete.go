package document

import (
	"net/http"
	"strconv"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary      Delete document
// @Description  Deletes a document owned by the authenticated user.
// @Tags         document
// @Produce      json
// @Param        id   path      int  true  "Document ID"
// @Success      204  "No Content"
// @Failure      400  {object}  routes.ErrorResponse "Invalid document ID"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403  {object}  routes.ErrorResponse "Forbidden"
// @Failure      404  {object}  routes.ErrorResponse "Document not found"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/documents/delete/{id} [delete]
// @Security     BearerAuth
func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		routes.JSONError(c, http.StatusBadRequest, "invalid document id")
		return
	}

	userID := c.GetInt("userId")

	doc, err := repo.Document.Get(id)
	if err != nil || doc == nil {
		routes.JSONError(c, http.StatusNotFound, "document not found")
		return
	}

	if doc.UserID != userID {
		routes.JSONError(c, http.StatusForbidden, "not authorized to delete this document")
		return
	}

	if err := repo.Document.Delete(id); err != nil {
		log.Errorf("failed to delete document %d: %v", id, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to delete document")
		return
	}

	c.Status(http.StatusNoContent)
}
