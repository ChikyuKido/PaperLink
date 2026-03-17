package pdfws

import (
	"errors"
	"net/http"

	"paperlink/server/routes"
	"paperlink/service/collabedit"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func Connect(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		routes.JSONError(c, http.StatusBadRequest, "document id required")
		return
	}

	token := c.Query("token")
	if err := collabedit.PDFCollab.ValidateConnection(documentID, token); err != nil {
		switch {
		case errors.Is(err, collabedit.ErrTokenRequired):
			routes.JSONError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, collabedit.ErrTokenInvalid), errors.Is(err, collabedit.ErrTokenExpired):
			routes.JSONError(c, http.StatusUnauthorized, err.Error())
		default:
			log.Errorf("failed to validate websocket connection for document %s: %v", documentID, err)
			routes.JSONError(c, http.StatusInternalServerError, "failed to validate websocket token")
		}
		return
	}

	websocket.Handler(func(ws *websocket.Conn) {
		if err := collabedit.PDFCollab.HandleConnection(documentID, token, ws); err != nil {
			log.Warnf("websocket closed for document %s: %v", documentID, err)
		}
	}).ServeHTTP(c.Writer, c.Request)
}
