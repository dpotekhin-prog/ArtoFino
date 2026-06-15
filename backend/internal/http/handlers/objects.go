package handlers

import (
	"context"
	"net/http"
	"time"

	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
)

// MongoObjectsRepository defines the contract required to fetch object parameters
type MongoObjectsRepository interface {
	FindByID(ctx context.Context, idStr string) (*mongo.Object, error)
}

type ObjectsHandler struct {
	repo MongoObjectsRepository
}

func NewObjectsHandler(repo MongoObjectsRepository) *ObjectsHandler {
	return &ObjectsHandler{repo: repo}
}

type ArtObjectResponse struct {
	ID              string    `json:"id" example:"64a7b3e1f1d2c3a4b5777777"`
	Title           string    `json:"title" example:"Symphony of Lights"`
	BasePrice       float64   `json:"basePrice" example:"5000.00"`
	CurrentPrice    float64   `json:"currentPrice" example:"5125.50"`
	Currency        string    `json:"currency" example:"EUR"`
	DailyGrowthRate float64   `json:"dailyGrowthRate" example:"0.00025"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (h *ObjectsHandler) GetArtObject(c *gin.Context) {
	idStr := c.Param("id")

	obj, err := h.repo.FindByID(c.Request.Context(), idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "art object not found or invalid id structure"})
		return
	}

	// Live economic age-based calculation engine
	daysPassed := time.Since(obj.CreatedAt).Hours() / 24
	if daysPassed < 0 {
		daysPassed = 0
	}
	currentPrice := obj.BasePrice * (1 + obj.DailyGrowthRate*daysPassed)

	c.JSON(http.StatusOK, ArtObjectResponse{
		ID:              obj.ID.Hex(),
		Title:           obj.Title,
		BasePrice:       obj.BasePrice,
		CurrentPrice:    currentPrice,
		Currency:        obj.Currency,
		DailyGrowthRate: obj.DailyGrowthRate,
		CreatedAt:       obj.CreatedAt,
	})
}
