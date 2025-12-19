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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DirectoryID *int     `json:"directoryId"`
	Tags        []string `json:"tags"`
	FileUUID    string   `json:"fileUUID"`
}

func Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	userID := c.GetInt("userId")

	if req.DirectoryID != nil {
		if _, err := repo.Directory.Get(*req.DirectoryID); err != nil {
			c.JSON(http.StatusBadRequest, routes.NewError(http.StatusBadRequest, "invalid directory"))
			return
		}
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
			c.JSON(http.StatusInternalServerError, routes.NewError(http.StatusInternalServerError, err.Error()))
			return
		}

		existing := make(map[string]entity.Tag, len(dbTags))
		for _, t := range dbTags {
			existing[t.Name] = t
		}

		var finalTags []entity.Tag

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
				c.JSON(http.StatusInternalServerError, routes.NewError(http.StatusInternalServerError, err.Error()))
				return
			}

			existing[name] = newTag
			finalTags = append(finalTags, newTag)
		}

		doc.Tags = finalTags
	}

	if err := repo.Document.Save(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, routes.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, routes.NewSuccess(gin.H{"msg": "ok"}))
}
