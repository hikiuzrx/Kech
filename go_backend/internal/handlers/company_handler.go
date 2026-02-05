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

// CompanyHandler handles company-related HTTP requests
type CompanyHandler struct {
	companyRepo    *repository.CompanyRepository
	pricingRepo    *repository.PricingRepository
	valuationSvc   *services.ValuationService
}

// NewCompanyHandler creates a new CompanyHandler
func NewCompanyHandler(
	companyRepo *repository.CompanyRepository,
	pricingRepo *repository.PricingRepository,
	valuationSvc *services.ValuationService,
) *CompanyHandler {
	return &CompanyHandler{
		companyRepo:  companyRepo,
		pricingRepo:  pricingRepo,
		valuationSvc: valuationSvc,
	}
}

// GetCompany retrieves a company by ID
// @Summary Get company by ID
// @Tags Companies
// @Produce json
// @Param id path string true "Company ID"
// @Success 200 {object} models.CompanyResponse
// @Failure 404 {object} utils.APIError
// @Router /api/v1/companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid company ID format")
		return
	}

	company, err := h.companyRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve company")
		return
	}

	if company == nil {
		utils.NotFound(c, "Company not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, company.ToResponse())
}

// CreateCompany creates a new company
// @Summary Create a new company
// @Tags Companies
// @Accept json
// @Produce json
// @Param company body models.CreateCompanyRequest true "Company data"
// @Success 201 {object} models.CompanyResponse
// @Failure 400 {object} utils.APIError
// @Router /api/v1/companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// Check if email already exists
	existing, err := h.companyRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		utils.InternalError(c, "Failed to check existing company")
		return
	}
	if existing != nil {
		utils.Conflict(c, "Email already registered")
		return
	}

	company := &models.Company{
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Phone,
		Address:            req.Address,
		City:               req.City,
		Country:            req.Country,
		RegistrationNumber: req.RegistrationNumber,
		IsActive:           true,
	}

	if err := h.companyRepo.Create(c.Request.Context(), company); err != nil {
		utils.InternalError(c, "Failed to create company")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, company.ToResponse())
}

// UpdateCompany updates a company
// @Summary Update company
// @Tags Companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company body models.UpdateCompanyRequest true "Company data"
// @Success 200 {object} models.CompanyResponse
// @Router /api/v1/companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid company ID format")
		return
	}

	var req models.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	company, err := h.companyRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve company")
		return
	}
	if company == nil {
		utils.NotFound(c, "Company not found")
		return
	}

	// Update fields
	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Email != nil {
		company.Email = *req.Email
	}
	if req.Phone != nil {
		company.Phone = req.Phone
	}
	if req.Address != nil {
		company.Address = req.Address
	}
	if req.City != nil {
		company.City = req.City
	}
	if req.Country != nil {
		company.Country = req.Country
	}
	if req.RegistrationNumber != nil {
		company.RegistrationNumber = req.RegistrationNumber
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}

	if err := h.companyRepo.Update(c.Request.Context(), company); err != nil {
		utils.InternalError(c, "Failed to update company")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, company.ToResponse())
}

// ListCompanies retrieves all companies with pagination
// @Summary List companies
// @Tags Companies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {array} models.CompanyResponse
// @Router /api/v1/companies [get]
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 20)
	offset := (page - 1) * perPage

	companies, err := h.companyRepo.List(c.Request.Context(), perPage, offset)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve companies")
		return
	}

	responses := make([]models.CompanyResponse, len(companies))
	for i, comp := range companies {
		responses[i] = *comp.ToResponse()
	}

	utils.SuccessResponseWithPagination(c, responses, &utils.Pagination{
		Page:    page,
		PerPage: perPage,
	})
}

// DeleteCompany deletes a company
// @Summary Delete company
// @Tags Companies
// @Param id path string true "Company ID"
// @Success 204 "No Content"
// @Router /api/v1/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid company ID format")
		return
	}

	if err := h.companyRepo.Delete(c.Request.Context(), id); err != nil {
		utils.InternalError(c, "Failed to delete company")
		return
	}

	c.Status(http.StatusNoContent)
}

// --- Pricing Rules ---

