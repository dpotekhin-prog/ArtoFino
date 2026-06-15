package db

import (
	"context"

	"ArtoFino/backend/internal/models"

	"gorm.io/gorm"
)

type TransactionsRepository struct {
	db *gorm.DB
}

// NewTransactionsRepository initializes Postgres repository for financial transactions
func NewTransactionsRepository(db *gorm.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db}
}

// Save inserts a new transaction record into PostgreSQL
func (r *TransactionsRepository) Save(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}
