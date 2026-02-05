//go:build !linux && !darwin
// +build !linux,!darwin

// Note: This file is excluded from standard builds to avoid TinyGo dependency errors on standard OS
// It is intended for TinyGo builds targeting microcontrollers or Linux on ARM, handled via build tags if needed.
// However, since we are developing on Mac, we wrap this to avoid editor errors if tinygo isn't installed.

package sensor

import (
	"machine"

	"tinygo.org/x/drivers/hcsr04"
)

// HCSR04 implements Sensor interface using HC-SR04 ultrasonic sensor
type HCSR04 struct {
	device hcsr04.Device
}

// NewHCSR04 creates a new HCSR04 sensor
// Pins depend on the board. Example for RPi or Arduino.
func NewHCSR04(trigPin, echoPin machine.Pin) *HCSR04 {
	dev := hcsr04.New(trigPin, echoPin)
	dev.Configure()
	return &HCSR04{
		device: dev,
	}
}

func (s *HCSR04) ReadDistance() (float64, error) {
	// Read distance in mm, convert to cm
	distMm := s.device.ReadDistance()
	return float64(distMm) / 10.0, nil
}

func (s *HCSR04) Close() error {
	return nil
}
