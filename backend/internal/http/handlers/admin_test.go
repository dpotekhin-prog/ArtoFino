package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAdminRepo struct{}

func (m *mockAdminRepo) Create(ctx context.Context, obj *mongo.Object) error {
	return nil // Mock successful write operation
}

func TestAdminHandler_Stats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{})
	r.GET("/admin/stats", h.Stats)

	req, _ := http.NewRequest(http.MethodGet, "/admin/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{})
	r.GET("/admin/ping", h.Ping)

	req, _ := http.NewRequest(http.MethodGet, "/admin/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApproveArtistApplication_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewAdminHandler(&mockAdminRepo{})
	r.POST("/admin/applications/:id/approve", h.ApproveArtistApplication)

	targetAppID := "app-test-uuid-999"
	req, _ := http.NewRequest(http.MethodPost, "/admin/applications/"+targetAppID+"/approve", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
