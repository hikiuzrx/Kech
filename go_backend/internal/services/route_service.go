package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"

	"github.com/google/uuid"
	"github.com/smartwaste/backend/internal/config"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
)

// RouteService handles route optimization for drivers
type RouteService struct {
	binRepo   *repository.BinRepository
	googleKey string
}

// NewRouteService creates a new RouteService
func NewRouteService(binRepo *repository.BinRepository, cfg *config.GoogleConfig) *RouteService {
	return &RouteService{
		binRepo:   binRepo,
		googleKey: cfg.MapsAPIKey,
	}
}

// OptimizeRoute calculates an optimized route for a driver
func (s *RouteService) OptimizeRoute(ctx context.Context, driverLat, driverLng float64, binIDs []uuid.UUID, optimizeBy string) (*models.DriverRoute, error) {
	// Get bins
	bins := make([]*models.Bin, 0, len(binIDs))
	for _, id := range binIDs {
		bin, err := s.binRepo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get bin %s: %w", id, err)
		}
		if bin != nil {
			bins = append(bins, bin)
		}
	}

	if len(bins) == 0 {
		return nil, fmt.Errorf("no valid bins found")
	}

	// Sort bins based on optimization criteria
	var waypoints []models.Waypoint
	switch optimizeBy {
	case "fill_level":
		waypoints = s.optimizeByFillLevel(bins, driverLat, driverLng)
	case "distance":
		fallthrough
	default:
		waypoints = s.optimizeByDistance(bins, driverLat, driverLng)
	}

	// Calculate total distance and duration
	totalDistance, duration := s.calculateRouteMetrics(driverLat, driverLng, waypoints)

	route := &models.DriverRoute{
		ID:                       uuid.New(),
		WaypointsList:            waypoints,
		TotalDistanceKm:          &totalDistance,
		EstimatedDurationMinutes: &duration,
		Status:                   models.RouteStatusPending,
	}

	// Marshal waypoints to JSON for storage
	waypointsJSON, err := json.Marshal(waypoints)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal waypoints: %w", err)
	}
	route.Waypoints = waypointsJSON

	// Try to get optimized route from Google Maps/OSRM
	if s.googleKey != "" {
		optimizedRoute, err := s.getGoogleMapsRoute(driverLat, driverLng, waypoints)
		if err != nil {
			log.Printf("Failed to get Google Maps route, using calculated distance: %v", err)
		} else if optimizedRoute != nil {
			route.TotalDistanceKm = &optimizedRoute.distance
			route.EstimatedDurationMinutes = &optimizedRoute.duration
		}
	}

	return route, nil
}

// optimizeByDistance sorts bins by distance from driver (nearest first)
func (s *RouteService) optimizeByDistance(bins []*models.Bin, driverLat, driverLng float64) []models.Waypoint {
	type binWithDistance struct {
		bin      *models.Bin
		distance float64
	}

	// Calculate distances
	binsWithDist := make([]binWithDistance, len(bins))
	for i, bin := range bins {
		binsWithDist[i] = binWithDistance{
			bin:      bin,
			distance: haversineDistance(driverLat, driverLng, bin.Latitude, bin.Longitude),
		}
	}

	// Use nearest neighbor algorithm
	waypoints := make([]models.Waypoint, 0, len(bins))
	currentLat, currentLng := driverLat, driverLng
	visited := make(map[uuid.UUID]bool)

	for len(waypoints) < len(bins) {
		var nearest *models.Bin
		minDist := math.MaxFloat64

		for _, bwd := range binsWithDist {
			if visited[bwd.bin.ID] {
				continue
			}
			dist := haversineDistance(currentLat, currentLng, bwd.bin.Latitude, bwd.bin.Longitude)
			if dist < minDist {
				minDist = dist
				nearest = bwd.bin
			}
		}

		if nearest != nil {
			visited[nearest.ID] = true
			waypoints = append(waypoints, models.Waypoint{
				BinID:     nearest.ID,
				DeviceID:  nearest.DeviceID,
				Latitude:  nearest.Latitude,
				Longitude: nearest.Longitude,
				FillLevel: nearest.FillLevel,
				Order:     len(waypoints) + 1,
			})
			currentLat, currentLng = nearest.Latitude, nearest.Longitude
		}
	}

	return waypoints
}

