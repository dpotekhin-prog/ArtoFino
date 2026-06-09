package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

type AdminStatsResponse struct {
	TotalUsers   int     `json:"totalUsers" example:"154"`
	ActiveOrders int     `json:"activeOrders" example:"12"`
	TotalRevenue float64 `json:"totalRevenue" example:"1450.50"`
}

// Stats godoc
// @Summary      Admin stats
// @Tags         admin
// @Security     BearerAuth
// @Success      200 {object} AdminStatsResponse
// @Router       /admin/stats [get]
func (h *AdminHandler) Stats(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "admin ok",
		"info":   "some admin stats here",
	})
}

type AdminPingResponse struct {
	Message string `json:"message" example:"admin pong"`
}

// Ping godoc
// @Summary      Admin ping
// @Tags         admin
// @Security     BearerAuth
// @Success      200 {object} AdminPingResponse
// @Router       /admin/ping [get]
func (h *AdminHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "admin pong",
	})
}

// AdminApplicationApprovalResponse describes the result of an artist application approval
type AdminApplicationApprovalResponse struct {
	ApplicationID string    `json:"applicationId" example:"app-555666"`
	UserID        string    `json:"userId" example:"user-keycloak-uuid"`
	Status        string    `json:"status" example:"approved"`
	RoleAssigned  string    `json:"roleAssigned" example:"artist"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ApproveArtistApplication godoc
// @Summary      Approve an artist application
// @Description  Allows an administrator to review and approve a pending creator application, upgrading the user's system role.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Application ID"
// @Success      200 {object} AdminApplicationApprovalResponse
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Router       /admin/applications/{id}/approve [post]
func (h *AdminHandler) ApproveArtistApplication(c *gin.Context) {
	appID := c.Param("id")

	// TODO: Fetch application from Postgres, update status to 'approved'
	// TODO: Call Keycloak API to assign the 'artist' client role to the target UserID

	c.JSON(http.StatusOK, AdminApplicationApprovalResponse{
		ApplicationID: appID,
		UserID:        "mock-user-applicant-id", // In reality, taken from the DB record
		Status:        "approved",
		RoleAssigned:  "artist",
		UpdatedAt:     time.Now(),
	})
}
