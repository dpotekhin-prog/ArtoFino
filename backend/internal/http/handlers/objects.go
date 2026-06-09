package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ArtObjectResponse describes the dynamic art object card from MongoDB
type ArtObjectResponse struct {
	ID              string    `json:"id" example:"64a7b3e1f1d2c3a4b5"`
	Title           string    `json:"title" example:"Sunset over Prague"`
	ArtistID        string    `json:"artistId" example:"artist-4321"`
	BasePrice       float64   `json:"basePrice" example:"5000.00"`
	Currency        string    `json:"currency" example:"EUR"`
	CurrentPrice    float64   `json:"currentPrice" example:"5075.50"`
	DailyGrowthRate float64   `json:"dailyGrowthRate" example:"0.0001"`
	CreatedAt       time.Time `json:"createdAt"`
	CurrentHolderID string    `json:"currentHolderId" example:"partner-9876"`
	Country         string    `json:"country" example:"CZ"`
}

// ObjectsHandler holds dependencies for art object endpoints
type ObjectsHandler struct{}

// NewObjectsHandler creates a new instance of ObjectsHandler
func NewObjectsHandler() *ObjectsHandler {
	return &ObjectsHandler{}
}

// GetArtObject godoc
// @Summary      Get art object with dynamic price
// @Description  Returns art object data with current price calculated based on daily linear growth, custom currency and dynamic rate
// @Tags         objects
// @Produce      json
// @Param        id   path      string  true  "Object ID"
// @Success      200 {object} ArtObjectResponse
// @Router       /objects/{id} [get]
func (h *ObjectsHandler) GetArtObject(c *gin.Context) {
	// Simulated database record (mocking MongoDB document)
	// This represents what the Artist actually saved + what the Platform calculated
	mockObjectFromDB := struct {
		ID              string
		Title           string
		ArtistID        string
		BasePrice       float64
		Currency        string
		DailyGrowthRate float64
		CreatedAt       time.Time
		CurrentHolderID string
		Country         string
	}{
		ID:              c.Param("id"),
		Title:           "Sunset over Prague",
		ArtistID:        "artist-4321",
		BasePrice:       5000.00,
		Currency:        "EUR",
		DailyGrowthRate: 0.00025,
		CreatedAt:       time.Now().AddDate(0, 0, -60), // 60 days ago
		CurrentHolderID: "partner-9876",
		Country:         "CZ",
	}

	// LINEAR GROWTH FORMULA USING DYNAMIC PARAMETERS
	daysPassed := time.Since(mockObjectFromDB.CreatedAt).Hours() / 24
	currentPrice := mockObjectFromDB.BasePrice * (1 + mockObjectFromDB.DailyGrowthRate*daysPassed)

	c.JSON(http.StatusOK, ArtObjectResponse{
		ID:              mockObjectFromDB.ID,
		Title:           mockObjectFromDB.Title,
		ArtistID:        mockObjectFromDB.ArtistID,
		BasePrice:       mockObjectFromDB.BasePrice,
		Currency:        mockObjectFromDB.Currency,
		CurrentPrice:    currentPrice,
		DailyGrowthRate: mockObjectFromDB.DailyGrowthRate,
		CreatedAt:       mockObjectFromDB.CreatedAt,
		CurrentHolderID: mockObjectFromDB.CurrentHolderID,
		Country:         mockObjectFromDB.Country,
	})
}
