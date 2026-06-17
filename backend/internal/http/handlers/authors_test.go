package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ArtoFino/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthorAppRepo struct{}

func (m *mockAuthorAppRepo) Create(ctx context.Context, app *models.AuthorApplication) error {
	app.ID = "mock-app-uuid-111"
	app.CreatedAt = time.Now() // Simulate database default timestamp behavior
	return nil
}

func TestApplyForCreator_Success(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize router and handler with mock dependency
	r := gin.Default()
	mockRepo := &mockAuthorAppRepo{}
	h := NewAuthorsHandler(mockRepo)
	r.POST("/authors/apply", h.ApplyForCreator)

	// Prepare valid creator application data
	input := CreatorApplicationInput{
		Bio:          "Prague-based independent artist specializing in digital renaissance concepts.",
		PortfolioURL: "https://www.behance.net/prague-digital-art",
		SocialLinks:  []string{"https://instagram.com/prague_digital_art"},
	}
	jsonPayload, _ := json.Marshal(input)

	// Create POST request with JSON payload
	req, _ := http.NewRequest(http.MethodPost, "/authors/apply", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	r.ServeHTTP(w, req)

	// Assert HTTP status code is 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Unmarshal JSON response
	var response CreatorApplicationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert application metadata and state mapping
	assert.Equal(t, "mock-app-uuid-111", response.ApplicationID)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", response.UserID)
	assert.Equal(t, input.Bio, response.Bio)
	assert.Equal(t, input.PortfolioURL, response.PortfolioURL)
	assert.Equal(t, "pending", response.Status)
	assert.NotEmpty(t, response.CreatedAt)
}
