package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
	"github.com/smartwaste/backend/internal/services"
	"github.com/smartwaste/backend/pkg/utils"
)

// DriverHandler handles driver-related HTTP requests
type DriverHandler struct {
	driverRepo     *repository.DriverRepository
	binRepo        *repository.BinRepository
	collectionRepo *repository.CollectionRepository
	routeService   *services.RouteService
}

// NewDriverHandler creates a new DriverHandler
func NewDriverHandler(
	driverRepo *repository.DriverRepository,
	binRepo *repository.BinRepository,
	collectionRepo *repository.CollectionRepository,
	routeService *services.RouteService,
) *DriverHandler {
	return &DriverHandler{
		driverRepo:     driverRepo,
		binRepo:        binRepo,
		collectionRepo: collectionRepo,
		routeService:   routeService,
	}
}

// GetDriver retrieves a driver by ID
// @Summary Get driver by ID
// @Tags Drivers
// @Produce json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.DriverResponse
// @Failure 404 {object} utils.APIError
// @Router /api/v1/drivers/{id} [get]
func (h *DriverHandler) GetDriver(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	driver, err := h.driverRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve driver")
		return
	}

	if driver == nil {
		utils.NotFound(c, "Driver not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, driver.ToResponse())
}

// CreateDriver creates a new driver
// @Summary Create a new driver
// @Tags Drivers
// @Accept json
// @Produce json
// @Param driver body models.CreateDriverRequest true "Driver data"
// @Success 201 {object} models.DriverResponse
// @Failure 400 {object} utils.APIError
// @Router /api/v1/drivers [post]
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req models.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// Check if email already exists
	existing, err := h.driverRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		utils.InternalError(c, "Failed to check existing driver")
		return
	}
	if existing != nil {
		utils.Conflict(c, "Email already registered")
		return
	}

	driver := &models.Driver{
		Email:         req.Email,
		PasswordHash:  req.Password, // In production, hash this!
		FullName:      req.FullName,
		Phone:         req.Phone,
		LicenseNumber: req.LicenseNumber,
		VehicleType:   req.VehicleType,
		VehiclePlate:  req.VehiclePlate,
		IsAvailable:   true,
	}

	if err := h.driverRepo.Create(c.Request.Context(), driver); err != nil {
		utils.InternalError(c, "Failed to create driver")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, driver.ToResponse())
}

// UpdateDriver updates a driver
// @Summary Update driver
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Param driver body models.UpdateDriverRequest true "Driver data"
// @Success 200 {object} models.DriverResponse
// @Router /api/v1/drivers/{id} [put]
func (h *DriverHandler) UpdateDriver(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	var req models.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	driver, err := h.driverRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve driver")
		return
	}
	if driver == nil {
		utils.NotFound(c, "Driver not found")
		return
	}

	// Update fields
	if req.FullName != nil {
		driver.FullName = *req.FullName
	}
	if req.Phone != nil {
		driver.Phone = *req.Phone
	}
	if req.VehicleType != nil {
		driver.VehicleType = req.VehicleType
	}
	if req.VehiclePlate != nil {
		driver.VehiclePlate = req.VehiclePlate
	}
	if req.IsAvailable != nil {
		driver.IsAvailable = *req.IsAvailable
	}

	if err := h.driverRepo.Update(c.Request.Context(), driver); err != nil {
		utils.InternalError(c, "Failed to update driver")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, driver.ToResponse())
}

// UpdateLocation updates a driver's location
// @Summary Update driver location
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Param location body models.UpdateDriverLocationRequest true "Location data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/drivers/{id}/location [put]
func (h *DriverHandler) UpdateLocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	var req models.UpdateDriverLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := h.driverRepo.UpdateLocation(c.Request.Context(), id, req.Latitude, req.Longitude); err != nil {
		utils.InternalError(c, "Failed to update location")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"driver_id": id,
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
		"message":   "Location updated successfully",
	})
}

