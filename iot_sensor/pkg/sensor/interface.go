package sensor

// Sensor is the interface for reading distance/waste level
type Sensor interface {
	// ReadDistance returns the distance to the waste in centimeters
	ReadDistance() (float64, error)
	// Close releases any resources
	Close() error
}
