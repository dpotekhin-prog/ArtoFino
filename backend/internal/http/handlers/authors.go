package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"ArtoFino/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AuthorApplicationsRepo interface {
	Create(ctx context.Context, app *models.AuthorApplication) error
}

type AuthorsHandler struct {
	repo AuthorApplicationsRepo
}

func NewAuthorsHandler(repo AuthorApplicationsRepo) *AuthorsHandler {
	return &AuthorsHandler{repo: repo}
}

type CreatorApplicationInput struct {
	Bio          string   `json:"bio" binding:"required,max=500" example:"Contemporary oil painter based in Prague, focusing on urban landscapes."`
	PortfolioURL string   `json:"portfolioUrl" binding:"required,url" example:"https://behance.net/my-art-portfolio"`
	SocialLinks  []string `json:"socialLinks" binding:"required,gt=0" example:"[\"https://instagram.com/artist_handle\"]"`
}

type CreatorApplicationResponse struct {
	ApplicationID string    `json:"applicationId" example:"uuid-string"`
	UserID        string    `json:"userId" example:"user-keycloak-uuid"`
	Bio           string    `json:"bio"`
	PortfolioURL  string    `json:"portfolioUrl"`
	Status        string    `json:"status" example:"pending"`
	CreatedAt     time.Time `json:"createdAt"`
}

// ApplyForCreator godoc
// @Summary      Apply for Artist (Creator) status
// @Description  Allows a registered user to submit portfolio and social links to get vetted and upgraded to the 'artist' role.
// @Tags         authors
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreatorApplicationInput true "Artist application data"
// @Success      201 {object} CreatorApplicationResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /authors/apply [post]
func (h *AuthorsHandler) ApplyForCreator(c *gin.Context) {
	var input CreatorApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("claims")
	var userID string
	if !exists {
		userID = "00000000-0000-0000-0000-000000000001"
	} else {
		userID = claims.(map[string]interface{})["sub"].(string)
	}

	// Склеиваем слайс ссылок в одну строку для БД
	linksStr := strings.Join(input.SocialLinks, ",")

	app := &models.AuthorApplication{
		UserID:       userID,
		Bio:          input.Bio,
		PortfolioURL: input.PortfolioURL,
		SocialLinks:  linksStr,
		Status:       "pending",
	}

	if err := h.repo.Create(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save author application"})
		return
	}

	c.JSON(http.StatusCreated, CreatorApplicationResponse{
		ApplicationID: app.ID,
		UserID:        app.UserID,
		Bio:           app.Bio,
		PortfolioURL:  app.PortfolioURL,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt,
	})
}
