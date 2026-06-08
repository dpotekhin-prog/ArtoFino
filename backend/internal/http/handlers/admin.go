package handlers

import "github.com/gin-gonic/gin"

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
