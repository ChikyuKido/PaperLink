package db

import (
	"fmt"
	"gorm.io/gorm/logger"
	"os"
	"paperlink/db/entity"
	"paperlink/util"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	instance *gorm.DB
)
var log = util.GroupLog("DATABASE")

func DB() *gorm.DB {
	once.Do(func() {
		err := os.MkdirAll("./data/log", 0755)
		if err != nil {
			logrus.Fatalf("Failed to create log directory: %v", err)
		}
		instance, err = gorm.Open(sqlite.Open("./data/app.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		instance.Exec("PRAGMA cache_size = -10240;")
		instance.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		err = instance.AutoMigrate(
			&entity.Annotation{}, &entity.AnnotationAction{}, &entity.FileDocument{},
			&entity.Document{}, &entity.DocumentUser{}, &entity.Notification{},
			&entity.Tag{}, &entity.User{},
		)
		if err != nil {
			log.Fatalf("Error migrating database: %v", err)
		}
		fmt.Println("Database connection established.")
	})

	return instance
}
