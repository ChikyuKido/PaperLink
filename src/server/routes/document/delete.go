package document

import (
	"paperlink/db/repo"
	"paperlink/server/routes"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, routes.NewError(400, "invalid document id"))
		return
	}

	userID := c.GetInt("userId")

	doc, err := repo.Document.Get(id)
	if err != nil {
		c.JSON(404, routes.NewError(404, "document not found"))
		return
	}
	if doc.UserID != userID {
		c.JSON(403, routes.NewError(403, "you are not authorized to delete this document"))
		return
	}

	if err := repo.Document.Delete(id); err != nil {
		c.JSON(404, routes.NewError(404, "failed to delete document: "+err.Error()))
		return
	}

	c.JSON(200, routes.NewSuccess(gin.H{"msg": "ok"}))
}
