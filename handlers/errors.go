package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  string      `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationErrorDetails represents details for validation errors
type ValidationErrorDetails struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Details: details,
	})
}

// HandleValidationErrors processes validation errors and returns a standardized response
func HandleValidationErrors(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Get the first validation error
		firstErr := validationErrors[0]

		// Extract field name (remove struct name prefix)
		fieldName := firstErr.Field()
		if idx := strings.LastIndex(fieldName, "."); idx != -1 {
			fieldName = fieldName[idx+1:]
		}

		// Convert to JSON field name (lowercase first letter)
		jsonFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]

		// Create error details
		details := ValidationErrorDetails{
			Field:  jsonFieldName,
			Reason: firstErr.Tag(),
		}

		// Create user-friendly message
		var message string
		switch firstErr.Tag() {
		case "required":
			message = "The " + jsonFieldName + " is required"
		default:
			message = "Validation failed for field " + jsonFieldName
		}

		RespondWithError(c, http.StatusBadRequest, "MISSING_REQUIRED_FIELD", message, details)
		return
	}

	// Handle non-validation errors
	RespondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), nil)
}
