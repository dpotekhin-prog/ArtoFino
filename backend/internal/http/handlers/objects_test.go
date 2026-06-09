package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetArtObject_PriceCalculation(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize router and handler
	r := gin.Default()
	h := NewObjectsHandler()
	r.GET("/objects/:id", h.GetArtObject)

	// Create a test request for a mock object ID
	req, _ := http.NewRequest(http.MethodGet, "/objects/test-uuid-123", nil)
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal JSON response
	var response ArtObjectResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Calculate expected price based on the same mock values used in handler
	// BasePrice: 5000.00, DailyGrowthRate: 0.00025, Age: 60 days
	expectedBasePrice := 5000.00
	expectedRate := 0.00025
	daysPassed := time.Since(response.CreatedAt).Hours() / 24
	expectedCurrentPrice := expectedBasePrice * (1 + expectedRate*daysPassed)

	// Assert correct mathematical behavior
	assert.Equal(t, "test-uuid-123", response.ID)
	assert.Equal(t, "EUR", response.Currency)
	assert.Equal(t, expectedBasePrice, response.BasePrice)
	assert.Equal(t, expectedRate, response.DailyGrowthRate)
	assert.InEpsilon(t, expectedCurrentPrice, response.CurrentPrice, 0.0001, "The dynamic current price calculation is incorrect")
}
