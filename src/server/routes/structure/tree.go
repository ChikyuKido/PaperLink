package structure

import (
	"net/http"

	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type FileNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type DirNode struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Files       []FileNode `json:"files"`
	Directories []DirNode  `json:"directories"`
}

// Tree godoc
// @Summary      Get directory tree
// @Description  Returns the full directory and document tree for the authenticated user.
// @Tags         structure
// @Produce      json
// @Success      200 {object} DirNode
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/structure/tree [get]
// @Security     BearerAuth
func Tree(c *gin.Context) {
	userID := c.GetInt("userId")

	directories, err := repo.Directory.GetAllByUserId(userID)
	if err != nil {
		log.Errorf("failed to fetch directories for user %d: %v", userID, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to fetch directories")
		return
	}

	documents, err := repo.Document.GetAllByUserIdWithFile(userID)
	if err != nil {
		log.Errorf("failed to fetch documents for user %d: %v", userID, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to fetch documents")
		return
	}

	directoriesMap := make(map[int][]entity.Directory)
	for _, directory := range directories {
		parentID := 0
		if directory.ParentID != nil {
			parentID = *directory.ParentID
		}
		directoriesMap[parentID] = append(directoriesMap[parentID], directory)
	}

	documentsMap := make(map[int][]entity.Document)
	for _, document := range documents {
		directoryID := 0
		if document.DirectoryID != nil {
			directoryID = *document.DirectoryID
		}
		documentsMap[directoryID] = append(documentsMap[directoryID], document)
	}

	tree := buildTree(
		"",
		0,
		documentsMap[0],
		directoriesMap[0],
		documentsMap,
		directoriesMap,
	)

	routes.JSONSuccess(c, http.StatusOK, tree)
}

func mapDocuments(docs []entity.Document) []FileNode {
	files := make([]FileNode, 0, len(docs))

	for _, d := range docs {
		files = append(files, FileNode{
			ID:   d.UUID,
			Name: d.Name,
			Size: d.File.Size,
		})
	}

	return files
}

func buildTree(
	name string,
	id int,
	documents []entity.Document,
	directories []entity.Directory,
	docMap map[int][]entity.Document,
	dirMap map[int][]entity.Directory,
) DirNode {

	node := DirNode{
		ID:          id,
		Name:        name,
		Files:       mapDocuments(documents),
		Directories: make([]DirNode, 0),
	}

	for _, directory := range directories {
		child := buildTree(
			directory.Name,
			directory.ID,
			docMap[directory.ID],
			dirMap[directory.ID],
			docMap,
			dirMap,
		)

		node.Directories = append(node.Directories, child)
	}

	return node
}
