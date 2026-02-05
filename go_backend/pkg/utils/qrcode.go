package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// GenerateQRCode generates a QR code string for bin verification
// In production, this would generate an actual QR code image
func GenerateQRCode(binID uuid.UUID, collectionID uuid.UUID) string {
	// Generate random bytes for uniqueness
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	// Create QR code content
	return fmt.Sprintf("SMARTWASTE:%s:%s:%s", binID.String(), collectionID.String(), randomHex)
}

// ValidateQRCode validates a QR code string
func ValidateQRCode(qrCode string, expectedBinID, expectedCollectionID uuid.UUID) bool {
	// Parse QR code
	var prefix, binIDStr, collectionIDStr, _ string
	_, err := fmt.Sscanf(qrCode, "%[^:]:%36s:%36s:%s", &prefix, &binIDStr, &collectionIDStr)
	if err != nil {
		return false
	}

	if prefix != "SMARTWASTE" {
		return false
	}

	// Validate IDs
	parsedBinID, err := uuid.Parse(binIDStr)
	if err != nil || parsedBinID != expectedBinID {
		return false
	}

	parsedCollectionID, err := uuid.Parse(collectionIDStr)
	if err != nil || parsedCollectionID != expectedCollectionID {
		return false
	}

	return true
}

// ExtractQRCodeData extracts bin and collection IDs from a QR code
func ExtractQRCodeData(qrCode string) (binID, collectionID uuid.UUID, err error) {
	var prefix string
	var binIDStr, collectionIDStr, randomStr string

	n, parseErr := fmt.Sscanf(qrCode, "SMARTWASTE:%36[^:]:%36[^:]:%s", &binIDStr, &collectionIDStr, &randomStr)
	if parseErr != nil || n < 2 {
		// Try alternative parsing
		_, parseErr = fmt.Sscanf(qrCode, "%[^:]:%36[^:]:%36[^:]:%s", &prefix, &binIDStr, &collectionIDStr, &randomStr)
		if parseErr != nil {
			return uuid.Nil, uuid.Nil, fmt.Errorf("invalid QR code format")
		}
	}

	binID, err = uuid.Parse(binIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid bin ID in QR code: %w", err)
	}

	collectionID, err = uuid.Parse(collectionIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid collection ID in QR code: %w", err)
	}

	return binID, collectionID, nil
}
