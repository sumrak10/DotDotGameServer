package database

import (
	"OnlineGame/internal/config"
	"fmt"
	"log"
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
		dbPath := fmt.Sprintf("db/%s", config.Database().Path)
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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
