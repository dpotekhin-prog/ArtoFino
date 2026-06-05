package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres opens a GORM connection to Postgres.
func ConnectPostgres(dsn string) (*gorm.DB, error) {
	log.Println("Connecting to Postgres...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
