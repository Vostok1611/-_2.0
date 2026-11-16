package db

import (
	"GOBAR/internal/taskService"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=yourpassword dbname=gobar port=5433 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := database.AutoMigrate(&taskService.Task{}); err != nil {
		return nil, err
	}

	return database, nil
}
