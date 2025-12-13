package document

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/pvf"
	"paperlink/server/routes"
	"paperlink/util"
	"path/filepath"
)

func Upload(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, routes.NewError(400, fmt.Sprintf("failed to get file from request: %v", err)))
		return
	}
	fileUUID := uuid.New().String()
	tmpDst := "./data/tmp/uploads/" + fileUUID
	if err := c.SaveUploadedFile(f, tmpDst); err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to upload file: %v", err)))
		return
	}
	pvfFile, err := pvf.WritePVFFromPDF(tmpDst)
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to convert pvfFile file: %v", err)))
		return
	}
	err = os.MkdirAll("./data/uploads", 0750)
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to make dir: %v", err)))
		return
	}
	err = util.CopyFile(pvfFile, "./data/uploads/"+fileUUID+".pvf")
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to copy file: %v", err)))
		return
	}
	err = os.RemoveAll(filepath.Dir(pvfFile))
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to remove file: %v", err)))
		return
	}
	err = os.Remove(tmpDst)
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to remove file: %v", err)))
		return
	}
	stat, err := os.Stat("./data/uploads/" + fileUUID + ".pvf")
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to stat file: %v", err)))
		return
	}
	metadata, err := pvf.ReadMetadata("./data/uploads/" + fileUUID + ".pvf")
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to read metadata file: %v", err)))
		return
	}
	err = repo.FileDocument.Save(&entity.FileDocument{
		UUID:  fileUUID,
		Path:  "./data/uploads/" + fileUUID + ".pvf",
		Size:  uint64(stat.Size()),
		Pages: metadata.PageCount,
	})
	if err != nil {
		c.JSON(500, routes.NewError(500, fmt.Sprintf("failed to save file: %v", err)))
		return
	}

	c.JSON(200, routes.NewSuccess(gin.H{
		"id": fileUUID,
	}))
}
