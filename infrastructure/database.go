package infrastructure

import (
	"github.com/cndrsdrmn/ihsan-assessment/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entities.Blog{}); err != nil {
		return nil, err
	}

	return db, nil
}
