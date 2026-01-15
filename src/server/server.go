package server

import (
	"paperlink/server/routes/admin"
	"paperlink/server/routes/auth"
	"paperlink/server/routes/d4s"
	"paperlink/server/routes/directory"
	"paperlink/server/routes/document"
	"paperlink/server/routes/invite"
	"paperlink/server/routes/pdf"
	"paperlink/server/routes/structure"
	"paperlink/server/routes/task"
	"paperlink/util"
	"strings"

	"github.com/gin-gonic/gin"
)

var log = util.GroupLog("SERVER")

func Start() {
	r := gin.New()
	r.Static("/assets", "./dist/assets")
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File("./dist/index.html")
	})

	auth.InitAuthRouter(r)
	admin.InitAdminRouter(r)
	pdf.InitPDFRouter(r)
	document.InitDocumentRouter(r)
	invite.InitInviteRouter(r)
	directory.InitDirectoryRouter(r)
	structure.InitStructureRoutes(r)
	d4s.InitDigi4SchoolRouter(r)
	task.InitTasksTasks(r)
	log.Info("starting server at port 8080")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
