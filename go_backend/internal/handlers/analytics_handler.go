package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwaste/backend/internal/services"
	"github.com/smartwaste/backend/pkg/utils"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsSvc *services.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(analyticsSvc *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsSvc: analyticsSvc}
}

// GetDashboardStats retrieves overall dashboard statistics
// @Summary Get dashboard statistics
// @Tags Analytics
// @Produce json
// @Success 200 {object} services.DashboardStats
// @Router /api/v1/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.analyticsSvc.GetDashboardStats(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to retrieve dashboard statistics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stats)
}

// GetBinAnalytics retrieves bin-specific analytics
// @Summary Get bin analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} services.BinAnalytics
// @Router /api/v1/analytics/bins [get]
func (h *AnalyticsHandler) GetBinAnalytics(c *gin.Context) {
	analytics, err := h.analyticsSvc.GetBinAnalytics(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to retrieve bin analytics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, analytics)
}

// GetDriverAnalytics retrieves driver performance analytics
// @Summary Get driver analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} services.DriverPerformance
// @Router /api/v1/analytics/drivers [get]
func (h *AnalyticsHandler) GetDriverAnalytics(c *gin.Context) {
	analytics, err := h.analyticsSvc.GetDriverAnalytics(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to retrieve driver analytics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, analytics)
}

// GetCollectionAnalytics retrieves collection analytics
// @Summary Get collection analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} services.CollectionAnalytics
// @Router /api/v1/analytics/collections [get]
func (h *AnalyticsHandler) GetCollectionAnalytics(c *gin.Context) {
	analytics, err := h.analyticsSvc.GetCollectionAnalytics(c.Request.Context())
	if err != nil {
		utils.InternalError(c, "Failed to retrieve collection analytics")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, analytics)
}

// Helper function to get query parameter as int
func getQueryInt(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
