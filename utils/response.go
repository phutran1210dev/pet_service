package utils

import (
	"net/http"
	"pet-service/dto"

	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// CreatedResponse sends a 201 Created JSON response
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// ErrorResponse sends an error JSON response with custom status code
func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// BadRequestError sends a 400 Bad Request error
func BadRequestError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusBadRequest, code, message)
}

// UnauthorizedError sends a 401 Unauthorized error
func UnauthorizedError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusUnauthorized, code, message)
}

// ForbiddenError sends a 403 Forbidden error
func ForbiddenError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusForbidden, code, message)
}

// NotFoundError sends a 404 Not Found error
func NotFoundError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusNotFound, code, message)
}

// ConflictError sends a 409 Conflict error
func ConflictError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusConflict, code, message)
}

// InternalServerError sends a 500 Internal Server Error
func InternalServerError(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusInternalServerError, code, message)
}

// ValidationError sends a validation error with field details
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, NewValidationErrorResponse(err))
}
