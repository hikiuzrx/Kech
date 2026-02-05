package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithPagination sends a successful response with pagination
func SuccessResponseWithPagination(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta:    pagination,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorResponseWithDetails sends an error response with additional details
func ErrorResponseWithDetails(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// Common error codes
const (
	ErrCodeBadRequest       = "BAD_REQUEST"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeValidationFailed = "VALIDATION_FAILED"
)

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, ErrCodeBadRequest, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, message)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, ErrCodeInternalError, message)
}

// ValidationError sends a 400 response for validation errors
func ValidationError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, ErrCodeValidationFailed, message)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, ErrCodeConflict, message)
}
