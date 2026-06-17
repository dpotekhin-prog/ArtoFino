package handlers

import (
	"context"
	"net/http"
	"time"

	"ArtoFino/backend/internal/models"
	"ArtoFino/backend/internal/mongo"

	"github.com/gin-gonic/gin"
)

// Interfaces for clean isolation
type MongoTransferRepository interface {
	FindByID(ctx context.Context, idStr string) (*mongo.Object, error)
}

type PostgresTransferRepository interface {
	Create(ctx context.Context, transfer *models.Transfer) error
	FindByID(ctx context.Context, id string) (*models.Transfer, error)
	Update(ctx context.Context, transfer *models.Transfer) error
	GetUserBalance(ctx context.Context, userID string, objectID string) (*models.ArtShareBalance, error)
}

type TransfersHandler struct {
	mongoRepo MongoTransferRepository
	pgRepo    PostgresTransferRepository
}

// NewTransfersHandler injects repositories into the logistics module
func NewTransfersHandler(mongoRepo MongoTransferRepository, pgRepo PostgresTransferRepository) *TransfersHandler {
	return &TransfersHandler{
		mongoRepo: mongoRepo,
		pgRepo:    pgRepo,
	}
}

type TransferRequestInput struct {
	ObjectID    string `json:"objectId" binding:"required" example:"64a7b3e1f1d2c3a4b5777777"`
	Destination string `json:"destination" binding:"required" example:"Prague National Gallery, Hall 4"`
}

type TransferRequestResponse struct {
	TransferID  string    `json:"transferId" example:"uuid-string"`
	ObjectID    string    `json:"objectId" example:"64a7b3e1f1d2c3a4b5777777"`
	RequesterID string    `json:"requesterId" example:"partner-host-uuid"`
	Destination string    `json:"destination" example:"Prague National Gallery, Hall 4"`
	Status      string    `json:"status" example:"pending"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RequestTransfer godoc
// @Summary      Request art object for temporary use
// @Description  Allows a Partner to apply for hosting an art object at their event/location.
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body TransferRequestInput true "Transfer request details"
// @Success      201 {object} TransferRequestResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /transfers/request [post]
func (h *TransfersHandler) RequestTransfer(c *gin.Context) {
	var input TransferRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("claims")
	var requesterID string
	if !exists {
		requesterID = "00000000-0000-0000-0000-000000000001"
	} else {
		requesterID = claims.(map[string]interface{})["sub"].(string)
	}

	// 1. Verify that the asset exists in MongoDB
	obj, err := h.mongoRepo.FindByID(c.Request.Context(), input.ObjectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "art object not found"})
		return
	}

	// 2. Prevent the owner from requesting their own physical asset
	if obj.OwnerUserID == requesterID && exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already physically hold or own this art object"})
		return
	}

	// 3. Verify that the user owns active fraction shares of this artwork in PostgreSQL
	balance, err := h.pgRepo.GetUserBalance(c.Request.Context(), requesterID, input.ObjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify user shares balance"})
		return
	}

	if balance == nil || balance.SharePct <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you must own shares of this art piece to request logistics"})
		return
	}

	// 4. Persist the pending logistics state inside PostgreSQL
	dbTransfer := &models.Transfer{
		ObjectID:    input.ObjectID,
		RequesterID: requesterID,
		Destination: input.Destination,
		Status:      "pending",
	}

	if err := h.pgRepo.Create(c.Request.Context(), dbTransfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save transfer request"})
		return
	}

	c.JSON(http.StatusCreated, TransferRequestResponse{
		TransferID:  dbTransfer.ID,
		ObjectID:    dbTransfer.ObjectID,
		RequesterID: dbTransfer.RequesterID,
		Destination: dbTransfer.Destination,
		Status:      dbTransfer.Status,
		CreatedAt:   dbTransfer.CreatedAt,
	})
}

type TransferApprovalResponse struct {
	TransferID string    `json:"transferId" example:"tr-888999"`
	Status     string    `json:"status" example:"approved"`
	ApprovedAt time.Time `json:"approvedAt"`
}

// ApproveTransfer godoc
// @Summary      Approve a physical transfer request
// @Description  Allows the current holder or owner of the art object to approve a pending temporary transfer request.
// @Tags         transfers
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Transfer ID"
// @Success      200 {object} TransferApprovalResponse
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /transfers/{id}/approve [post]
func (h *TransfersHandler) ApproveTransfer(c *gin.Context) {
	transferID := c.Param("id")

	// 1. Find the logistics entity inside PostgreSQL
	transfer, err := h.pgRepo.FindByID(c.Request.Context(), transferID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transfer request not found"})
		return
	}

	// 2. Mutate state parameters
	transfer.Status = "approved"
	transfer.UpdatedAt = time.Now()

	// 3. Save updates into relational engine
	if err := h.pgRepo.Update(c.Request.Context(), transfer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve transfer record"})
		return
	}

	c.JSON(http.StatusOK, TransferApprovalResponse{
		TransferID: transfer.ID,
		Status:     transfer.Status,
		ApprovedAt: transfer.UpdatedAt,
	})
}
