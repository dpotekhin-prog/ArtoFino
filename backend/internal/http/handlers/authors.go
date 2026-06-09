package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreatorApplicationInput describes the payload for an artist profile upgrade
type CreatorApplicationInput struct {
	Bio          string   `json:"bio" binding:"required,max=500" example:"Contemporary oil painter based in Prague, focusing on urban landscapes."`
	PortfolioURL string   `json:"portfolioUrl" binding:"required,url" example:"https://behance.net/my-art-portfolio"`
	SocialLinks  []string `json:"socialLinks" binding:"required,gt=0" example:"['https://instagram.com/artist_handle']"`
}

// CreatorApplicationResponse describes the state of the submitted application
type CreatorApplicationResponse struct {
	ApplicationID string    `json:"applicationId" example:"app-111222"`
	UserID        string    `json:"userId" example:"user-keycloak-uuid"`
	Bio           string    `json:"bio"`
	PortfolioURL  string    `json:"portfolioUrl"`
	Status        string    `json:"status" example:"pending"` // pending, approved, rejected
	CreatedAt     time.Time `json:"createdAt"`
}

// AuthorsHandler holds dependencies for artist profile and verification logic
type AuthorsHandler struct{}

// NewAuthorsHandler creates a new instance of AuthorsHandler
func NewAuthorsHandler() *AuthorsHandler {
	return &AuthorsHandler{}
}

// ApplyForCreator godoc
// @Summary      Apply for Artist (Creator) status
// @Description  Allows a registered user to submit portfolio and social links to get vetted and upgraded to the 'artist' role.
// @Tags         authors
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreatorApplicationInput true "Artist application data"
// @Success      202 {object} CreatorApplicationResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /authors/apply [post]
func (h *AuthorsHandler) ApplyForCreator(c *gin.Context) {
	var input CreatorApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract User ID from Keycloak JWT claims
	claims, exists := c.Get("claims")
	var userID string
	if !exists {
		userID = "mock-user-applicant-id"
	} else {
		userID = claims.(map[string]interface{})["sub"].(string)
	}

	// TODO: Save to Postgres database verification infrastructure:
	// INSERT INTO artist_applications (user_id, bio, portfolio_url, status) VALUES ...

	c.JSON(http.StatusAccepted, CreatorApplicationResponse{
		ApplicationID: "app-mock-uuid-555666",
		UserID:        userID,
		Bio:           input.Bio,
		PortfolioURL:  input.PortfolioURL,
		Status:        "pending", // Requires admin approval via /admin/applications/:id/approve
		CreatedAt:     time.Now(),
	})
}
