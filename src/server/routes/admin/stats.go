package admin

import (
	"net/http"
	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type StatsResponse struct {
	UserCount       int64  `json:"userCount"`
	DocumentCount   int64  `json:"documentCount"`
	TotalDocSize    uint64 `json:"totalDocSize"`
	TotalPages      uint64 `json:"totalPages"`
	D4SBookCount    int64  `json:"d4sBookCount"`
	D4SAccountCount int64  `json:"d4sAccountCount"`
}

// Stats godoc
// @Summary      Admin statistics
// @Description  Returns basic statistics like user count, document sizes/pages.
// @Tags         admin
// @Produce      json
// @Success      200 {object} StatsResponse
// @Failure      401 {object} routes.ErrorResponse "Unauthorized"
// @Failure      403 {object} routes.ErrorResponse "Forbidden"
// @Failure      500 {object} routes.ErrorResponse "Internal server error"
// @Router       /api/v1/admin/stats [get]
// @Security     BearerAuth
func Stats(c *gin.Context) {
	// Users
	users, err := repo.User.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list users")
		return
	}

	// Documents
	docs, err := repo.Document.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list documents")
		return
	}

	// Sum file sizes/pages (file_documents are the source of truth)
	fileDocs, err := repo.FileDocument.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list file documents")
		return
	}

	var totalSize uint64
	var totalPages uint64
	for _, f := range fileDocs {
		totalSize += f.Size
		totalPages += f.Pages
	}

	// D4S books/accounts
	d4sBooks, err := repo.Digi4SchoolBook.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list d4s books")
		return
	}
	accs, err := repo.Digi4SchoolAccount.GetList()
	if err != nil {
		routes.JSONError(c, http.StatusInternalServerError, "failed to list d4s accounts")
		return
	}

	routes.JSONSuccessOK(c, StatsResponse{
		UserCount:       int64(len(users)),
		DocumentCount:   int64(len(docs)),
		TotalDocSize:    totalSize,
		TotalPages:      totalPages,
		D4SBookCount:    int64(len(d4sBooks)),
		D4SAccountCount: int64(len(accs)),
	})
}
