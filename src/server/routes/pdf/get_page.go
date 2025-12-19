package pdf

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"strconv"
	"strings"

	"paperlink/db/repo"
	"paperlink/pvf"

	"github.com/gin-gonic/gin"
)

// GetPage godoc
// @Summary      Fetch PDF page or page range
// @Description  Returns a custom binary format:
// @Description  flag byte (0=single,1=multi). Single: raw page bytes follow.
// @Description  Multi: repeated blocks of uint64 size + page bytes.
// @Tags         pdf
// @Param        id    path string true "Document UUID"
// @Param        page  path string true "Page or range (e.g. 3 or 2-5)"
// @Produce      application/octet-stream
// @Failure      400 {string} string "invalid page or format"
// @Failure      403 {string} string "forbidden"
// @Failure      404 {string} string "document not found"
// @Failure      500 {string} string "failed to read page(s)"
// @Router       /pdf/{id}/{page} [get]
// @Security     BearerAuth
func GetPage(c *gin.Context) {
	docUUID := c.Param("id")
	pageParam := c.Param("page")
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

	file := doc.FileDocument

	var from, to int

	if strings.Contains(pageParam, "-") {
		parts := strings.SplitN(pageParam, "-", 2)
		f, err1 := strconv.Atoi(parts[0])
		t, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || f <= 0 || t < f {
			c.String(http.StatusBadRequest, "invalid page range")
			return
		}
		from, to = f, t
	} else {
		n, err := strconv.Atoi(pageParam)
		if err != nil || n <= 0 {
			c.String(http.StatusBadRequest, "invalid page")
			return
		}
		from, to = n, n
	}

	if from > int(file.Pages) {
		c.String(http.StatusBadRequest, "page out of range")
		return
	}
	if to > int(file.Pages) {
		to = int(file.Pages)
	}

	var out bytes.Buffer

	if from == to {
		data, err := pvf.ReadPage(file.Path, uint64(from))
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to read page")
			return
		}

		out.WriteByte(0)
		out.Write(data)

		c.Data(http.StatusOK, "application/octet-stream", out.Bytes())
		return
	}

	pages, err := pvf.ReadPages(file.Path, uint64(from), uint64(to))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to read pages")
		return
	}

	out.WriteByte(1)

	for _, data := range pages {
		size := uint64(len(data))
		var header [8]byte
		binary.BigEndian.PutUint64(header[:], size)

		out.Write(header[:])
		out.Write(data)
	}

	c.Data(http.StatusOK, "application/octet-stream", out.Bytes())
}
