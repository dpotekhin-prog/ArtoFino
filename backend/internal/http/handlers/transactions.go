package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BuyShareInput describes the request body for purchasing an art object share
type BuyShareInput struct {
	ObjectID string  `json:"objectId" binding:"required" example:"64a7b3e1f1d2c3a4b5"`
	SharePct float64 `json:"sharePct" binding:"required,gte=1,lte=10" example:"5"`
}

// TransactionResponse describes the result of a successful transaction recorded in Postgres
type TransactionResponse struct {
	TransactionID string    `json:"transactionId" example:"tx-abc-999888"`
	ObjectID      string    `json:"objectId" example:"64a7b3e1f1d2c3a4b5"`
	BuyerID       string    `json:"buyerId" example:"user-buyer-uuid"`
	SellerID      string    `json:"sellerId" example:"artist-seller-uuid"`
	SharePct      float64   `json:"sharePct" example:"5"`
	AmountPaid    float64   `json:"amountPaid" example:"253.75"`
	Currency      string    `json:"currency" example:"EUR"`
	CreatedAt     time.Time `json:"createdAt"`
}

// TransactionsHandler holds dependencies for financial operations
type TransactionsHandler struct{}

// NewTransactionsHandler creates a new instance of TransactionsHandler
func NewTransactionsHandler() *TransactionsHandler {
	return &TransactionsHandler{}
}

// BuyShare godoc
// @Summary      Buy a share in an art object
// @Description  Purchase between 1% and 10% of an art object. Price is locked at the current calculated rate using the linear growth formula.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body BuyShareInput true "Purchase data"
// @Success      201 {object} TransactionResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
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
		buyerID = "mock-buyer-id"
	} else {
		buyerID = claims.(map[string]interface{})["sub"].(string)
	}

	// 1. Simulating fetching the specific object from MongoDB to get its raw economic parameters
	// In the next step, this will be: object, err := h.mongoRepo.FindByID(input.ObjectID)
	mockObjectFromDB := struct {
		BasePrice       float64
		Currency        string
		DailyGrowthRate float64
		CreatedAt       time.Time
		ArtistID        string
	}{
		BasePrice:       12000.00,                       // For example, this specific artist set 12000
		Currency:        "CZK",                          // in Czech Korunas
		DailyGrowthRate: 0.0001,                         // with base 0.01% daily rate
		CreatedAt:       time.Now().AddDate(0, 0, -100), // published 100 days ago
		ArtistID:        "artist-specific-uuid",
	}

	// 2. THE EXACT SAME LIVE FORMULA
	// We calculate the object's real price at the exact second of the transaction
	daysPassed := time.Since(mockObjectFromDB.CreatedAt).Hours() / 24
	currentPrice := mockObjectFromDB.BasePrice * (1 + mockObjectFromDB.DailyGrowthRate*daysPassed)

	// 3. Calculate how much the buyer must pay for their specific percentage share
	amountToPay := currentPrice * (input.SharePct / 100)

	// TODO: Execute transactional database steps via Postgres ACID:
	// - Ensure total shares for this object do not exceed 100%
	// - Update ownership state, insert transaction record

	c.JSON(http.StatusCreated, TransactionResponse{
		TransactionID: "tx-live-calculated-111",
		ObjectID:      input.ObjectID,
		BuyerID:       buyerID,
		SellerID:      mockObjectFromDB.ArtistID,
		SharePct:      input.SharePct,
		AmountPaid:    amountToPay, // Locked and calculated dynamically!
		Currency:      mockObjectFromDB.Currency,
		CreatedAt:     time.Now(),
	})
}
