package db

import (
	"log"

	"ArtoFino/backend/internal/models" // Импортируем наши модели

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgres opens a GORM connection to Postgres and ensures schemas are migrated.
func ConnectPostgres(dsn string) (*gorm.DB, error) {
	log.Println("Connecting to Postgres...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Заставляем GORM автоматически создать таблицу transactions, если её ещё нет
	log.Println("Running PostgreSQL database migrations...")
	err = db.AutoMigrate(&models.Transaction{}, &models.Transfer{}, &models.AuthorApplication{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
