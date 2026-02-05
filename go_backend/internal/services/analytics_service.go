package services

import (
	"context"
	"time"

	"github.com/smartwaste/backend/internal/repository"
)

// AnalyticsService handles analytics and reporting
type AnalyticsService struct {
	binRepo        *repository.BinRepository
	collectionRepo *repository.CollectionRepository
	driverRepo     *repository.DriverRepository
}

// NewAnalyticsService creates a new AnalyticsService
func NewAnalyticsService(
	binRepo *repository.BinRepository,
	collectionRepo *repository.CollectionRepository,
	driverRepo *repository.DriverRepository,
) *AnalyticsService {
	return &AnalyticsService{
		binRepo:        binRepo,
		collectionRepo: collectionRepo,
		driverRepo:     driverRepo,
	}
}

// DashboardStats represents overall dashboard statistics
type DashboardStats struct {
	TotalBins           int                    `json:"total_bins"`
	BinsNeedingCollection int                  `json:"bins_needing_collection"`
	AverageFillLevel    float64                `json:"average_fill_level"`
	TodayCollections    int                    `json:"today_collections"`
	TodayWeightKg       float64                `json:"today_weight_kg"`
	MonthCollections    int                    `json:"month_collections"`
	ActiveDrivers       int                    `json:"active_drivers"`
	Timestamp           time.Time              `json:"timestamp"`
	BinStats            map[string]interface{} `json:"bin_stats,omitempty"`
	CollectionStats     map[string]interface{} `json:"collection_stats,omitempty"`
}

// GetDashboardStats retrieves comprehensive dashboard statistics
func (s *AnalyticsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{
		Timestamp: time.Now(),
	}

	// Get bin statistics
	binStats, err := s.binRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}
	stats.BinStats = binStats
	stats.TotalBins = binStats["total_bins"].(int)
	stats.BinsNeedingCollection = binStats["needs_collection"].(int)
	stats.AverageFillLevel = binStats["average_fill_level"].(float64)

	// Get collection statistics
	collectionStats, err := s.collectionRepo.GetCollectionStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.CollectionStats = collectionStats
	stats.TodayCollections = collectionStats["today_collections"].(int)
	stats.TodayWeightKg = collectionStats["today_weight_kg"].(float64)
	stats.MonthCollections = collectionStats["month_collections"].(int)

	// Get active drivers count
	drivers, err := s.driverRepo.GetAvailableDrivers(ctx)
	if err != nil {
		return nil, err
	}
	stats.ActiveDrivers = len(drivers)

	return stats, nil
}

// BinAnalytics represents bin-specific analytics
type BinAnalytics struct {
	TotalBins         int              `json:"total_bins"`
	ActiveBins        int              `json:"active_bins"`
	AverageFillLevel  float64          `json:"average_fill_level"`
	BinsByFillRange   []FillRangeCount `json:"bins_by_fill_range"`
	BinsNeedingAction int              `json:"bins_needing_action"`
}

// FillRangeCount represents count of bins in a fill level range
type FillRangeCount struct {
	Range string `json:"range"`
	Count int    `json:"count"`
}

// GetBinAnalytics retrieves bin-specific analytics
func (s *AnalyticsService) GetBinAnalytics(ctx context.Context) (*BinAnalytics, error) {
	stats, err := s.binRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return &BinAnalytics{
		TotalBins:         stats["total_bins"].(int),
		ActiveBins:        stats["total_bins"].(int), // Same for now
		AverageFillLevel:  stats["average_fill_level"].(float64),
		BinsNeedingAction: stats["needs_collection"].(int),
		BinsByFillRange: []FillRangeCount{
			{Range: "0-25%", Count: 0},   // Would query from DB
			{Range: "26-50%", Count: 0},  // Would query from DB
			{Range: "51-75%", Count: 0},  // Would query from DB
			{Range: "76-100%", Count: stats["needs_collection"].(int)},
		},
	}, nil
}

// DriverPerformance represents driver performance metrics
type DriverPerformance struct {
	TotalDrivers      int     `json:"total_drivers"`
	AvailableDrivers  int     `json:"available_drivers"`
	AverageRating     float64 `json:"average_rating"`
	TotalCollections  int     `json:"total_collections"`
	AveragePerDriver  float64 `json:"average_per_driver"`
}

// GetDriverAnalytics retrieves driver-specific analytics
func (s *AnalyticsService) GetDriverAnalytics(ctx context.Context) (*DriverPerformance, error) {
	drivers, err := s.driverRepo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	availableDrivers, err := s.driverRepo.GetAvailableDrivers(ctx)
	if err != nil {
		return nil, err
	}

	totalCollections := 0
	totalRating := 0.0
	for _, d := range drivers {
		totalCollections += d.TotalCollections
		totalRating += d.AverageRating
	}

	avgRating := 0.0
	avgPerDriver := 0.0
	if len(drivers) > 0 {
		avgRating = totalRating / float64(len(drivers))
		avgPerDriver = float64(totalCollections) / float64(len(drivers))
	}

	return &DriverPerformance{
		TotalDrivers:     len(drivers),
		AvailableDrivers: len(availableDrivers),
		AverageRating:    avgRating,
		TotalCollections: totalCollections,
		AveragePerDriver: avgPerDriver,
	}, nil
}

// CollectionAnalytics represents collection analytics
type CollectionAnalytics struct {
	TodayCollections    int     `json:"today_collections"`
	WeekCollections     int     `json:"week_collections"`
	MonthCollections    int     `json:"month_collections"`
	TotalWeightToday    float64 `json:"total_weight_today_kg"`
	TotalWeightMonth    float64 `json:"total_weight_month_kg"`
	AverageCollectionTime string `json:"average_collection_time"`
}

// GetCollectionAnalytics retrieves collection analytics
func (s *AnalyticsService) GetCollectionAnalytics(ctx context.Context) (*CollectionAnalytics, error) {
	stats, err := s.collectionRepo.GetCollectionStats(ctx)
	if err != nil {
		return nil, err
	}

	return &CollectionAnalytics{
		TodayCollections: stats["today_collections"].(int),
		MonthCollections: stats["month_collections"].(int),
		TotalWeightToday: stats["today_weight_kg"].(float64),
	}, nil
}
