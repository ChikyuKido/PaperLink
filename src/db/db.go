package db

import (
	"log"
	"paperlink/db/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open("paperlink.db"), &gorm.Config{})
	_ = db.Exec(`
		PRAGMA journal_mode=WAL;
	`)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	err = db.AutoMigrate(&entity.Annotation{}, &entity.AnnotationAction{}, &entity.Category{}, &entity.Document{},
		&entity.DocumentUser{}, &entity.Notification{}, &entity.Tag{}, &entity.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
		return
	}

	DB = db
}
