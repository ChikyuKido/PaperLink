package document

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"
	"strings"

	"github.com/gin-gonic/gin"
)

type FilterDocumentItem struct {
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DirectoryID *int     `json:"directoryId"`
	Tags        []string `json:"tags"`

	FileUUID string `json:"fileUUID"`
	Pages    uint64 `json:"pages"`
	Size     uint64 `json:"size"`
}

func Filter(c *gin.Context) {
	userID := c.GetInt("userId")

	tagsParam := strings.TrimSpace(c.Query("tags"))
	search := strings.TrimSpace(c.Query("search"))

	var tags []string
	if tagsParam != "" {
		parts := strings.Split(tagsParam, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				tags = append(tags, p)
			}
		}
	}

	docs, err := repo.Document.Filter(userID, tags, search)
	if err != nil {
		log.Errorf("filter query failed: %v", err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to filter documents")
		return
	}

	out := make([]FilterDocumentItem, 0, len(docs))
	for _, d := range docs {
		tagNames := make([]string, 0, len(d.Tags))
		for _, t := range d.Tags {
			tagNames = append(tagNames, t.Name)
		}

		out = append(out, FilterDocumentItem{
			UUID:        d.UUID,
			Name:        d.Name,
			Description: d.Description,
			DirectoryID: d.DirectoryID,
			Tags:        tagNames,
			FileUUID:    d.FileUUID,
			Pages:       d.File.Pages,
			Size:        d.File.Size,
		})
	}

	routes.JSONSuccess(c, http.StatusOK, out)
}