// GetPricingRule retrieves a pricing rule by ID
// @Summary Get pricing rule by ID
// @Tags Pricing Rules
// @Produce json
// @Param id path string true "Pricing Rule ID"
// @Success 200 {object} models.PricingRuleResponse
// @Router /api/v1/pricing-rules/{id} [get]
func (h *CompanyHandler) GetPricingRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid pricing rule ID format")
		return
	}

	rule, err := h.pricingRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve pricing rule")
		return
	}

	if rule == nil {
		utils.NotFound(c, "Pricing rule not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, rule.ToResponse())
}

// CreatePricingRule creates a new pricing rule
// @Summary Create a new pricing rule
// @Tags Pricing Rules
// @Accept json
// @Produce json
// @Param rule body models.CreatePricingRuleRequest true "Pricing rule data"
// @Success 201 {object} models.PricingRuleResponse
// @Router /api/v1/pricing-rules [post]
func (h *CompanyHandler) CreatePricingRule(c *gin.Context) {
	var req models.CreatePricingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	rule := &models.PricingRule{
		WasteType:   req.WasteType,
		Condition:   req.Condition,
		PricePerKg:  req.PricePerKg,
		Currency:    req.Currency,
		MinWeightKg: req.MinWeightKg,
		MaxWeightKg: req.MaxWeightKg,
		CompanyID:   req.CompanyID,
		IsActive:    true,
	}

	if err := h.pricingRepo.Create(c.Request.Context(), rule); err != nil {
		utils.InternalError(c, "Failed to create pricing rule")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, rule.ToResponse())
}

// UpdatePricingRule updates a pricing rule
// @Summary Update pricing rule
// @Tags Pricing Rules
// @Accept json
// @Produce json
// @Param id path string true "Pricing Rule ID"
// @Param rule body models.UpdatePricingRuleRequest true "Pricing rule data"
// @Success 200 {object} models.PricingRuleResponse
// @Router /api/v1/pricing-rules/{id} [put]
func (h *CompanyHandler) UpdatePricingRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid pricing rule ID format")
		return
	}

	var req models.UpdatePricingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	rule, err := h.pricingRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve pricing rule")
		return
	}
	if rule == nil {
		utils.NotFound(c, "Pricing rule not found")
		return
	}

	// Update fields
	if req.WasteType != nil {
		rule.WasteType = *req.WasteType
	}
	if req.Condition != nil {
		rule.Condition = *req.Condition
	}
	if req.PricePerKg != nil {
		rule.PricePerKg = *req.PricePerKg
	}
	if req.Currency != nil {
		rule.Currency = *req.Currency
	}
	if req.MinWeightKg != nil {
		rule.MinWeightKg = *req.MinWeightKg
	}
	if req.MaxWeightKg != nil {
		rule.MaxWeightKg = req.MaxWeightKg
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	if err := h.pricingRepo.Update(c.Request.Context(), rule); err != nil {
		utils.InternalError(c, "Failed to update pricing rule")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, rule.ToResponse())
}

// ListPricingRules retrieves all pricing rules
// @Summary List pricing rules
// @Tags Pricing Rules
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(50)
// @Success 200 {array} models.PricingRuleResponse
// @Router /api/v1/pricing-rules [get]
func (h *CompanyHandler) ListPricingRules(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 50)
	offset := (page - 1) * perPage

	rules, err := h.pricingRepo.List(c.Request.Context(), perPage, offset)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve pricing rules")
		return
	}

	responses := make([]models.PricingRuleResponse, len(rules))
	for i, r := range rules {
		responses[i] = *r.ToResponse()
	}

	utils.SuccessResponseWithPagination(c, responses, &utils.Pagination{
		Page:    page,
		PerPage: perPage,
	})
}

// DeletePricingRule deletes a pricing rule
// @Summary Delete pricing rule
// @Tags Pricing Rules
// @Param id path string true "Pricing Rule ID"
// @Success 204 "No Content"
// @Router /api/v1/pricing-rules/{id} [delete]
func (h *CompanyHandler) DeletePricingRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid pricing rule ID format")
		return
	}

	if err := h.pricingRepo.Delete(c.Request.Context(), id); err != nil {
		utils.InternalError(c, "Failed to delete pricing rule")
		return
	}

	c.Status(http.StatusNoContent)
}

// CalculateValuation calculates waste valuation
// @Summary Calculate waste valuation
// @Tags Pricing Rules
// @Accept json
// @Produce json
// @Param request body models.ValuationRequest true "Valuation request"
// @Success 200 {object} models.ValuationResponse
// @Router /api/v1/valuations [post]
func (h *CompanyHandler) CalculateValuation(c *gin.Context) {
	var req models.ValuationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	result, err := h.valuationSvc.CalculateValue(c.Request.Context(), &req)
	if err != nil {
		utils.InternalError(c, "Failed to calculate valuation")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}
