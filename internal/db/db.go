package db

import (
	"backend/internal/calculationService"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=pw dbname=postgres port=5433 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("couldnt connect to db: %v", err)
	}
	if err := db.AutoMigrate(&calculationService.Calculation{}); err != nil {
		log.Fatalf("couldnt migrate: %v", err)
	}
	return db, nil
}
