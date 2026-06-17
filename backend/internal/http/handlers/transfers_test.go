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

// Mock repositories for logistics
type mockMongoTransferRepo struct{}

func (m *mockMongoTransferRepo) FindByID(ctx context.Context, idStr string) (*mongo.Object, error) {
	hexID, _ := primitive.ObjectIDFromHex("64a7b3e1f1d2c3a4b5777777")
	return &mongo.Object{
		ID:          hexID,
		Title:       "Logistics Test Asset",
		OwnerUserID: "00000000-0000-0000-0000-000000000002",
		CreatedAt:   time.Now().AddDate(0, 0, -10),
	}, nil
}

type mockPostgresTransferRepo struct {
	mockBalanceSharePct float64
}

func (m *mockPostgresTransferRepo) Create(ctx context.Context, transfer *models.Transfer) error {
	transfer.ID = "generated-transfer-uuid-111"
	transfer.CreatedAt = time.Now()
	return nil
}

func (m *mockPostgresTransferRepo) FindByID(ctx context.Context, id string) (*models.Transfer, error) {
	return &models.Transfer{
		ID:          id,
		ObjectID:    "64a7b3e1f1d2c3a4b5777777",
		RequesterID: "00000000-0000-0000-0000-000000000001",
		Destination: "Prague Gallery",
		Status:      "pending",
	}, nil
}

func (m *mockPostgresTransferRepo) Update(ctx context.Context, transfer *models.Transfer) error {
	return nil
}

func (m *mockPostgresTransferRepo) GetUserBalance(ctx context.Context, userID string, objectID string) (*models.ArtShareBalance, error) {
	if m.mockBalanceSharePct <= 0 {
		return nil, nil
	}
	return &models.ArtShareBalance{
		UserID:   userID,
		ObjectID: objectID,
		SharePct: m.mockBalanceSharePct,
	}, nil
}

func TestRequestTransfer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockMongo := &mockMongoTransferRepo{}
	mockPg := &mockPostgresTransferRepo{mockBalanceSharePct: 10.0} // User owns 10% of the artwork
	h := NewTransfersHandler(mockMongo, mockPg)
	r.POST("/transfers/request", h.RequestTransfer)

	inputJSON := `{"objectId": "64a7b3e1f1d2c3a4b5777777", "destination": "Prague Gallery"}`

	req, _ := http.NewRequest(http.MethodPost, "/transfers/request", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response TransferRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "64a7b3e1f1d2c3a4b5777777", response.ObjectID)
	assert.Equal(t, "Prague Gallery", response.Destination)
	assert.Equal(t, "pending", response.Status)
	assert.NotEmpty(t, response.TransferID)
}

func TestRequestTransfer_Forbidden_NoShares(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	mockMongo := &mockMongoTransferRepo{}
	mockPg := &mockPostgresTransferRepo{mockBalanceSharePct: 0.0} // User has no shares
	h := NewTransfersHandler(mockMongo, mockPg)
	r.POST("/transfers/request", h.RequestTransfer)

	inputJSON := `{"objectId": "64a7b3e1f1d2c3a4b5777777", "destination": "Prague Gallery"}`

	req, _ := http.NewRequest(http.MethodPost, "/transfers/request", strings.NewReader(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
