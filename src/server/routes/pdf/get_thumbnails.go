package pdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"os"
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

func readThumbnailRange(path string, start, end int) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out bytes.Buffer
	index := 0
	selected := 0

	for {
		var size uint64
		if err := binary.Read(f, binary.LittleEndian, &size); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if index >= start && index <= end {
			if err := binary.Write(&out, binary.LittleEndian, size); err != nil {
				return nil, err
			}
			if _, err := io.CopyN(&out, f, int64(size)); err != nil {
				return nil, err
			}
			selected++
		} else {
			if _, err := f.Seek(int64(size), io.SeekCurrent); err != nil {
				return nil, err
			}
		}

		index++
		if index > end {
			break
		}
	}

	if selected == 0 {
		return nil, os.ErrNotExist
	}
	return out.Bytes(), nil
}

// GetThumbnails godoc
// @Summary      Fetch document thumbnails file
// @Description  Returns the full thumbnails file in PTF format.
// @Tags         pdf
// @Param        id    path string true "Document ID"
// @Produce      application/octet-stream
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "document not found"
// @Failure      500 {string} string "failed to read thumbnails"
// @Router       /pdf/thumbnails/{id} [get]
// @Security     BearerAuth
func GetThumbnails(c *gin.Context) {
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

	thumbPath := getThumbPath(doc.File.Path)
	data, err := os.ReadFile(thumbPath)
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
	data, err := readThumbnailRange(thumbPath, start, end)
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
