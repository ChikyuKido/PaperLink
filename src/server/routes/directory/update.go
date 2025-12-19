package directory

import (
	"net/http"
	"strconv"

	"paperlink/db/repo"
	"paperlink/server/routes"

	"github.com/gin-gonic/gin"
)

type UpdateDirectoryRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parentId"`
}

// Update godoc
// @Summary      Update directory
// @Description  Updates a directory owned by the authenticated user.
// @Tags         directory
// @Accept       json
// @Produce      json
// @Param        id      path      int                    true  "Directory ID"
// @Param        request body      UpdateDirectoryRequest true  "Update directory payload"
// @Success      204     "No Content"
// @Failure      400     {object}  routes.ErrorResponse "Invalid request"
// @Failure      401     {object}  routes.ErrorResponse "Unauthorized"
// @Failure      403     {object}  routes.ErrorResponse "Forbidden"
// @Failure      404     {object}  routes.ErrorResponse "Directory not found"
// @Failure      409     {object}  routes.ErrorResponse "Invalid parent directory"
// @Failure      500     {object}  routes.ErrorResponse "Internal server error"
// @Router       /api/v1/directories/{id} [patch]
// @Security     BearerAuth
func Update(c *gin.Context) {
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
		routes.JSONError(c, http.StatusForbidden, "not authorized to update this directory")
		return
	}

	var req UpdateDirectoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("invalid update directory body: %v", err)
		routes.JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name != nil {
		dir.Name = *req.Name
	}

	if req.ParentID != nil {
		newParentID := *req.ParentID
		if newParentID == id {
			routes.JSONError(c, http.StatusConflict, "directory cannot be its own parent")
			return
		}
		parent, err := repo.Directory.Get(newParentID)
		if err != nil || parent == nil {
			routes.JSONError(c, http.StatusConflict, "parent directory not found")
			return
		}
		if parent.UserID != userID {
			routes.JSONError(c, http.StatusForbidden, "parent directory not owned by user")
			return
		}
		// do not allow a directory cycle
		for p := parent; p.ParentID != nil; {
			if *p.ParentID == id {
				routes.JSONError(c, http.StatusConflict, "directory cycle detected")
				return
			}
			p, err = repo.Directory.Get(*p.ParentID)
			if err != nil || p == nil {
				break
			}
		}

		dir.ParentID = req.ParentID
	}

	if err := repo.Directory.Save(dir); err != nil {
		log.Errorf("failed to update directory %d: %v", id, err)
		routes.JSONError(c, http.StatusInternalServerError, "failed to update directory")
		return
	}

	c.Status(http.StatusNoContent)
}
