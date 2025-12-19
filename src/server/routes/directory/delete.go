package directory

import (
	"net/http"
	"strconv"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary      Delete directory
// @Description  Deletes a directory owned by the authenticated user.
// @Tags         directory
// @Produce      json
// @Param        id   path      int  true  "Directory ID"
// @Success      204  "No Content"
// @Failure      400  {object}  routes.ErrorResponse "Invalid directory ID"
// @Failure      401  {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403  {object}  routes.ErrorResponse "Forbidden"
// @Failure      404  {object}  routes.ErrorResponse "Directory not found"
// @Failure      500  {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/directories/{id} [delete]
// @Security     BearerAuth
func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		routes.JSONError(c, http.StatusBadRequest, "invalid directory id")
		return
	}

	userID := c.GetInt("userId")

	dir, err := repo.Directory.Get(id)
	if err != nil || dir == nil {
		routes.JSONError(c, http.StatusNotFound, "directory not found")
		return
	}

	if dir.UserID != userID {
		routes.JSONError(c, http.StatusForbidden, "not authorized to delete this directory")
		return
	}

	if err := deleteDirectoryTree(id); err != nil {
		log.Errorf("failed to recursively delete directory %d: %v", id, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to delete directory")
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteDirectoryTree(dirID int) error {
	children, err := repo.Directory.GetChildren(dirID)
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := deleteDirectoryTree(child.ID); err != nil {
			return err
		}
	}

	return repo.Directory.Delete(dirID)
}
