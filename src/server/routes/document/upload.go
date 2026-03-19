package document

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/ptf"
	"paperlink/pvf"
	"paperlink/server/routes"
	"paperlink/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadDocumentResponse struct {
	FileUUID string `json:"fileUUID"`
}

// Upload godoc
// @Summary      Upload document file
// @Description  Uploads a PDF, converts it to PVF, generates thumbnail PTF, and stores it.
// @Tags         document
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "PDF file"
// @Success      200 {object} UploadDocumentResponse
// @Failure      400 {object} routes.ErrorResponse "Invalid upload"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/documents/upload [post]
// @Security     BearerAuth
func Upload(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		log.Warnf("failed to get file from request: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "failed to read uploaded file")
		return
	}

	fileUUID := uuid.New().String()

	if err := os.MkdirAll("./data/tmp/uploads", 0750); err != nil {
		log.Errorf("failed to create temp upload dir: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to prepare upload")
		return
	}

	if err := os.MkdirAll("./data/uploads", 0750); err != nil {
		log.Errorf("failed to create upload dir: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	tmpSrc := "./data/tmp/uploads/" + fileUUID + ".pdf"
	tmpLinearized := "./data/tmp/uploads/" + fileUUID + ".linearized.pdf"
	dst := "./data/uploads/" + fileUUID + ".pvf"
	thumbDst := "./data/uploads/" + fileUUID + "_thumb.ptf"

	cleanupPaths := []string{tmpSrc, tmpLinearized}
	defer func() {
		for _, p := range cleanupPaths {
			_ = os.Remove(p)
		}
	}()

	if err := c.SaveUploadedFile(f, tmpSrc); err != nil {
		log.Errorf("failed to save uploaded file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to upload file")
		return
	}

	cmd := exec.Command(
		"qpdf",
		"--warning-exit-0",
		"--linearize",
		"--object-streams=generate",
		"--stream-data=compress",
		tmpSrc,
		tmpLinearized,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("failed to process pdf with qpdf: %v, output: %s", err, string(output))
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}

	viewPVFFile, err := pvf.WritePVFFromPDF(tmpLinearized)
	if err != nil {
		log.Errorf("failed to convert pdf to pvf: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	defer func() {
		_ = os.RemoveAll(filepath.Dir(viewPVFFile))
	}()

	if err := util.CopyFile(viewPVFFile, dst); err != nil {
		log.Errorf("failed to copy pvf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	thumbPTFFile, err := ptf.WriteThumbnailPTFFromPDF(tmpLinearized)
	if err != nil {
		log.Errorf("failed to convert pdf thumbnails to ptf: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	defer func() {
		_ = os.RemoveAll(filepath.Dir(thumbPTFFile))
	}()

	if err := util.CopyFile(thumbPTFFile, thumbDst); err != nil {
		log.Errorf("failed to copy thumbnail ptf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	stat, err := os.Stat(dst)
	if err != nil {
		log.Errorf("failed to stat pvf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	metadata, err := pvf.ReadMetadata(dst)
	if err != nil {
		log.Errorf("failed to read pvf metadata: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to read file metadata")
		return
	}

	if err := repo.FileDocument.Save(&entity.FileDocument{
		UUID:  fileUUID,
		Path:  dst,
		Size:  uint64(stat.Size()),
		Pages: metadata.PageCount,
	}); err != nil {
		log.Errorf("failed to save file document: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	routes.JSONSuccess(c, http.StatusOK, UploadDocumentResponse{
		FileUUID: fileUUID,
	})
}
