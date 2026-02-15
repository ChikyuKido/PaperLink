package pdf

import (
	"errors"
	"net/http"
	"os"
	"paperlink/ptf"
	"path/filepath"
	"strconv"
	"strings"

	"paperlink/db/repo"

	"github.com/gin-gonic/gin"
)

func getThumbPath(docFilePath string) string {
	return strings.TrimSuffix(docFilePath, filepath.Ext(docFilePath)) + "_thumb.ptf"
}

func parseThumbnailRange(raw string) (int, int, error) {
	parts := strings.SplitN(raw, "-", 2)
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid range")
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil || start < 0 {
		return 0, 0, errors.New("invalid start")
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil || end < start {
		return 0, 0, errors.New("invalid end")
	}
	return start, end, nil
}

// GetThumbnailsRange godoc
// @Summary      Fetch document thumbnail range
// @Description  Returns a PTF chunk for the requested zero-based inclusive range (e.g. 0-50).
// @Tags         pdf
// @Param        id     path string true "Document ID"
// @Param        range  path string true "Thumbnail index range (start-end)"
// @Produce      application/octet-stream
// @Failure      400 {string} string "invalid range"
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "document not found"
// @Failure      500 {string} string "failed to read thumbnails"
// @Router       /pdf/thumbnails/{id}/{range} [get]
// @Security     BearerAuth
func GetThumbnailsRange(c *gin.Context) {
	docUUID := c.Param("id")
	rangeParam := c.Param("range")
	userID := c.GetInt("userId")

	start, end, err := parseThumbnailRange(rangeParam)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid range")
		return
	}

	doc := repo.Document.GetByUUIDWithFile(docUUID)
	if doc == nil {
		c.String(http.StatusNotFound, "document not found")
		return
	}

	if doc.UserID != userID {
		c.String(http.StatusForbidden, "forbidden")
		return
	}

	thumbPath := getThumbPath(doc.File.Path)
	data, err := ptf.Read(thumbPath, ptf.ReadOptions{
		HasRange: true,
		Start:    uint64(start),
		End:      uint64(end),
	})
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "thumbnails not found")
			return
		}
		c.String(http.StatusInternalServerError, "failed to read thumbnails")
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}
