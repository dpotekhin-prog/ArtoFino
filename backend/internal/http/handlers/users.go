package handlers

import (
	"github.com/gin-gonic/gin"
)

type UsersHandler struct{}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}

// Me godoc
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /users/me [get]
func (h *UsersHandler) Me(c *gin.Context) {
	claims := c.MustGet("claims").(map[string]interface{})

	c.JSON(200, gin.H{
		"id":        claims["sub"],
		"email":     claims["email"],
		"username":  claims["preferred_username"],
		"firstName": claims["given_name"],
		"lastName":  claims["family_name"],
	})
}
