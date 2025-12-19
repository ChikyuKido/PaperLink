package document

import (
	"math/rand"
	"net/http"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	DirectoryID *int     `json:"directoryId"`
	Tags        []string `json:"tags"`
	FileUUID    string   `json:"fileUUID" binding:"required"`
}

type DocumentResponse struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DirectoryID *int     `json:"directoryId"`
	Tags        []string `json:"tags"`
}

// Create godoc
// @Summary      Create document
// @Description  Creates a new document for the authenticated user.
// @Tags         document
// @Accept       json
// @Produce      json
// @Param        request body CreateRequest true "Create document payload"
// @Success      201 {object} DocumentResponse
// @Failure      400 {object} routes.ErrorResponse "Invalid request"
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/documents/create [post]
// @Security     BearerAuth
func Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid create document body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetInt("userId")

	if req.DirectoryID != nil {
		dir, err := repo.Directory.Get(*req.DirectoryID)
		if err != nil || dir == nil {
			routes.JSONError(c, http.StatusBadRequest, "invalid directory")
			return
		}
		if dir.UserID != userID {
			routes.JSONError(c, http.StatusForbidden, "directory not owned by user")
			return
		}
	}

	file := repo.FileDocument.GetByUUID(req.FileUUID)
	if file == nil {
		routes.JSONError(c, http.StatusBadRequest, "file does not exist")
		return
	}

	doc := entity.Document{
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
		FileUUID:    req.FileUUID,
		DirectoryID: req.DirectoryID,
	}

	if len(req.Tags) > 0 {
		dbTags, err := repo.Tag.GetList()
		if err != nil {
			log.Errorf("failed to fetch tags: %v", err)
			routes.JSONError(c, http.StatusInternalServerError, "failed to process tags")
			return
		}

		existing := make(map[string]entity.Tag, len(dbTags))
		for _, t := range dbTags {
			existing[t.Name] = t
		}

		finalTags := make([]entity.Tag, 0, len(req.Tags))

		for _, name := range req.Tags {
			if t, ok := existing[name]; ok {
				finalTags = append(finalTags, t)
				continue
			}

			newTag := entity.Tag{
				Name:  name,
				Color: entity.TagColors[rand.Intn(len(entity.TagColors))],
			}

			if err := repo.Tag.Save(&newTag); err != nil {
				log.Errorf("failed to create tag %s: %v", name, err)
				routes.JSONError(c, http.StatusInternalServerError, "failed to create tag")
				return
			}

			existing[name] = newTag
			finalTags = append(finalTags, newTag)
		}

		doc.Tags = finalTags
	}

	if err := repo.Document.Save(&doc); err != nil {
		log.Errorf("failed to create document: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to create document")
		return
	}

	tagNames := make([]string, 0, len(doc.Tags))
	for _, t := range doc.Tags {
		tagNames = append(tagNames, t.Name)
	}

	routes.JSONSuccess(c, http.StatusCreated, DocumentResponse{
		UUID:        doc.UUID,
		Name:        doc.Name,
		Description: doc.Description,
		DirectoryID: doc.DirectoryID,
		Tags:        tagNames,
	})
}
