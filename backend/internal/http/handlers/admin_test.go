package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_Stats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	h := NewAdminHandler()
	r.GET("/admin/stats", h.Stats)

	req, _ := http.NewRequest(http.MethodGet, "/admin/stats", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin ok")
}

func TestAdminHandler_Ping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	h := NewAdminHandler()
	r.GET("/admin/ping", h.Ping)

	req, _ := http.NewRequest(http.MethodGet, "/admin/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin pong")
}

func TestApproveArtistApplication_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	h := NewAdminHandler()
	r.POST("/admin/applications/:id/approve", h.ApproveArtistApplication)

	targetAppID := "app-test-uuid-999"

	req, _ := http.NewRequest(http.MethodPost, "/admin/applications/"+targetAppID+"/approve", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AdminApplicationApprovalResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, targetAppID, response.ApplicationID)
	assert.Equal(t, "mock-user-applicant-id", response.UserID)
	assert.Equal(t, "approved", response.Status)
	assert.Equal(t, "artist", response.RoleAssigned)
	assert.NotEmpty(t, response.UpdatedAt)
}
