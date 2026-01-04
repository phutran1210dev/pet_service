package utils

import (
	"fmt"
	"pet-service/dto"

	"github.com/go-playground/validator/v10"
)

// Error codes following Golang naming conventions
const (
	// Authentication & Authorization
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeInvalidToken     = "INVALID_TOKEN"
	ErrCodeTokenExpired     = "TOKEN_EXPIRED"
	ErrCodeInvalidPassword  = "INVALID_PASSWORD"
	ErrCodePermissionDenied = "PERMISSION_DENIED"

	// User errors
	ErrCodeUserNotFound = "USER_NOT_FOUND"
	ErrCodeEmailTaken   = "EMAIL_TAKEN"

	// Validation errors
	ErrCodeValidationFailed = "VALIDATION_ERROR"
	ErrCodeInvalidInput     = "INVALID_INPUT"

	// Resource errors
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodePetNotFound   = "PET_NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"

	// Server errors
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeServiceError  = "SERVICE_ERROR"
)

// Error messages
const (
	UserIsNotExist      = "User does not exist"
	PasswordInvalid     = "Invalid password"
	LoginError          = "Login error"
	ErrorInvalidToken   = "Invalid token"
	TokenExpired        = "Token has expired"
	JTINotExist         = "JTI does not exist"
	JTIInBlacklist      = "Token has been revoked"
	ServiceError        = "Service error"
	PetIDNotExist = "Pet ID does not exist"
	EmailTaken    = "Email is already taken"
	PermissionDenied    = "Permission denied"
	UserHasNoPermission = "User has no permissions"
	InvalidRequestBody  = "Invalid request body"
	ValidationFailed    = "Validation failed"
)

// NewErrorResponse creates a standard error response
func NewErrorResponse(code, message string) dto.ErrorResponse {
	return dto.ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewValidationErrorResponse creates a validation error response with field details
func NewValidationErrorResponse(err error) dto.ErrorResponse {
	details := FormatValidationErrors(err)
	return dto.ErrorResponse{
		Code:    ErrCodeValidationFailed,
		Message: ValidationFailed,
		Details: details,
	}
}

// FormatValidationErrors formats Gin validation errors to user-friendly messages
func FormatValidationErrors(err error) []dto.ErrorDetail {
	var errors []dto.ErrorDetail
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			var message string
			field := fieldError.Field()
			tag := fieldError.Tag()
			
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				message = fmt.Sprintf("%s must not exceed %s characters", field, fieldError.Param())
			default:
				message = fmt.Sprintf("%s is invalid", field)
			}
			
			errors = append(errors, dto.ErrorDetail{
				Field:   field,
				Message: message,
			})
		}
	} else {
		errors = append(errors, dto.ErrorDetail{
			Field:   "body",
			Message: err.Error(),
		})
	}
	
	return errors
}
