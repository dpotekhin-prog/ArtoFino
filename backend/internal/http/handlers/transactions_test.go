package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ArtoFino/backend/internal/models"
	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repositories for isolation testing
type mockMongoTxRepo struct{}

func (m *mockMongoTxRepo) FindByID(ctx context.Context, idStr string) (*mongo.Object, error) {
	hexID, _ := primitive.ObjectIDFromHex("64a7b3e1f1d2c3a4b5777777")
	return &mongo.Object{
		ID:              hexID,
		Title:           "Test Canvas",
		BasePrice:       10000.00,
		Currency:        "EUR",
		DailyGrowthRate: 0.0002,
		OwnerUserID:     "00000000-0000-0000-0000-000000000002",
		CreatedAt:       time.Now().AddDate(0, 0, -5), // 5 days ago
	}, nil
}

type mockPostgresTxRepo struct{}

func (m *mockPostgresTxRepo) Save(ctx context.Context, tx *models.Transaction) error {
	tx.ID = "generated-test-uuid-000999" // Simulate DB default hook
	return nil
}

func TestBuyShare_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockMongo := &mockMongoTxRepo{}
	mockPg := &mockPostgresTxRepo{}
	h := NewTransactionsHandler(mockMongo, mockPg)
	r.POST("/transactions/buy", h.BuyShare)

	inputJSON := `{"objectId": "64a7b3e1f1d2c3a4b5777777", "sharePct": 5}`

	req, _ := http.NewRequest(http.MethodPost, "/transactions/buy", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response TransactionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "64a7b3e1f1d2c3a4b5777777", response.ObjectID)
	assert.Equal(t, 5.0, response.SharePct)
	assert.Equal(t, "EUR", response.Currency)
	assert.NotEmpty(t, response.TransactionID)
	assert.True(t, response.AmountPaid > 0)
}
