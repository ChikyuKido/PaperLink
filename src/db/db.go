package db

import (
	"fmt"
	"os"
	"paperlink/db/entity"
	"paperlink/util"
	"sync"

	"gorm.io/gorm/logger"

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
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		err = ApplySQLiteConfig(instance)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}
		err = instance.AutoMigrate(
			&entity.Annotation{}, &entity.AnnotationAction{}, &entity.FileDocument{},
			&entity.Document{}, &entity.DocumentUser{}, &entity.Notification{},
			&entity.Tag{}, &entity.User{}, &entity.Directory{},
			&entity.RegistrationInvite{},
		)
		if err != nil {
			log.Fatalf("Error migrating database: %v", err)
		}
		fmt.Println("Database connection established.")
	})

	return instance
}
func ApplySQLiteConfig(instance *gorm.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA cache_size = -10240;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA wal_autocheckpoint = 1000;",
	}

	for _, p := range pragmas {
		if err := instance.Exec(p).Error; err != nil {
			return err
		}
	}
	return nil
}
