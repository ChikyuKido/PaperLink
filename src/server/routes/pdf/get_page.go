package pdf

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"paperlink/db/repo"

	"github.com/gin-gonic/gin"
)

// GetPage godoc
// @Summary      Fetch PDF file
// @Description  Returns the stored PDF file and supports byte range requests.
// @Tags         pdf
// @Param        id    path string true "Document ID"
// @Produce      application/pdf
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "document not found"
// @Failure      500 {string} string "failed to open pdf"
// @Router       /pdf/{id} [get]
// @Security     BearerAuth
func GetPage(c *gin.Context) {
	docUUID := c.Param("id")
	userID := c.GetInt("userId")

	doc := repo.Document.GetByUUIDWithFile(docUUID)
	if doc == nil {
		c.String(http.StatusNotFound, "document not found")
		return
	}

	if doc.UserID != userID {
		c.String(http.StatusForbidden, "forbidden")
		return
	}

	file := doc.File

	f, err := os.Open(file.Path)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to open pdf")
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to open pdf")
		return
	}

	filename := filepath.Base(file.Path)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	c.Header("Accept-Ranges", "bytes")
	http.ServeContent(c.Writer, c.Request, filename, stat.ModTime(), f)
}
