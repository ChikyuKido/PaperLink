package main

import (
	"github.com/sirupsen/logrus"
	"paperlink/db"
	"paperlink/server"
	"paperlink/service/task"
	"paperlink/util"
)

func main() {
	logrus.SetFormatter(&util.GroupFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	db.DB()
	task.Init()
	server.Start()
}
