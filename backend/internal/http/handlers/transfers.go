package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TransferRequestInput describes the payload to request an art object for temporary use
type TransferRequestInput struct {
	ObjectID     string `json:"objectId" binding:"required" example:"64a7b3e1f1d2c3a4b5"`
	EventDetails string `json:"eventDetails" binding:"required" example:"Corporate pop-up exhibition or apartment art party"`
	DurationDays int    `json:"durationDays" binding:"required,gte=1" example:"7"`
}

// TransferRequestResponse describes the created logistics transfer state
type TransferRequestResponse struct {
	TransferID   string    `json:"transferId" example:"tr-999111"`
	ObjectID     string    `json:"objectId" example:"64a7b3e1f1d2c3a4b5"`
	RequesterID  string    `json:"requesterId" example:"partner-host-uuid"`
	FromHolderID string    `json:"fromHolderId" example:"partner-current-owner-uuid"`
	Status       string    `json:"status" example:"pending"` // pending, approved, rejected, completed
	EventDetails string    `json:"eventDetails" example:"Corporate pop-up exhibition or apartment art party"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// TransfersHandler holds dependencies for art object physical movements
type TransfersHandler struct{}

// NewTransfersHandler creates a new instance of TransfersHandler
func NewTransfersHandler() *TransfersHandler {
	return &TransfersHandler{}
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
// @Failure      401 {object} map[string]string
// @Router       /transfers/request [post]
func (h *TransfersHandler) RequestTransfer(c *gin.Context) {
	var input TransferRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Requester ID from Keycloak JWT tokens
	claims, exists := c.Get("claims")
	var requesterID string
	if !exists {
		requesterID = "mock-host-partner-id"
	} else {
		requesterID = claims.(map[string]interface{})["sub"].(string)
	}

	// 1. Simulating fetching the current asset state from MongoDB
	// We need to know who currently physically holds the painting (FromHolderID)
	mockCurrentHolderID := "partner-9876" // Currently sitting at this partner's location

	// 2. Logic: Ensure the requester is not trying to request from themselves
	if requesterID == mockCurrentHolderID && exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already hold this art object"})
		return
	}

	// TODO: Save to Postgres log infrastructure:
	// INSERT INTO art_object_transfers (object_id, requester_id, from_holder_id, status, event_details) ...

	c.JSON(http.StatusCreated, TransferRequestResponse{
		TransferID:   "tr-mock-uuid-888999",
		ObjectID:     input.ObjectID,
		RequesterID:  requesterID,
		FromHolderID: mockCurrentHolderID,
		Status:       "pending",
		EventDetails: input.EventDetails,
		ExpiresAt:    time.Now().AddDate(0, 0, 3), // Request open for 3 days
		CreatedAt:    time.Now(),
	})
}

// TransferApprovalResponse describes the result of a transfer request being approved by the holder
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
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /transfers/{id}/approve [post]
func (h *TransfersHandler) ApproveTransfer(c *gin.Context) {
	transferID := c.Param("id")

	// TODO: Verify via Postgres that the current authenticated user matches 'FromHolderID'
	// TODO: Update transfer record status to 'approved'
	// TODO: Trigger dynamic dailyGrowthRate recalculation inside MongoDB for this object

	c.JSON(http.StatusOK, TransferApprovalResponse{
		TransferID: transferID,
		Status:     "approved",
		ApprovedAt: time.Now(),
	})
}
