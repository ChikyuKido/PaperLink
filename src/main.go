package main

import (
	"paperlink/db"
	"paperlink/server"
)

func main() {
	db.Init()
	server.Start()
}
