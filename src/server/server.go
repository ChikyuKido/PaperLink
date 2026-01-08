package server

import (
	"paperlink/server/routes/auth"
	"paperlink/server/routes/digi4school"
	"paperlink/server/routes/directory"
	"paperlink/server/routes/document"
	"paperlink/server/routes/invite"
	"paperlink/server/routes/pdf"
	"paperlink/server/routes/structure"
	"paperlink/util"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("SERVER")

func Start() {
	r := gin.Default()
	auth.InitAuthRouter(r)
	pdf.InitPDFRouter(r)
	document.InitDocumentRouter(r)
	invite.InitInviteRouter(r)
	directory.InitDirectoryRouter(r)
	structure.InitStructureRoutes(r)
	digi4school.InitDigi4SchoolRouter(r)
	log.Info("starting server at port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