// optimizeByFillLevel prioritizes bins with higher fill levels first, then by distance
func (s *RouteService) optimizeByFillLevel(bins []*models.Bin, driverLat, driverLng float64) []models.Waypoint {
	// Sort by fill level (descending), then by distance
	sort.Slice(bins, func(i, j int) bool {
		if bins[i].FillLevel != bins[j].FillLevel {
			return bins[i].FillLevel > bins[j].FillLevel
		}
		distI := haversineDistance(driverLat, driverLng, bins[i].Latitude, bins[i].Longitude)
		distJ := haversineDistance(driverLat, driverLng, bins[j].Latitude, bins[j].Longitude)
		return distI < distJ
	})

	waypoints := make([]models.Waypoint, len(bins))
	for i, bin := range bins {
		waypoints[i] = models.Waypoint{
			BinID:     bin.ID,
			DeviceID:  bin.DeviceID,
			Latitude:  bin.Latitude,
			Longitude: bin.Longitude,
			FillLevel: bin.FillLevel,
			Order:     i + 1,
		}
	}

	return waypoints
}

// calculateRouteMetrics calculates approximate distance and duration
func (s *RouteService) calculateRouteMetrics(startLat, startLng float64, waypoints []models.Waypoint) (float64, int) {
	if len(waypoints) == 0 {
		return 0, 0
	}

	totalDistance := 0.0
	currentLat, currentLng := startLat, startLng

	for _, wp := range waypoints {
		totalDistance += haversineDistance(currentLat, currentLng, wp.Latitude, wp.Longitude)
		currentLat, currentLng = wp.Latitude, wp.Longitude
	}

	// Estimate duration: assume average speed of 30 km/h in urban areas + 2 min per stop
	durationMinutes := int((totalDistance/30)*60) + len(waypoints)*2

	return totalDistance, durationMinutes
}

// haversineDistance calculates distance between two points using Haversine formula
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// googleMapsRouteResult holds the result from Google Maps API
type googleMapsRouteResult struct {
	distance float64 // km
	duration int     // minutes
}

// getGoogleMapsRoute fetches optimized route from Google Maps Directions API
func (s *RouteService) getGoogleMapsRoute(startLat, startLng float64, waypoints []models.Waypoint) (*googleMapsRouteResult, error) {
	if s.googleKey == "" {
		return nil, fmt.Errorf("google Maps API key not configured")
	}

	if len(waypoints) == 0 {
		return nil, fmt.Errorf("no waypoints provided")
	}

	// Build waypoints string
	waypointStrs := make([]string, len(waypoints))
	for i, wp := range waypoints {
		waypointStrs[i] = fmt.Sprintf("%f,%f", wp.Latitude, wp.Longitude)
	}

	// Last waypoint is destination
	destination := waypointStrs[len(waypointStrs)-1]
	intermediateWaypoints := ""
	if len(waypointStrs) > 1 {
		intermediateWaypoints = "optimize:true|" + url.QueryEscape(waypointStrs[0])
		for i := 1; i < len(waypointStrs)-1; i++ {
			intermediateWaypoints += "|" + waypointStrs[i]
		}
	}

	apiURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/directions/json?origin=%f,%f&destination=%s&waypoints=%s&key=%s",
		startLat, startLng, destination, intermediateWaypoints, s.googleKey,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Maps API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Status string `json:"status"`
		Routes []struct {
			Legs []struct {
				Distance struct {
					Value int `json:"value"` // meters
				} `json:"distance"`
				Duration struct {
					Value int `json:"value"` // seconds
				} `json:"duration"`
			} `json:"legs"`
		} `json:"routes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "OK" || len(result.Routes) == 0 {
		return nil, fmt.Errorf("no routes found: %s", result.Status)
	}

	// Sum up all legs
	totalDistance := 0
	totalDuration := 0
	for _, leg := range result.Routes[0].Legs {
		totalDistance += leg.Distance.Value
		totalDuration += leg.Duration.Value
	}

	return &googleMapsRouteResult{
		distance: float64(totalDistance) / 1000, // Convert to km
		duration: totalDuration / 60,             // Convert to minutes
	}, nil
}

// GetBinsForRoute retrieves bins that need collection
func (s *RouteService) GetBinsForRoute(ctx context.Context, threshold int) ([]models.Bin, error) {
	return s.binRepo.GetBinsNeedingCollection(ctx, threshold)
}
