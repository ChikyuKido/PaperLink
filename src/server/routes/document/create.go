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
	Path        string   `json:"path"`
	Tags        []string `json:"tags"`
}

func Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, routes.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	doc := entity.Document{
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Path:        req.Path,
		OwnerID:     c.GetInt("userId"),
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

			newTag := entity.Tag{Name: name, Color: entity.TagColors[rand.Intn(len(entity.TagColors))]}
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

	c.JSON(http.StatusOK, doc)
}
