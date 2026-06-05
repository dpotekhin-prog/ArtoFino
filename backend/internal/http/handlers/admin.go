package handlers

import "github.com/gin-gonic/gin"

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// Stats godoc
// @Summary Admin stats
// @Tags admin
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/stats [get]
func (h *AdminHandler) Stats(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "admin ok",
		"info":   "some admin stats here",
	})
}

// Ping godoc
// @Summary Admin ping
// @Tags admin
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /admin/ping [get]
func (h *AdminHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "admin pong",
	})
}
