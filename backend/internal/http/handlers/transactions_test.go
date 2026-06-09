package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBuyShare_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize router and handler
	r := gin.Default()
	h := NewTransactionsHandler()
	r.POST("/transactions/buy", h.BuyShare)

	// Prepare payload for purchasing 5% of the art object
	input := BuyShareInput{
		ObjectID: "64a7b3e1f1d2c3a4b5",
		SharePct: 5.0,
	}
	jsonPayload, _ := json.Marshal(input)

	// Create POST request with JSON body
	req, _ := http.NewRequest(http.MethodPost, "/transactions/buy", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert HTTP status code is 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Unmarshal JSON response
	var response TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Mathematical evaluation verification:
	// BasePrice: 12000.00, DailyGrowthRate: 0.0001, Age: 100 days
	// Current price equals 12000 * (1 + 0.0001 * 100) = 12120.00 CZK
	// 5% share of 12120.00 equals exactly 606.00 CZK
	expectedBasePrice := 12000.00
	expectedRate := 0.0001
	daysPassed := time.Since(time.Now().AddDate(0, 0, -100)).Hours() / 24
	expectedCurrentPrice := expectedBasePrice * (1 + expectedRate*daysPassed)
	expectedAmountPaid := expectedCurrentPrice * (input.SharePct / 100)

	// Assert final dynamic numbers match expectations perfectly
	assert.Equal(t, "tx-live-calculated-111", response.TransactionID)
	assert.Equal(t, input.ObjectID, response.ObjectID)
	assert.Equal(t, "artist-specific-uuid", response.SellerID)
	assert.Equal(t, input.SharePct, response.SharePct)
	assert.Equal(t, "CZK", response.Currency)
	assert.InEpsilon(t, expectedAmountPaid, response.AmountPaid, 0.0001, "The share cost calculation does not match the live logic asset rate")
}
