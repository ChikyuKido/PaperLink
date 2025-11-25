package db

import (
	"paperlink/db/entity"
	"paperlink/util"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var log = util.GroupLog("DATABASE")

func Init() {
	db, err := gorm.Open(sqlite.Open("paperlink.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	_ = db.Exec(`
		PRAGMA journal_mode=WAL;
	`)

	err = db.AutoMigrate(
		&entity.Annotation{}, &entity.AnnotationAction{}, &entity.Category{},
		&entity.Document{}, &entity.DocumentUser{}, &entity.Notification{},
		&entity.Tag{}, &entity.User{},
	)
	if err != nil {
		logrus.WithError(err).Fatal("failed to migrate database")
		return
	}

	log.Info("database initialized successfully")
	DB = db
}
