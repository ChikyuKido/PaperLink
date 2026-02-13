package document

import (
	"net/http"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Get godoc
// @Summary      Get document
// @Description  Returns a document owned by the authenticated user.
// @Tags         document
// @Produce      json
// @Param        uuid   path      string  true  "Document UUID"
// @Success      200    {object}  entity.Document
// @Failure      400    {object}  routes.ErrorResponse "Invalid document UUID"
// @Failure      401    {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403    {object}  routes.ErrorResponse "Forbidden"
// @Failure      404    {object}  routes.ErrorResponse "Document not found"
// @Failure      500    {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/documents/{uuid} [get]
// @Security     BearerAuth
func Get(c *gin.Context) {
	uuid := c.Param("id")
	if uuid == "" {
		routes.JSONError(c, http.StatusBadRequest, "invalid document uuid")
		return
	}

	userID := c.GetInt("userId")

	doc := repo.Document.GetByUUIDWithTagsAndFile(uuid)
	if doc == nil {
		routes.JSONError(c, http.StatusNotFound, "document not found")
		return
	}

	if doc.UserID != userID {
		routes.JSONError(c, http.StatusForbidden, "not authorized to access this document")
		return
	}

	c.JSON(http.StatusOK, doc)
}
