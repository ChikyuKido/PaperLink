package main

import (
	"paperlink/db"
	"paperlink/server"
	"paperlink/util"

	"github.com/sirupsen/logrus"
)

func main() {
	// Set custom formatter
	logrus.SetFormatter(&util.GroupFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	db.DB()
	server.Start()
}
