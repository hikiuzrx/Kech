package services

import (
	"context"
	"fmt"

	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
)

// ValuationService handles waste valuation based on pricing rules
type ValuationService struct {
	pricingRepo *repository.PricingRepository
}

// NewValuationService creates a new ValuationService
func NewValuationService(pricingRepo *repository.PricingRepository) *ValuationService {
	return &ValuationService{
		pricingRepo: pricingRepo,
	}
}

// CalculateValue calculates the value of waste based on type, condition, and weight
func (s *ValuationService) CalculateValue(ctx context.Context, req *models.ValuationRequest) (*models.ValuationResponse, error) {
	// Find applicable pricing rule
	rule, err := s.pricingRepo.GetByTypeAndCondition(ctx, req.WasteType, req.Condition)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pricing rule: %w", err)
	}

	if rule == nil {
		// No specific rule found, return default pricing
		return &models.ValuationResponse{
			WasteType:  req.WasteType,
			Condition:  req.Condition,
			WeightKg:   req.WeightKg,
			PricePerKg: 0,
			TotalPrice: 0,
			Currency:   "USD",
			Message:    "No pricing rule found for this waste type and condition",
		}, nil
	}

	// Validate weight against rule constraints
	if req.WeightKg < rule.MinWeightKg {
		return &models.ValuationResponse{
			WasteType:  req.WasteType,
			Condition:  req.Condition,
			WeightKg:   req.WeightKg,
			PricePerKg: rule.PricePerKg,
			TotalPrice: 0,
			Currency:   rule.Currency,
			Message:    fmt.Sprintf("Weight below minimum threshold of %.2f kg", rule.MinWeightKg),
		}, nil
	}

	if rule.MaxWeightKg != nil && req.WeightKg > *rule.MaxWeightKg {
		return &models.ValuationResponse{
			WasteType:  req.WasteType,
			Condition:  req.Condition,
			WeightKg:   req.WeightKg,
			PricePerKg: rule.PricePerKg,
			TotalPrice: 0,
			Currency:   rule.Currency,
			Message:    fmt.Sprintf("Weight exceeds maximum threshold of %.2f kg", *rule.MaxWeightKg),
		}, nil
	}

	// Calculate total value
	totalPrice := req.WeightKg * rule.PricePerKg
	ruleID := rule.ID.String()

	return &models.ValuationResponse{
		WasteType:     req.WasteType,
		Condition:     req.Condition,
		WeightKg:      req.WeightKg,
		PricePerKg:    rule.PricePerKg,
		TotalPrice:    totalPrice,
		Currency:      rule.Currency,
		PricingRuleID: &ruleID,
		Message:       "Valuation calculated successfully",
	}, nil
}

// ValuateWasteMetadata valuates waste based on AI-detected metadata
func (s *ValuationService) ValuateWasteMetadata(ctx context.Context, metadata *models.WasteMetadata, weightKg float64) (*models.ValuationResponse, error) {
	req := &models.ValuationRequest{
		WasteType: metadata.WasteType,
		Condition: metadata.Condition,
		WeightKg:  weightKg,
	}
	return s.CalculateValue(ctx, req)
}

// GetPricingRules returns all active pricing rules
func (s *ValuationService) GetPricingRules(ctx context.Context, limit, offset int) ([]models.PricingRule, error) {
	return s.pricingRepo.List(ctx, limit, offset)
}

// Common waste types for reference
const (
	WasteTypePlastic    = "plastic"
	WasteTypePaper      = "paper"
	WasteTypeGlass      = "glass"
	WasteTypeMetal      = "metal"
	WasteTypeOrganic    = "organic"
	WasteTypeElectronic = "electronic"
	WasteTypeTextile    = "textile"
	WasteTypeGeneral    = "general"
)

// Common waste conditions for reference
const (
	ConditionExcellent = "excellent"
	ConditionGood      = "good"
	ConditionFair      = "fair"
	ConditionPoor      = "poor"
	ConditionDamaged   = "damaged"
)
