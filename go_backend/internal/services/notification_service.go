package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
)

// NotificationService handles notifications to drivers
type NotificationService struct {
	driverRepo       *repository.DriverRepository
	notificationRepo *NotificationRepository
}

// NotificationRepository handles notification data operations
type NotificationRepository struct {
	// This would be implemented similar to other repositories
	// For now, we'll log notifications as a placeholder
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(driverRepo *repository.DriverRepository) *NotificationService {
	return &NotificationService{
		driverRepo:       driverRepo,
		notificationRepo: &NotificationRepository{},
	}
}

// NotifyNearestDriver finds the nearest driver and sends them a notification
func (s *NotificationService) NotifyNearestDriver(ctx context.Context, bin *models.Bin) error {
	log.Printf("Finding nearest driver for bin %s at location (%.6f, %.6f)",
		bin.DeviceID, bin.Latitude, bin.Longitude)

	// Find nearest available driver
	driver, err := s.driverRepo.GetNearestDriver(ctx, bin.Latitude, bin.Longitude)
	if err != nil {
		return fmt.Errorf("failed to find nearest driver: %w", err)
	}

	if driver == nil {
		log.Printf("No available drivers found for bin %s", bin.DeviceID)
		return nil
	}

	// Create notification
	notification := &models.Notification{
		ID:       uuid.New(),
		DriverID: &driver.ID,
		BinID:    &bin.ID,
		Type:     models.NotificationTypeBinFull,
		Title:    "Bin Collection Required",
		Message: fmt.Sprintf(
			"Bin %s at %s is %d%% full and requires collection.",
			bin.DeviceID,
			*bin.LocationName,
			bin.FillLevel,
		),
	}

	// Send FCM notification (placeholder)
	if err := s.sendFCMNotification(driver, notification); err != nil {
		log.Printf("Failed to send FCM notification: %v", err)
		// Continue even if FCM fails - save notification for later retrieval
	}

	log.Printf("Notification sent to driver %s (%s) for bin %s",
		driver.ID, driver.FullName, bin.DeviceID)

	return nil
}

// sendFCMNotification sends a push notification via Firebase Cloud Messaging
// This is a placeholder implementation - in production, integrate with FCM SDK
func (s *NotificationService) sendFCMNotification(driver *models.Driver, notification *models.Notification) error {
	// Placeholder for FCM integration
	// In production:
	// 1. Use firebase.google.com/go/messaging
	// 2. Create message with driver.FCMToken
	// 3. Send via messaging.Client.Send()

	if driver.FCMToken == nil || *driver.FCMToken == "" {
		log.Printf("Driver %s has no FCM token, skipping push notification", driver.ID)
		return nil
	}

	log.Printf("[FCM PLACEHOLDER] Sending notification to driver %s:", driver.ID)
	log.Printf("  Token: %s", *driver.FCMToken)
	log.Printf("  Title: %s", notification.Title)
	log.Printf("  Message: %s", notification.Message)

	// In production, implement actual FCM sending:
	/*
		msg := &messaging.Message{
			Notification: &messaging.Notification{
				Title: notification.Title,
				Body:  notification.Message,
			},
			Token: *driver.FCMToken,
			Data: map[string]string{
				"bin_id": notification.BinID.String(),
				"type":   string(notification.Type),
			},
		}
		_, err := fcmClient.Send(ctx, msg)
		return err
	*/

	return nil
}

// NotifyDriver sends a notification to a specific driver
func (s *NotificationService) NotifyDriver(ctx context.Context, driverID uuid.UUID, notification *models.Notification) error {
	driver, err := s.driverRepo.GetByID(ctx, driverID)
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}
	if driver == nil {
		return fmt.Errorf("driver not found: %s", driverID)
	}

	notification.DriverID = &driverID
	return s.sendFCMNotification(driver, notification)
}

// NotifyAllAvailableDrivers broadcasts a notification to all available drivers
func (s *NotificationService) NotifyAllAvailableDrivers(ctx context.Context, notification *models.Notification) error {
	drivers, err := s.driverRepo.GetAvailableDrivers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get available drivers: %w", err)
	}

	for _, driver := range drivers {
		notificationCopy := *notification
		notificationCopy.ID = uuid.New()
		notificationCopy.DriverID = &driver.ID

		go func(d models.Driver, n *models.Notification) {
			if err := s.sendFCMNotification(&d, n); err != nil {
				log.Printf("Failed to notify driver %s: %v", d.ID, err)
			}
		}(driver, &notificationCopy)
	}

	log.Printf("Broadcast notification sent to %d available drivers", len(drivers))
	return nil
}
