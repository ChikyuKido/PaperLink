package pdf

import (
	"bytes"
	"encoding/binary"
	"github.com/gin-gonic/gin"
	"paperlink/db/repo"
	"paperlink/pvf"
	"strconv"
	"strings"
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
// @Failure      404 {string} string "document not found"
// @Failure      500 {string} string "failed to read page(s)"
// @Router       /pdf/{id}/{page} [get]
func GetPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		page := c.Param("page")

		doc := repo.FileDocument.GetByUUID(id)
		if doc == nil {
			c.String(404, "document not found")
			return
		}

		var from, to int
		if strings.Contains(page, "-") {
			parts := strings.SplitN(page, "-", 2)
			f, err1 := strconv.Atoi(parts[0])
			t, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				c.String(400, "invalid page format")
				return
			}
			from, to = f, t
		} else {
			n, err := strconv.Atoi(page)
			if err != nil {
				c.String(400, "invalid page")
				return
			}
			from, to = n, n
		}

		var out bytes.Buffer
		if from == to {
			data, err := pvf.ReadPage(doc.Path, uint64(from))
			if err != nil {
				c.String(500, "failed to read page")
				return
			}

			out.WriteByte(0)
			out.Write(data)

			c.Data(200, "application/octet-stream", out.Bytes())
			return
		}
		pages, err := pvf.ReadPages(doc.Path, uint64(from), uint64(to))
		if err != nil {
			c.String(500, "failed to read pages")
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

		c.Data(200, "application/octet-stream", out.Bytes())
	}
}
