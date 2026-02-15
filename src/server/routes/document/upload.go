package document

import (
	"net/http"
	"os"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/ptf"
	"paperlink/server/routes"
	"paperlink/util"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type UploadDocumentResponse struct {
	FileUUID string `json:"fileUUID"`
}

// Upload godoc
// @Summary      Upload document file
// @Description  Uploads a PDF, generates thumbnail PTF, and stores the PDF.
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
	tmpDst := "./data/tmp/uploads/" + fileUUID

	if err := c.SaveUploadedFile(f, tmpDst); err != nil {
		log.Errorf("failed to save uploaded file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to upload file")
		return
	}
	thumbPTFFile, err := ptf.WriteThumbnailPTFFromPDF(tmpDst)
	if err != nil {
		log.Errorf("failed to convert pdf thumbnails to ptf: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	if err := os.MkdirAll("./data/uploads", 0750); err != nil {
		log.Errorf("failed to create upload dir: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	dst := "./data/uploads/" + fileUUID + ".pdf"
	if err := util.CopyFile(tmpDst, dst); err != nil {
		log.Errorf("failed to copy pdf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	thumbDst := "./data/uploads/" + fileUUID + "_thumb.ptf"
	if err := util.CopyFile(thumbPTFFile, thumbDst); err != nil {
		log.Errorf("failed to copy thumbnail ptf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	_ = os.RemoveAll(filepath.Dir(thumbPTFFile))
	_ = os.Remove(tmpDst)

	stat, err := os.Stat(dst)
	if err != nil {
		log.Errorf("failed to stat pdf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	pageCount, err := api.PageCountFile(dst)
	if err != nil {
		log.Errorf("failed to read pdf page count: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to read file metadata")
		return
	}

	if err := repo.FileDocument.Save(&entity.FileDocument{
		UUID:  fileUUID,
		Path:  dst,
		Size:  uint64(stat.Size()),
		Pages: uint64(pageCount),
	}); err != nil {
		log.Errorf("failed to save file document: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	routes.JSONSuccess(c, http.StatusOK, UploadDocumentResponse{
		FileUUID: fileUUID,
	})
}
