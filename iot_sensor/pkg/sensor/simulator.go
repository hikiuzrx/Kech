package sensor

import (
	"math/rand"
)

// Simulator implements Sensor interface with fake data
type Simulator struct {
	currentLevel float64
	binHeight    float64
}

// NewSimulator creates a new simulated sensor
func NewSimulator(binHeight float64) *Simulator {
	return &Simulator{
		currentLevel: binHeight * 0.2, // Start 20% full (distance is large)
		binHeight:    binHeight,
	}
}

// ReadDistance simulates reading distance
// Returns distance from top (Sensor) to waste
func (s *Simulator) ReadDistance() (float64, error) {
	// Simulate filling up
	// Decrease distance (waste level rising)
	change := (rand.Float64() * 5) - 1 // Mostly filling (positive fill = negative distance change)
	s.currentLevel -= change

	if s.currentLevel < 10 {
		s.currentLevel = s.binHeight // Reset if full (emptied)
	}
	if s.currentLevel > s.binHeight {
		s.currentLevel = s.binHeight
	}

	return s.currentLevel, nil
}

func (s *Simulator) Close() error {
	return nil
}
