package db

import (
	"context"
	"errors"

	"ArtoFino/backend/internal/models"

	"gorm.io/gorm"
)

type TransfersRepository struct {
	db *gorm.DB
}

func NewTransfersRepository(db *gorm.DB) *TransfersRepository {
	return &TransfersRepository{db: db}
}

// Create inserts a new transfer request
func (r *TransfersRepository) Create(ctx context.Context, transfer *models.Transfer) error {
	return r.db.WithContext(ctx).Create(transfer).Error
}

// FindByID fetches a transfer request by its UUID
func (r *TransfersRepository) FindByID(ctx context.Context, id string) (*models.Transfer, error) {
	var transfer models.Transfer
	err := r.db.WithContext(ctx).First(&transfer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// Update saves modified transfer fields (like status)
func (r *TransfersRepository) Update(ctx context.Context, transfer *models.Transfer) error {
	return r.db.WithContext(ctx).Save(transfer).Error
}

// GetUserBalance retrieves the current share balance for a specific user and object
func (r *TransfersRepository) GetUserBalance(ctx context.Context, userID string, objectID string) (*models.ArtShareBalance, error) {
	var balance models.ArtShareBalance
	err := r.db.WithContext(ctx).Where("user_id = ? AND object_id = ?", userID, objectID).First(&balance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &balance, nil
}