// GetRoutes retrieves optimized routes for a driver
// @Summary Get optimized routes
// @Tags Drivers
// @Produce json
// @Param id path string true "Driver ID"
// @Param optimize_by query string false "Optimization criteria: distance or fill_level" default(distance)
// @Success 200 {object} models.RouteResponse
// @Router /api/v1/drivers/{id}/routes [get]
func (h *DriverHandler) GetRoutes(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	driver, err := h.driverRepo.GetByID(c.Request.Context(), id)
	if err != nil || driver == nil {
		utils.NotFound(c, "Driver not found")
		return
	}

	// Get bins needing collection (>80% full)
	bins, err := h.routeService.GetBinsForRoute(c.Request.Context(), 80)
	if err != nil {
		utils.InternalError(c, "Failed to get bins for route")
		return
	}

	if len(bins) == 0 {
		utils.SuccessResponse(c, http.StatusOK, gin.H{
			"message":   "No bins need collection",
			"waypoints": []models.Waypoint{},
		})
		return
	}

	// Extract bin IDs
	binIDs := make([]uuid.UUID, len(bins))
	for i, b := range bins {
		binIDs[i] = b.ID
	}

	// Get driver location (use default if not set)
	driverLat := 0.0
	driverLng := 0.0
	if driver.Latitude != nil && driver.Longitude != nil {
		driverLat = *driver.Latitude
		driverLng = *driver.Longitude
	}

	optimizeBy := c.DefaultQuery("optimize_by", "distance")
	route, err := h.routeService.OptimizeRoute(c.Request.Context(), driverLat, driverLng, binIDs, optimizeBy)
	if err != nil {
		utils.InternalError(c, "Failed to calculate route")
		return
	}

	route.DriverID = id
	utils.SuccessResponse(c, http.StatusOK, route.ToResponse())
}

// VerifyTask verifies a collection task via QR code
// @Summary Verify task via QR code
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Param request body models.VerifyTaskRequest true "QR code data"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/drivers/{id}/verify [post]
func (h *DriverHandler) VerifyTask(c *gin.Context) {
	idParam := c.Param("id")
	driverID, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	var req models.VerifyTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	collectionID, err := uuid.Parse(req.CollectionID)
	if err != nil {
		utils.BadRequest(c, "Invalid collection ID format")
		return
	}

	// Get collection
	collection, err := h.collectionRepo.GetByID(c.Request.Context(), collectionID)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve collection")
		return
	}
	if collection == nil {
		utils.NotFound(c, "Collection not found")
		return
	}

	// Verify driver is assigned to this collection
	if collection.DriverID != driverID {
		utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "You are not assigned to this collection")
		return
	}

	// Extract and verify QR code
	binID, qrCollectionID, err := utils.ExtractQRCodeData(req.QRCode)
	if err != nil {
		utils.BadRequest(c, "Invalid QR code format")
		return
	}

	if qrCollectionID != collectionID || binID != collection.BinID {
		utils.BadRequest(c, "QR code does not match collection")
		return
	}

	// Mark as verified
	if err := h.collectionRepo.VerifyQRCode(c.Request.Context(), collectionID); err != nil {
		utils.InternalError(c, "Failed to verify collection")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"collection_id": collectionID,
		"verified":      true,
		"message":       "Task verified successfully",
	})
}

// GetDriverStats retrieves driver performance statistics
// @Summary Get driver statistics
// @Tags Drivers
// @Produce json
// @Param id path string true "Driver ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/drivers/{id}/stats [get]
func (h *DriverHandler) GetDriverStats(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid driver ID format")
		return
	}

	stats, err := h.collectionRepo.GetDriverStats(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve driver statistics")
		return
	}

	stats["driver_id"] = id
	utils.SuccessResponse(c, http.StatusOK, stats)
}

// ListDrivers retrieves all drivers with pagination
// @Summary List drivers
// @Tags Drivers
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {array} models.DriverResponse
// @Router /api/v1/drivers [get]
func (h *DriverHandler) ListDrivers(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 20)
	offset := (page - 1) * perPage

	drivers, err := h.driverRepo.List(c.Request.Context(), perPage, offset)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve drivers")
		return
	}

	responses := make([]models.DriverResponse, len(drivers))
	for i, d := range drivers {
		responses[i] = *d.ToResponse()
	}

	utils.SuccessResponseWithPagination(c, responses, &utils.Pagination{
		Page:    page,
		PerPage: perPage,
	})
}
