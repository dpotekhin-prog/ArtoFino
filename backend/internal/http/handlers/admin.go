package handlers

import (
	"context"
	"net/http"
	"time"

	"ArtoFino/backend/internal/models"
	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
)

// MongoAdminRepository defines the contract for inserting new assets
type MongoAdminRepository interface {
	Create(ctx context.Context, obj *mongo.Object) error
}

// AdminAuthorApplicationsRepo defines the contract for GORM Author applications management
type AdminAuthorApplicationsRepo interface {
	FindByID(ctx context.Context, id string) (*models.AuthorApplication, error)
	Update(ctx context.Context, app *models.AuthorApplication) error
}

type AdminHandler struct {
	mongoRepo MongoAdminRepository
	appRepo   AdminAuthorApplicationsRepo // Подключили Postgres-репозиторий заявок
}

// NewAdminHandler initializes AdminHandler according to Variant A (Explicit DI)
func NewAdminHandler(mongoRepo MongoAdminRepository, appRepo AdminAuthorApplicationsRepo) *AdminHandler {
	return &AdminHandler{
		mongoRepo: mongoRepo,
		appRepo:   appRepo,
	}
}

type CreateObjectInput struct {
	Title           string   `json:"title" binding:"required" example:"Symphony of Lights"`
	Description     string   `json:"description" example:"An incredible modern oil canvas exploration."`
	Tags            []string `json:"tags" example:"['oil', 'canvas', 'modern']"`
	BasePrice       float64  `json:"basePrice" binding:"required,gt=0" example:"5000.00"`
	Currency        string   `json:"currency" binding:"required" example:"EUR"`
	DailyGrowthRate float64  `json:"dailyGrowthRate" binding:"required" example:"0.00025"`
	ArtistID        string   `json:"artistId" binding:"required" example:"artist-uuid-111"`
}

type AdminStatsResponse struct {
	TotalUsers   int     `json:"totalUsers" example:"154"`
	ActiveOrders int     `json:"activeOrders" example:"12"`
	TotalRevenue float64 `json:"totalRevenue" example:"1450.50"`
}

func (h *AdminHandler) Stats(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "admin ok",
		"info":   "some admin stats here",
	})
}

type AdminPingResponse struct {
	Message string `json:"message" example:"admin pong"`
}

func (h *AdminHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "admin pong",
	})
}

type AdminApplicationApprovalResponse struct {
	ApplicationID string    `json:"applicationId" example:"app-555666"`
	UserID        string    `json:"userId" example:"user-keycloak-uuid"`
	Status        string    `json:"status" example:"approved"`
	RoleAssigned  string    `json:"roleAssigned" example:"artist"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ApproveArtistApplication godoc
// @Summary      Approve an author onboarding request
// @Description  Allows an Admin to review and approve a pending author application inside PostgreSQL.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Application ID"
// @Success      200 {object} AdminApplicationApprovalResponse
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/applications/{id}/approve [post]
func (h *AdminHandler) ApproveArtistApplication(c *gin.Context) {
	appID := c.Param("id")

	// 1. Извлекаем реальную заявку из PostgreSQL
	app, err := h.appRepo.FindByID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "author application not found"})
		return
	}

	// 2. Мутируем статус на approved
	app.Status = "approved"
	app.UpdatedAt = time.Now()

	// 3. Сохраняем изменения в реляционную БД
	if err := h.appRepo.Update(c.Request.Context(), app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update application status"})
		return
	}

	// Отдаем красивый оригинальный JSON-ответ фронтенду
	c.JSON(http.StatusOK, AdminApplicationApprovalResponse{
		ApplicationID: app.ID,
		UserID:        app.UserID,
		Status:        app.Status,
		RoleAssigned:  "artist",
		UpdatedAt:     app.UpdatedAt,
	})
}

// CreateArtObject godoc
// @Summary      Create a new art object asset
// @Description  Allows administrators to mint a new digital art object metadata entry inside MongoDB with base pricing configuration.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreateObjectInput true "Art Object configuration state"
// @Success      201 {object} mongo.Object
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/objects [post]
func (h *AdminHandler) CreateArtObject(c *gin.Context) {
	var input CreateObjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newObj := &mongo.Object{
		Title:           input.Title,
		Description:     input.Description,
		Tags:            input.Tags,
		BasePrice:       input.BasePrice,
		Currency:        input.Currency,
		DailyGrowthRate: input.DailyGrowthRate,
		OwnerUserID:     input.ArtistID,
		CurrentHolderID: input.ArtistID,
	}

	if err := h.mongoRepo.Create(c.Request.Context(), newObj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist object inside database"})
		return
	}

	c.JSON(http.StatusCreated, newObj)
}
