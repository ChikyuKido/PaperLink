package main

import (
	"github.com/sirupsen/logrus"
	"paperlink/db"
	"paperlink/server"
	"paperlink/util"
)

func main() {
	// Set custom formatter
	logrus.SetFormatter(&util.GroupFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	db.DB()
	server.Start()
}
