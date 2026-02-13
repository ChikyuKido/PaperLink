package document

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"paperlink/db/entity"
	"paperlink/db/repo"
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
// @Description  Uploads a PDF, converts it to PVF, and stores it.
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
	uploadStart := time.Now()
	f, err := c.FormFile("file")
	if err != nil {
		log.Warnf("failed to get file from request: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "failed to read uploaded file")
		return
	}
	log.Infof("upload start filename=%s size=%dB", f.Filename, f.Size)

	fileUUID := uuid.New().String()
	tmpDst := "./data/tmp/uploads/" + fileUUID

	saveStart := time.Now()
	if err := c.SaveUploadedFile(f, tmpDst); err != nil {
		log.Errorf("failed to save uploaded file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to upload file")
		return
	}
	log.Infof("upload saved tmp=%s took=%s", tmpDst, time.Since(saveStart))

	pvfStart := time.Now()
	pvfFile, err := pvf.WritePVFFromPDF(tmpDst)
	if err != nil {
		log.Errorf("failed to convert pdf to pvf: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	log.Infof("upload pvf conversion done path=%s took=%s", pvfFile, time.Since(pvfStart))

	ptfStart := time.Now()
	thumbPTFFile, err := pvf.WriteThumbnailPTFFromPDF(tmpDst)
	if err != nil {
		log.Errorf("failed to convert pdf thumbnails to ptf: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to process file")
		return
	}
	log.Infof("upload thumbnail conversion done path=%s took=%s", thumbPTFFile, time.Since(ptfStart))

	if err := os.MkdirAll("./data/uploads", 0750); err != nil {
		log.Errorf("failed to create upload dir: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}

	copyPVFStart := time.Now()
	dst := "./data/uploads/" + fileUUID + ".pvf"
	if err := util.CopyFile(pvfFile, dst); err != nil {
		log.Errorf("failed to copy pvf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}
	log.Infof("upload copied pvf dst=%s took=%s", dst, time.Since(copyPVFStart))

	copyPTFStart := time.Now()
	thumbDst := "./data/uploads/" + fileUUID + "_thumb.ptf"
	if err := util.CopyFile(thumbPTFFile, thumbDst); err != nil {
		log.Errorf("failed to copy thumbnail ptf file: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to store file")
		return
	}
	log.Infof("upload copied ptf dst=%s took=%s", thumbDst, time.Since(copyPTFStart))

	_ = os.RemoveAll(filepath.Dir(pvfFile))
	_ = os.RemoveAll(filepath.Dir(thumbPTFFile))
	_ = os.Remove(tmpDst)

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
	if thumbStat, err := os.Stat(thumbDst); err == nil {
		log.Infof("upload sizes pvf=%dB ptf=%dB pages=%d", stat.Size(), thumbStat.Size(), metadata.PageCount)
	} else {
		log.Infof("upload sizes pvf=%dB pages=%d (ptf stat failed: %v)", stat.Size(), metadata.PageCount, err)
	}
	log.Infof("upload done fileUUID=%s total=%s", fileUUID, time.Since(uploadStart))

	routes.JSONSuccess(c, http.StatusOK, UploadDocumentResponse{
		FileUUID: fileUUID,
	})
}
