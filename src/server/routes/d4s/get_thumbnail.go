package d4s

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"paperlink/db/repo"

	"github.com/gin-gonic/gin"
)

const d4sThumbCacheDir = "./data/d4s/thumbs"
const d4sThumbDPI = "80"
const d4sThumbWebPQuality = "60"

func d4sThumbCachePath(uuid string) string {
	return path.Join(d4sThumbCacheDir, uuid+".webp")
}

func ensureD4SThumbnail(uuid string, pdfPath string) (string, error) {
	thumbPath := d4sThumbCachePath(uuid)
	if st, err := os.Stat(thumbPath); err == nil && !st.IsDir() {
		return thumbPath, nil
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if err := os.MkdirAll(d4sThumbCacheDir, 0o750); err != nil {
		return "", err
	}

	stamp := time.Now().UnixNano()
	tmpPNGPath := fmt.Sprintf("%s.%d.tmp.png", thumbPath, stamp)
	tmpWebPPath := fmt.Sprintf("%s.%d.tmp.webp", thumbPath, stamp)
	_ = os.Remove(tmpPNGPath)
	_ = os.Remove(tmpWebPPath)

	gsCmd := exec.Command(
		"gs",
		"-sDEVICE=png16m",
		"-r"+d4sThumbDPI,
		"-dFirstPage=1",
		"-dLastPage=1",
		"-dTextAlphaBits=4",
		"-dGraphicsAlphaBits=4",
		"-dBATCH",
		"-dNOPAUSE",
		"-sOutputFile="+tmpPNGPath,
		pdfPath,
	)
	out, err := gsCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ghostscript failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	webpCmd := exec.Command("cwebp", "-quiet", "-q", d4sThumbWebPQuality, tmpPNGPath, "-o", tmpWebPPath)
	out, err = webpCmd.CombinedOutput()
	if err != nil {
		_ = os.Remove(tmpPNGPath)
		_ = os.Remove(tmpWebPPath)
		return "", fmt.Errorf("cwebp failed: %w: %s", err, strings.TrimSpace(string(out)))
	}

	_ = os.Remove(tmpPNGPath)

	if err := os.Rename(tmpWebPPath, thumbPath); err != nil {
		_ = os.Remove(tmpWebPPath)
		return "", err
	}
	return thumbPath, nil
}

// GetThumbnail godoc
// @Summary      Get Digi4School thumbnail
// @Description  Returns the first thumbnail page as WebP for a Digi4School book.
// @Tags         digi4school
// @Produce      image/webp
// @Param        id   path      int  true  "Digi4School book ID"
// @Failure      400  {object}  routes.ErrorResponse "Invalid book ID"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      404  {object}  routes.ErrorResponse "Book not found"
// @Failure      500  {object}  routes.ErrorResponse "Failed to load thumbnail"
// @Router       /api/v1/d4s/thumbnail/{id} [get]
// @Security     BearerAuth
func GetThumbnail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid book id")
		return
	}

	book, err := repo.Digi4SchoolBook.Get(id)
	if err != nil || book == nil {
		c.String(http.StatusNotFound, "book not found")
		return
	}
	if book.FileUUID == "" {
		c.String(http.StatusNotFound, "book has no file")
		return
	}

	file := repo.FileDocument.GetByUUID(book.FileUUID)
	if file == nil {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	thumbPath, err := ensureD4SThumbnail(book.FileUUID, file.Path)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to create thumbnail")
		return
	}

	data, err := os.ReadFile(thumbPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to read thumbnail")
		return
	}

	c.Data(http.StatusOK, "image/webp", data)
}
