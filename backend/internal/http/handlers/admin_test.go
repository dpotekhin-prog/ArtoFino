package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ArtoFino/backend/internal/models"
	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Оригинальный мок для Mongo
type mockAdminRepo struct{}

func (m *mockAdminRepo) Create(ctx context.Context, obj *mongo.Object) error {
	return nil // Mock successful write operation
}

// Новый мок для Postgres-репозитория заявок авторов
type mockAdminAuthorAppRepo struct{}

func (m *mockAdminAuthorAppRepo) FindByID(ctx context.Context, id string) (*models.AuthorApplication, error) {
	return &models.AuthorApplication{
		ID:     id,
		UserID: "mock-user-applicant-id",
		Status: "pending",
	}, nil
}

func (m *mockAdminAuthorAppRepo) Update(ctx context.Context, app *models.AuthorApplication) error {
	return nil
}

func TestAdminHandler_Stats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{}, &mockAdminAuthorAppRepo{}) // Передали оба мока
	r.GET("/admin/stats", h.Stats)

	req, _ := http.NewRequest(http.MethodGet, "/admin/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{}, &mockAdminAuthorAppRepo{}) // Передали оба мока
	r.GET("/admin/ping", h.Ping)

	req, _ := http.NewRequest(http.MethodGet, "/admin/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApproveArtistApplication_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{}, &mockAdminAuthorAppRepo{}) // Передали оба мока
	r.POST("/admin/applications/:id/approve", h.ApproveArtistApplication)

	targetAppID := "app-test-uuid-999"
	req, _ := http.NewRequest(http.MethodPost, "/admin/applications/"+targetAppID+"/approve", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Дополнительно проверяем, что эндпоинт отдал корректно сматченный JSON
	var response AdminApplicationApprovalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, targetAppID, response.ApplicationID)
	assert.Equal(t, "mock-user-applicant-id", response.UserID)
	assert.Equal(t, "approved", response.Status)
}
