package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartwaste/backend/internal/models"
	"github.com/smartwaste/backend/internal/repository"
	"github.com/smartwaste/backend/pkg/utils"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} utils.APIError
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve user")
		return
	}

	if user == nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// CreateUser creates a new user
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} utils.APIError
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	// Check if email already exists
	existing, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		utils.InternalError(c, "Failed to check existing user")
		return
	}
	if existing != nil {
		utils.Conflict(c, "Email already registered")
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: req.Password, // In production, hash this!
		FullName:     req.FullName,
		Phone:        req.Phone,
		Address:      req.Address,
		RewardPoints: 0,
	}

	if err := h.repo.Create(c.Request.Context(), user); err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, user.ToResponse())
}

// UpdateUser updates a user
// @Summary Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserRequest true "User data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} utils.APIError
// @Failure 404 {object} utils.APIError
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve user")
		return
	}
	if user == nil {
		utils.NotFound(c, "User not found")
		return
	}

	// Update fields
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Address != nil {
		user.Address = req.Address
	}

	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		utils.InternalError(c, "Failed to update user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// GetRewardPoints retrieves a user's reward points
// @Summary Get user reward points
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]int
// @Failure 404 {object} utils.APIError
// @Router /api/v1/users/{id}/rewards [get]
func (h *UserHandler) GetRewardPoints(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	points, err := h.repo.GetRewardPoints(c.Request.Context(), id)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve reward points")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user_id":       id,
		"reward_points": points,
	})
}

// AddRewardPoints adds reward points to a user
// @Summary Add reward points
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.AddRewardPointsRequest true "Points to add"
// @Success 200 {object} map[string]int
// @Failure 400 {object} utils.APIError
// @Router /api/v1/users/{id}/rewards [post]
func (h *UserHandler) AddRewardPoints(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	var req models.AddRewardPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := h.repo.UpdateRewardPoints(c.Request.Context(), id, req.Points); err != nil {
		utils.InternalError(c, "Failed to update reward points")
		return
	}

	// Get updated points
	points, _ := h.repo.GetRewardPoints(c.Request.Context(), id)

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"user_id":       id,
		"points_added":  req.Points,
		"reason":        req.Reason,
		"total_points":  points,
	})
}

// ListUsers retrieves all users with pagination
// @Summary List users
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {array} models.UserResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page := getQueryInt(c, "page", 1)
	perPage := getQueryInt(c, "per_page", 20)
	offset := (page - 1) * perPage

	users, err := h.repo.List(c.Request.Context(), perPage, offset)
	if err != nil {
		utils.InternalError(c, "Failed to retrieve users")
		return
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *u.ToResponse()
	}

	utils.SuccessResponseWithPagination(c, responses, &utils.Pagination{
		Page:    page,
		PerPage: perPage,
	})
}

// DeleteUser deletes a user
// @Summary Delete user
// @Tags Users
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} utils.APIError
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		utils.InternalError(c, "Failed to delete user")
		return
	}

	c.Status(http.StatusNoContent)
}
