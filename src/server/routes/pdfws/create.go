package pdfws

import (
	"errors"
	"net/http"

	"paperlink/server/routes"
	"paperlink/service/collabedit"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		routes.JSONError(c, http.StatusBadRequest, "document id required")
		return
	}

	userID := c.GetInt("userId")
	result, err := collabedit.PDFCollab.CreateSingleUseToken(documentID, userID)
	if err != nil {
		switch {
		case errors.Is(err, collabedit.ErrDocumentNotFound):
			routes.JSONError(c, http.StatusNotFound, err.Error())
		case errors.Is(err, collabedit.ErrForbidden):
			routes.JSONError(c, http.StatusForbidden, err.Error())
		default:
			log.Errorf("failed to create websocket token for document %s: %v", documentID, err)
			routes.JSONError(c, http.StatusInternalServerError, "failed to create websocket token")
		}
		return
	}

	routes.JSONSuccessOK(c, result)
}
