package directory

import (
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type CreateDirectoryRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int   `json:"parentId"`
}

type CreateDirectoryResponse struct {
	ID int `json:"id"`
}

// Create godoc
// @Summary      Create directory
// @Description  Creates a new directory for the authenticated user.
// @Tags         directory
// @Accept       json
// @Produce      json
// @Param        request body CreateDirectoryRequest true "Create directory payload"
// @Success      201 {object} CreateDirectoryResponse
// @Failure      400 {object} routes.ErrorResponse "Invalid request body"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      409 {object} routes.ErrorResponse "Failed to create directory"
// @Router       /api/v1/directories [post]
// @Security     BearerAuth
func Create(c *gin.Context) {
	var req CreateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid create directory body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetInt("userId")

	dir := entity.Directory{
		UserID:   userID,
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if req.ParentID != nil {
		parent, err := repo.Directory.Get(*req.ParentID)
		if err != nil || parent.UserID != userID {
			routes.JSONError(c, http.StatusForbidden, "invalid parent directory")
			return
		}
	}

	if err := repo.Directory.Save(&dir); err != nil {
		log.Errorf("failed to save directory for user %d: %v", userID, err)
		routes.JSONError(c, http.StatusConflict, "failed to create directory")
		return
	}

	routes.JSONSuccess(c, http.StatusCreated, CreateDirectoryResponse{
		ID: dir.ID,
	})
}
