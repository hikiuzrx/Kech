package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartwaste/shipment-tracker/internal/models"
	"github.com/smartwaste/shipment-tracker/internal/services"
)

// ShipmentHandler handles HTTP requests for shipments
type ShipmentHandler struct {
	service *services.ShipmentService
}

// NewShipmentHandler creates a new ShipmentHandler
func NewShipmentHandler(service *services.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}

// CreateShipment handles creating a new shipment
func (h *ShipmentHandler) CreateShipment(c *gin.Context) {
	var req models.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shipment, err := h.service.CreateShipment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shipment.ToResponse())
}

// GetShipment handles retrieving a shipment by ID
func (h *ShipmentHandler) GetShipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	shipment, err := h.service.GetShipment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if shipment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
		return
	}

	c.JSON(http.StatusOK, shipment.ToResponse())
}

// AssignDriver handles assigning a driver to a shipment
func (h *ShipmentHandler) AssignDriver(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req models.AssignDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignDriver(id, req.DriverID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Driver assigned successfully"})
}
