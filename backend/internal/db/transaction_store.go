package db

import (
	"errors"
	"fmt"
	"time"

	"ArtoFino/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{db: db}
}

// ExecutePurchase atomically transfers share percentages and logs the transaction ledger
func (s *TransactionStore) ExecutePurchase(
	objectID string,
	fromUserID string,
	toUserID string,
	sharePct float64,
	amountCents int64,
	currency string,
) (*models.Transaction, error) {

	if sharePct <= 0 {
		return nil, errors.New("share percentage must be greater than zero")
	}
	if fromUserID == toUserID {
		return nil, errors.New("buyer cannot be the seller")
	}

	var txRecord models.Transaction

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var sellerBalance models.ArtShareBalance

		// 1. Lock seller balance row to prevent race conditions (SELECT ... FOR UPDATE)
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("object_id = ? AND user_id = ?", objectID, fromUserID).
			First(&sellerBalance).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("seller has no active shares for object %s", objectID)
			}
			return err
		}

		if sellerBalance.SharePct < sharePct {
			return fmt.Errorf("insufficient shares balance: seller has %v, requested %v", sellerBalance.SharePct, sharePct)
		}

		// 2. Deduct shares from the seller
		sellerBalance.SharePct -= sharePct
		sellerBalance.UpdatedAt = time.Now()
		if err := tx.Save(&sellerBalance).Error; err != nil {
			return err
		}

		// 3. Credit shares to the buyer (Upsert balance record)
		var buyerBalance models.ArtShareBalance
		err = tx.Where("object_id = ? AND user_id = ?", objectID, toUserID).
			First(&buyerBalance).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				buyerBalance = models.ArtShareBalance{
					ObjectID:  objectID,
					UserID:    toUserID,
					SharePct:  sharePct,
					UpdatedAt: time.Now(),
				}
				if err := tx.Create(&buyerBalance).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			buyerBalance.SharePct += sharePct
			buyerBalance.UpdatedAt = time.Now()
			if err := tx.Save(&buyerBalance).Error; err != nil {
				return err
			}
		}

		// 4. Generate transaction history record
		txRecord = models.Transaction{
			ObjectID:    objectID,
			FromUserID:  fromUserID,
			ToUserID:    toUserID,
			AmountCents: amountCents,
			SharePct:    sharePct,
			Currency:    currency,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &txRecord, nil
}
