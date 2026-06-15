package handlers

import (
	"context"
	"net/http"
	"time"

	"ArtoFino/backend/internal/models"
	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
)

// Contracts for safe decoupling
type MongoTxRepository interface {
	FindByID(ctx context.Context, idStr string) (*mongo.Object, error)
}

type PostgresTxRepository interface {
	Save(ctx context.Context, tx *models.Transaction) error
}

type TransactionsHandler struct {
	mongoRepo MongoTxRepository
	pgRepo    PostgresTxRepository
}

func NewTransactionsHandler(mongoRepo MongoTxRepository, pgRepo PostgresTxRepository) *TransactionsHandler {
	return &TransactionsHandler{
		mongoRepo: mongoRepo,
		pgRepo:    pgRepo,
	}
}

type BuyShareInput struct {
	ObjectID string  `json:"objectId" binding:"required" example:"64a7b3e1f1d2c3a4b5777777"`
	SharePct float64 `json:"sharePct" binding:"required,gte=1,lte=10" example:"5"`
}

type TransactionResponse struct {
	TransactionID string    `json:"transactionId" example:"uuid-string-here"`
	ObjectID      string    `json:"objectId" example:"64a7b3e1f1d2c3a4b5777777"`
	BuyerID       string    `json:"buyerId" example:"user-buyer-uuid"`
	SellerID      string    `json:"sellerId" example:"artist-seller-uuid"`
	SharePct      float64   `json:"sharePct" example:"5"`
	AmountPaid    float64   `json:"amountPaid" example:"253.75"`
	Currency      string    `json:"currency" example:"EUR"`
	CreatedAt     time.Time `json:"createdAt"`
}

// BuyShare godoc
// @Summary      Buy a share in an art object
// @Description  Purchase a fractional percentage of an art object. Price is locked and calculated dynamically using the live growth formula from MongoDB.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body BuyShareInput true "Purchase percentage details"
// @Success      201 {object} TransactionResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /transactions/buy [post]
func (h *TransactionsHandler) BuyShare(c *gin.Context) {
	var input BuyShareInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("claims")
	var buyerID string
	if !exists {
		// Mock valid UUID string format for constraints if claims don't exist
		buyerID = "00000000-0000-0000-0000-000000000001"
	} else {
		buyerID = claims.(map[string]interface{})["sub"].(string)
	}

	// 1. Fetch asset metadata from MongoDB
	obj, err := h.mongoRepo.FindByID(c.Request.Context(), input.ObjectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "target art object not found"})
		return
	}

	// 2. LIVE ECONOMIC FORMULA EVALUATION
	daysPassed := time.Since(obj.CreatedAt).Hours() / 24
	if daysPassed < 0 {
		daysPassed = 0
	}
	currentPrice := obj.BasePrice * (1 + obj.DailyGrowthRate*daysPassed)

	// 3. Calculate dynamic cost and convert to Cents (int64)
	amountToPay := currentPrice * (input.SharePct / 100)
	amountCents := int64(amountToPay * 100)

	// Fallback mock for seller if owner_user_id is not a valid UUID in test state
	sellerID := obj.OwnerUserID
	if sellerID == "" {
		sellerID = "00000000-0000-0000-0000-000000000002"
	}

	// 4. Build and save database transaction record
	dbTx := &models.Transaction{
		ObjectID:    input.ObjectID,
		FromUserID:  sellerID, // From whom (Seller/Artist)
		ToUserID:    buyerID,  // To whom (Buyer)
		AmountCents: amountCents,
		SharePct:    input.SharePct,
		Currency:    obj.Currency,
	}

	if err := h.pgRepo.Save(c.Request.Context(), dbTx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "financial ledger recording failed"})
		return
	}

	c.JSON(http.StatusCreated, TransactionResponse{
		TransactionID: dbTx.ID, // GORM will populate this via gen_random_uuid()
		ObjectID:      dbTx.ObjectID,
		BuyerID:       dbTx.ToUserID,
		SellerID:      dbTx.FromUserID,
		SharePct:      dbTx.SharePct,
		AmountPaid:    float64(dbTx.AmountCents) / 100, // Return as normal standard float to API
		Currency:      dbTx.Currency,
		CreatedAt:     dbTx.CreatedAt,
	})
}
