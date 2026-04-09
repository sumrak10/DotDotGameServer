package database

import (
	"log"
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

func GetDB() *gorm.DB {
	dbOnce.Do(func() {
		db, err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}
		err = db.AutoMigrate(
			&User{},
			&Match{},
		)
		if err != nil {
			return
		}
		dbInstance = db
	})
	return dbInstance
}
