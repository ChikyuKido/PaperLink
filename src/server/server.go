package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"paperlink/server/routes/auth"
)

func Start() {
	r := gin.Default()

	auth.InitAuthRouter(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
