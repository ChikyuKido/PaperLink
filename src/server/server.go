package server

import (
	"paperlink/server/routes/auth"
	"paperlink/server/routes/directory"
	"paperlink/server/routes/document"
	"paperlink/server/routes/pdf"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("SERVER")

func Start() {
	r := gin.Default()
	auth.InitAuthRouter(r)
	pdf.InitPDFRouter(r)
	document.InitDocumentRouter(r)
	directory.InitDirectoryRouter(r)
	log.Info("starting server at port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
