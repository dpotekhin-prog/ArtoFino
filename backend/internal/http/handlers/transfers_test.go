package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestTransfer_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize router and handler
	r := gin.Default()
	h := NewTransfersHandler()
	r.POST("/transfers/request", h.RequestTransfer)

	// Prepare payload for booking an art object for a local apartment exhibition
	input := TransferRequestInput{
		ObjectID:     "64a7b3e1f1d2c3a4b5",
		EventDetails: "Apartment art party in Prague",
		DurationDays: 5,
	}
	jsonPayload, _ := json.Marshal(input)

	// Create POST request with JSON payload
	req, _ := http.NewRequest(http.MethodPost, "/transfers/request", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert HTTP status code is 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Unmarshal JSON response
	var response TransferRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert correct operational state mapping
	assert.Equal(t, "tr-mock-uuid-888999", response.TransferID)
	assert.Equal(t, input.ObjectID, response.ObjectID)
	assert.Equal(t, "mock-host-partner-id", response.RequesterID)
	assert.Equal(t, "partner-9876", response.FromHolderID)
	assert.Equal(t, "pending", response.Status)
	assert.Equal(t, input.EventDetails, response.EventDetails)
	assert.NotEmpty(t, response.ExpiresAt)
	assert.NotEmpty(t, response.CreatedAt)
}

func TestApproveTransfer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	h := NewTransfersHandler()
	r.POST("/transfers/:id/approve", h.ApproveTransfer)

	targetTransferID := "tr-test-123"

	req, _ := http.NewRequest(http.MethodPost, "/transfers/"+targetTransferID+"/approve", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response TransferApprovalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, targetTransferID, response.TransferID)
	assert.Equal(t, "approved", response.Status)
	assert.NotEmpty(t, response.ApprovedAt)
}
