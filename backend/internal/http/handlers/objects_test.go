package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockObjectsRepo struct{}

func (m *mockObjectsRepo) FindByID(ctx context.Context, idStr string) (*mongo.Object, error) {
	hexID, _ := primitive.ObjectIDFromHex("64a7b3e1f1d2c3a4b5777777")
	return &mongo.Object{
		ID:              hexID,
		Title:           "Test Masterpiece",
		BasePrice:       5000.00,
		Currency:        "EUR",
		DailyGrowthRate: 0.00025,
		CreatedAt:       time.Now().AddDate(0, 0, -10), // 10 days ago
	}, nil
}

func TestGetArtObject_PriceCalculation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockRepo := &mockObjectsRepo{}
	h := NewObjectsHandler(mockRepo)
	r.GET("/objects/:id", h.GetArtObject)

	req, _ := http.NewRequest(http.MethodGet, "/objects/64a7b3e1f1d2c3a4b5777777", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ArtObjectResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedBasePrice := 5000.00
	expectedRate := 0.00025
	daysPassed := time.Since(response.CreatedAt).Hours() / 24
	expectedCurrentPrice := expectedBasePrice * (1 + expectedRate*daysPassed)

	assert.Equal(t, "64a7b3e1f1d2c3a4b5777777", response.ID)
	assert.Equal(t, "EUR", response.Currency)
	assert.Equal(t, expectedBasePrice, response.BasePrice)
	assert.Equal(t, expectedRate, response.DailyGrowthRate)
	assert.InEpsilon(t, expectedCurrentPrice, response.CurrentPrice, 0.0001)
}
