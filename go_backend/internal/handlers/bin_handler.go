package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
	"github.com/smartwaste/backend/pkg/utils"
)

// BinHandler handles bin-related HTTP requests
type BinHandler struct {
	repo *repository.BinRepository
}

// NewBinHandler creates a new BinHandler
func NewBinHandler(repo *repository.BinRepository) *BinHandler {
	return &BinHandler{repo: repo}
}

// GetBin retrieves a bin by ID
// @Summary Get bin by ID
// @Tags Bins
// @Produce json
// @Param id path string true "Bin ID"
// @Success 200 {object} models.BinResponse
// @Failure 404 {object} utils.APIError
// @Router /api/v1/bins/{id} [get]
func (h *BinHandler) GetBin(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid bin ID format")
		return
	}

	bin, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bin")
		return
	}

	if bin == nil {
		utils.NotFound(c, "Bin not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, bin.ToResponse())
}

// CreateBin creates a new bin
// @Summary Register a new bin
// @Tags Bins
// @Accept json
// @Produce json
// @Param bin body models.CreateBinRequest true "Bin data"
// @Success 201 {object} models.BinResponse
// @Failure 400 {object} utils.APIError
// @Router /api/v1/bins [post]
func (h *BinHandler) CreateBin(c *gin.Context) {
	var req models.CreateBinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// Check if device ID already exists
	existing, err := h.repo.GetByDeviceID(c.Request.Context(), req.DeviceID)
	if err != nil {
		utils.InternalError(c, "Failed to check existing bin")
		return
	}
	if existing != nil {
		utils.Conflict(c, "Device ID already registered")
		return
	}

	bin := &models.Bin{
		DeviceID:       req.DeviceID,
		LocationName:   req.LocationName,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		WasteType:      req.WasteType,
		CapacityLiters: req.CapacityLiters,
		CompanyID:      req.CompanyID,
		IsActive:       true,
	}

	if err := h.repo.Create(c.Request.Context(), bin); err != nil {
		utils.InternalError(c, "Failed to create bin")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, bin.ToResponse())
}

// UpdateBin updates a bin
// @Summary Update bin
// @Tags Bins
// @Accept json
// @Produce json
// @Param id path string true "Bin ID"
// @Param bin body models.UpdateBinRequest true "Bin data"
// @Success 200 {object} models.BinResponse
// @Router /api/v1/bins/{id} [put]
func (h *BinHandler) UpdateBin(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid bin ID format")
		return
	}

	var req models.UpdateBinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	bin, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bin")
		return
	}
	if bin == nil {
		utils.NotFound(c, "Bin not found")
		return
	}

	// Update fields
	if req.LocationName != nil {
		bin.LocationName = req.LocationName
	}
	if req.Latitude != nil {
		bin.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		bin.Longitude = *req.Longitude
	}
	if req.WasteType != nil {
		bin.WasteType = *req.WasteType
	}
	if req.CapacityLiters != nil {
		bin.CapacityLiters = *req.CapacityLiters
	}
	if req.IsActive != nil {
		bin.IsActive = *req.IsActive
	}
	if req.CompanyID != nil {
		bin.CompanyID = req.CompanyID
	}

	if err := h.repo.Update(c.Request.Context(), bin); err != nil {
		utils.InternalError(c, "Failed to update bin")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, bin.ToResponse())
}

// ListBins retrieves all bins with pagination
// @Summary List bins
// @Tags Bins
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {array} models.BinResponse
// @Router /api/v1/bins [get]
func (h *BinHandler) ListBins(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 20)
	offset := (page - 1) * perPage

	bins, err := h.repo.List(c.Request.Context(), perPage, offset)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bins")
		return
	}

	responses := make([]models.BinResponse, len(bins))
	for i, b := range bins {
		responses[i] = *b.ToResponse()
	}

	utils.SuccessResponseWithPagination(c, responses, &utils.Pagination{
		Page:    page,
		PerPage: perPage,
	})
}

// GetBinsNeedingCollection retrieves bins with fill level above threshold
// @Summary Get bins needing collection
// @Tags Bins
// @Produce json
// @Param threshold query int false "Fill level threshold" default(80)
// @Success 200 {array} models.BinResponse
// @Router /api/v1/bins/needs-collection [get]
func (h *BinHandler) GetBinsNeedingCollection(c *gin.Context) {
	threshold := getQueryInt(c, "threshold", 80)

	bins, err := h.repo.GetBinsNeedingCollection(c.Request.Context(), threshold)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bins")
		return
	}

	responses := make([]models.BinResponse, len(bins))
	for i, b := range bins {
		responses[i] = *b.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"threshold": threshold,
		"count":     len(bins),
		"bins":      responses,
	})
}

// GetBinStatistics retrieves bin statistics
// @Summary Get bin statistics
// @Tags Bins
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/bins/statistics [get]
func (h *BinHandler) GetBinStatistics(c *gin.Context) {
	stats, err := h.repo.GetStatistics(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bin statistics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stats)
}

// DeleteBin deletes a bin (soft delete)
// @Summary Delete bin
// @Tags Bins
// @Param id path string true "Bin ID"
// @Success 204 "No Content"
// @Failure 404 {object} utils.APIError
// @Router /api/v1/bins/{id} [delete]
func (h *BinHandler) DeleteBin(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid bin ID format")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		utils.InternalError(c, "Failed to delete bin")
		return
	}

	c.Status(http.StatusNoContent)
}
