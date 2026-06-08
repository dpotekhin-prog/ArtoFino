package handlers

import (
	"github.com/gin-gonic/gin"
)

type UsersHandler struct{}

func NewUsersHandler() *UsersHandler {
	return &UsersHandler{}
}

// UserProfileResponse представляет структуру ответа профиля пользователя для Swagger документации
type UserProfileResponse struct {
	ID        string `json:"id" example:"4321-abcd-1234"`
	Email     string `json:"email" example:"user@artofino.com"`
	Username  string `json:"username" example:"johndoe"`
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName" example:"Doe"`
}

// Me godoc
// @Summary      Get current user profile
// @Tags         users
// @Security     BearerAuth
// @Success      200  {object}  UserProfileResponse
// @Router       /users/me [get]
func (h *UsersHandler) Me(c *gin.Context) {
	// ВОТ ЭТОЙ СТРОКИ СЕЙЧАС НЕТ ИЛИ ОНА СЛОМАЛАСЬ:
	claims := c.MustGet("claims").(map[string]interface{})

	c.JSON(200, gin.H{
		"id":        claims["sub"],
		"email":     claims["email"],
		"username":  claims["preferred_username"],
		"firstName": claims["given_name"],
		"lastName":  claims["family_name"],
	})
}
