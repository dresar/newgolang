package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gowavuelite/models"
)

var DB *gorm.DB

func ConnectDatabase() error {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.AutoReply{},
		&models.MessageLog{},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

