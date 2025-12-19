package structure

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paperlink/db/entity"
	"paperlink/db/repo"
	"paperlink/server/routes"
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

func Tree(c *gin.Context) {
	userID := c.GetInt("userId")

	directories, err := repo.Directory.GetAllByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to fetch directories"))
		return
	}
	documents, err := repo.Document.GetAllByUserId(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, routes.NewError(500, "failed to fetch documents"))
		return
	}
	directoriesMap := make(map[int][]entity.Directory)
	for _, directory := range directories {
		parentID := 0
		if directory.ParentID != nil {
			parentID = *directory.ParentID
		}
		if _, ok := directoriesMap[parentID]; !ok {
			directoriesMap[parentID] = []entity.Directory{}
		}
		directoriesMap[parentID] = append(directoriesMap[parentID], directory)
	}
	documentsMap := make(map[int][]entity.Document)
	for _, document := range documents {
		directoryID := 0
		if document.DirectoryID != nil {
			directoryID = *document.DirectoryID
		}
		if _, ok := documentsMap[directoryID]; !ok {
			documentsMap[directoryID] = []entity.Document{}
		}
		documentsMap[directoryID] = append(documentsMap[directoryID], document)
	}
	c.JSON(http.StatusOK, routes.NewSuccess(
		buildTree("", 0, documentsMap[0], directoriesMap[0], documentsMap, directoriesMap)))
}
func mapDocuments(docs []entity.Document) []FileNode {
	files := make([]FileNode, 0, len(docs))

	for _, d := range docs {
		files = append(files, FileNode{
			ID:   d.UUID,
			Name: d.Name,
			Size: d.FileDocument.Size,
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
		childrenDirs := dirMap[directory.ID]
		childrenDocs := docMap[directory.ID]

		child := buildTree(
			directory.Name,
			directory.ID,
			childrenDocs,
			childrenDirs,
			docMap,
			dirMap,
		)

		node.Directories = append(node.Directories, child)
	}

	return node
}
