package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Sentry holds a set of common, user-friendly error messages so that
// internal / technical details never leak to the API consumer.
var (
	// ── generic ────────────────────────────────────────────────
	ErrInternal     = errors.New("An internal error occurred. Please try again later")
	ErrBadRequest   = errors.New("Invalid request. Please check your input and try again")
	ErrNotFound     = errors.New("The requested resource was not found")
	ErrUnauthorized = errors.New("You are not authorized to perform this action")
	ErrForbidden    = errors.New("You do not have permission to access this resource")
	ErrConflict     = errors.New("The resource already exists")
	ErrTimeout      = errors.New("The request timed out. Please try again")

	// ── auth ──────────────────────────────────────────────────
	ErrInvalidCredentials = errors.New("Email or password is incorrect")
	ErrAccountInactive    = errors.New("Your account is inactive. Please contact administrator")
	ErrEmailRegistered    = errors.New("This email is already registered")
	ErrNIKRegistered      = errors.New("This NIK is already registered")
	ErrTokenInvalid       = errors.New("Your session has expired. Please login again")
	ErrTokenMissing       = errors.New("Authentication is required. Please login first")
	ErrPasswordWeak       = errors.New("Password must be at least 8 characters and contain both letters and numbers")
	ErrPasswordTooLong    = errors.New("Password must not exceed 72 characters")
	ErrNameTooLong        = errors.New("Name must not exceed 100 characters")
	ErrNIKTooLong         = errors.New("NIK must not exceed 9 characters")

	// ── employee ──────────────────────────────────────────────
	ErrEmployeeNotFound = errors.New("Employee not found")
	ErrEmployeeExists   = errors.New("Employee with this NIK already exists")
	ErrNIKInvalid       = errors.New("NIK must be exactly 9 digits")
	ErrNameRequired     = errors.New("Employee name is required")
	ErrDeptRequired     = errors.New("Department is required")
	ErrPosRequired      = errors.New("Position is required")
	ErrEmailRequired    = errors.New("Email address is required")
	ErrEmailInvalid     = errors.New("Invalid email format")

	// ── file upload ───────────────────────────────────────────
	ErrFileRequired        = errors.New("Please upload a file")
	ErrFileTooLarge        = errors.New("File size exceeds the maximum limit")
	ErrFileTypeNotAllowed  = errors.New("File type is not supported")
	ErrPhotoRequired       = errors.New("Please upload a photo file")
	ErrPhotoTypeNotAllowed = errors.New("Only JPEG and PNG files are allowed")
	ErrPhotoTooLarge       = errors.New("Photo size exceeds the maximum limit of 5MB")

	// ── attendance ────────────────────────────────────────────
	ErrNIKRequired       = errors.New("NIK is required")
	ErrDateRangeRequired = errors.New("Start date ('from') and end date ('to') are required")

	// ── fleet ─────────────────────────────────────────────────
	ErrUnitCodeRequired        = errors.New("Unit code is required")
	ErrStatusRequired          = errors.New("Status is required")
	ErrStatusInvalid           = errors.New("Status must be one of: ready, breakdown, standby")
	ErrBreakdownReasonRequired = errors.New("Breakdown reason is required")
	ErrDiggerCodeRequired      = errors.New("Digger code is required")

	// ── FTW ───────────────────────────────────────────────────
	ErrFTWNIKRequired   = errors.New("NIK is required")
	ErrFTWShiftRequired = errors.New("Shift is required")
	ErrFTWNIKQuery      = errors.New("Please provide an employee NIK to view history")

	// ── roster ────────────────────────────────────────────────
	ErrRosterKeyRequired    = errors.New("Roster key is required")
	ErrSubmissionIDRequired = errors.New("Submission ID is required")
	ErrInvalidRevisionID    = errors.New("Invalid revision ID")

	// ── general validation ────────────────────────────────────
	ErrValidationFailed = errors.New("Validation failed. Please check your input")
	ErrInvalidJSON      = errors.New("Invalid request format. Please check your input and try again")
	ErrInvalidPayload   = errors.New("Invalid request body")
)

// ErrorTitles maps error types to user-friendly titles.
var ErrorTitles = map[error]string{
	ErrInternal:         "Internal Error",
	ErrBadRequest:       "Bad Request",
	ErrNotFound:         "Not Found",
	ErrUnauthorized:     "Unauthorized",
	ErrForbidden:        "Forbidden",
	ErrConflict:         "Conflict",
	ErrTimeout:          "Request Timeout",
	ErrValidationFailed: "Validation Error",
	ErrInvalidJSON:      "Invalid Format",
	ErrInvalidPayload:   "Invalid Request",
}

// AppError wraps an user-facing error with an optional code.
type AppError struct {
	Err        error
	UserMsg    string
	StatusCode int
	Field      string // optional field name for validation errors
}

func (e *AppError) Error() string { return e.UserMsg }

// MapHTTPError takes a domain error and returns the appropriate HTTP status code.
func MapHTTPError(err error) int {
	if err == nil {
		return fiber.StatusOK
	}

	msg := err.Error()
	msgLower := strings.ToLower(msg)

	// 4xx
	if strings.Contains(msgLower, "not found") || strings.Contains(msgLower, "not_found") {
		return fiber.StatusNotFound
	}
	if strings.Contains(msgLower, "invalid") || strings.Contains(msgLower, "required") || strings.Contains(msgLower, "must be") {
		return fiber.StatusBadRequest
	}
	if strings.Contains(msgLower, "unauthorized") || strings.Contains(msgLower, "invalid credentials") || strings.Contains(msgLower, "invalid or expired token") {
		return fiber.StatusUnauthorized
	}
	if strings.Contains(msgLower, "forbidden") || strings.Contains(msgLower, "inactive") {
		return fiber.StatusForbidden
	}
	if strings.Contains(msgLower, "conflict") || strings.Contains(msgLower, "already exists") || strings.Contains(msgLower, "already registered") {
		return fiber.StatusConflict
	}
	if strings.Contains(msgLower, "validation") {
		return fiber.StatusUnprocessableEntity
	}

	// 5xx
	return fiber.StatusInternalServerError
}

// IsContextError returns true when the HTTP context was cancelled / deadline exceeded,
// so callers can skip sending a second response.
func IsContextError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "context canceled") ||
		strings.Contains(msg, "context deadline exceeded") ||
		strings.Contains(msg, "request canceled")
}

// ── helper: send a predefined error ──────────────────────────────────────────

// SendError sends a user-friendly error response using a predefined AppError or error constant.
// If the error does not match any known error, it falls back to the generic internal error.
func SendError(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// If it's already an AppError, use its fields directly.
	var appErr *AppError
	if errors.As(err, &appErr) {
		return sendTypedError(c, appErr.StatusCode, appErr.UserMsg, appErr.Field)
	}

	// Map known errors to friendly messages and status codes.
	friendly, status := FriendlyError(err)
	return sendTypedError(c, status, friendly, "")
}

// FriendlyError converts any error into a user-friendly message and HTTP status code.
func FriendlyError(err error) (message string, statusCode int) {
	if err == nil {
		return "", fiber.StatusOK
	}

	statusCode = MapHTTPError(err)
	msg := err.Error()

	// ── known sentinel errors ─────────────────────────────────
	switch {
	// auth
	case errors.Is(err, ErrInvalidCredentials):
		return ErrInvalidCredentials.Error(), fiber.StatusUnauthorized
	case errors.Is(err, ErrAccountInactive):
		return ErrAccountInactive.Error(), fiber.StatusForbidden
	case errors.Is(err, ErrEmailRegistered):
		return ErrEmailRegistered.Error(), fiber.StatusConflict
	case errors.Is(err, ErrNIKRegistered):
		return ErrNIKRegistered.Error(), fiber.StatusConflict
	case errors.Is(err, ErrTokenInvalid):
		return ErrTokenInvalid.Error(), fiber.StatusUnauthorized
	case errors.Is(err, ErrTokenMissing):
		return ErrTokenMissing.Error(), fiber.StatusUnauthorized
	case errors.Is(err, ErrPasswordWeak):
		return sendValidationMsg("password", ErrPasswordWeak.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrPasswordTooLong):
		return sendValidationMsg("password", ErrPasswordTooLong.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrNameTooLong):
		return sendValidationMsg("name", ErrNameTooLong.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrNIKTooLong):
		return sendValidationMsg("nik", ErrNIKTooLong.Error()), fiber.StatusUnprocessableEntity

	// employee
	case errors.Is(err, ErrEmployeeNotFound):
		return ErrEmployeeNotFound.Error(), fiber.StatusNotFound
	case errors.Is(err, ErrEmployeeExists):
		return ErrEmployeeExists.Error(), fiber.StatusConflict
	case errors.Is(err, ErrNIKInvalid):
		return sendValidationMsg("nik", ErrNIKInvalid.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrNameRequired):
		return sendValidationMsg("name", ErrNameRequired.Error()), fiber.StatusUnprocessableEntity

	// file
	case errors.Is(err, ErrFileRequired):
		return ErrFileRequired.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrPhotoRequired):
		return ErrPhotoRequired.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrPhotoTypeNotAllowed):
		return ErrPhotoTypeNotAllowed.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrPhotoTooLarge):
		return ErrPhotoTooLarge.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrFileTooLarge):
		return ErrFileTooLarge.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrFileTypeNotAllowed):
		return ErrFileTypeNotAllowed.Error(), fiber.StatusBadRequest

	// attendance
	case errors.Is(err, ErrNIKRequired):
		return sendValidationMsg("nik", ErrNIKRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrDateRangeRequired):
		return ErrDateRangeRequired.Error(), fiber.StatusBadRequest

	// fleet
	case errors.Is(err, ErrUnitCodeRequired):
		return sendValidationMsg("code", ErrUnitCodeRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrStatusRequired):
		return sendValidationMsg("status", ErrStatusRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrStatusInvalid):
		return sendValidationMsg("status", ErrStatusInvalid.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrBreakdownReasonRequired):
		return sendValidationMsg("reason", ErrBreakdownReasonRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrDiggerCodeRequired):
		return sendValidationMsg("digger", ErrDiggerCodeRequired.Error()), fiber.StatusUnprocessableEntity

	// FTW
	case errors.Is(err, ErrFTWNIKRequired):
		return sendValidationMsg("nik", ErrFTWNIKRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrFTWShiftRequired):
		return sendValidationMsg("shift", ErrFTWShiftRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrFTWNIKQuery):
		return ErrFTWNIKQuery.Error(), fiber.StatusBadRequest

	// roster
	case errors.Is(err, ErrRosterKeyRequired):
		return ErrRosterKeyRequired.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrSubmissionIDRequired):
		return sendValidationMsg("sid", ErrSubmissionIDRequired.Error()), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrInvalidRevisionID):
		return ErrInvalidRevisionID.Error(), fiber.StatusBadRequest

	// validation / json
	case errors.Is(err, ErrValidationFailed):
		return ErrValidationFailed.Error(), fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrInvalidJSON), errors.Is(err, ErrInvalidPayload):
		return ErrInvalidPayload.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrBadRequest):
		return ErrBadRequest.Error(), fiber.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		return ErrNotFound.Error(), fiber.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return ErrUnauthorized.Error(), fiber.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return ErrForbidden.Error(), fiber.StatusForbidden
	case errors.Is(err, ErrConflict):
		return ErrConflict.Error(), fiber.StatusConflict
	case errors.Is(err, ErrInternal):
		return ErrInternal.Error(), fiber.StatusInternalServerError
	}

	// ── string-based fallback ─────────────────────────────────
	msgLower := strings.ToLower(msg)
	switch {
	case strings.Contains(msgLower, "email is already registered"):
		return ErrEmailRegistered.Error(), fiber.StatusConflict
	case strings.Contains(msgLower, "nik is already registered"):
		return ErrNIKRegistered.Error(), fiber.StatusConflict
	case strings.Contains(msgLower, "invalid credentials"):
		return ErrInvalidCredentials.Error(), fiber.StatusUnauthorized
	case strings.Contains(msgLower, "account is inactive"):
		return ErrAccountInactive.Error(), fiber.StatusForbidden
	case strings.Contains(msgLower, "invalid or expired token"):
		return ErrTokenInvalid.Error(), fiber.StatusUnauthorized
	case strings.Contains(msgLower, "user not found"), strings.Contains(msgLower, "employee not found"):
		return ErrEmployeeNotFound.Error(), fiber.StatusNotFound
	case strings.Contains(msgLower, "nik must be exactly 9 digits"):
		return sendValidationMsg("nik", ErrNIKInvalid.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "name is required"):
		return sendValidationMsg("name", ErrNameRequired.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "password must be at least 8 characters"):
		return sendValidationMsg("password", ErrPasswordWeak.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "password must not exceed 72 characters"):
		return sendValidationMsg("password", ErrPasswordTooLong.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "name must not exceed 100 characters"):
		return sendValidationMsg("name", ErrNameTooLong.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "nik must not exceed 9 characters"):
		return sendValidationMsg("nik", ErrNIKTooLong.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "photo file is required"):
		return ErrPhotoRequired.Error(), fiber.StatusBadRequest
	case strings.Contains(msgLower, "only jpeg and png files are allowed"):
		return ErrPhotoTypeNotAllowed.Error(), fiber.StatusBadRequest
	case strings.Contains(msgLower, "file size exceeds") && strings.Contains(msgLower, "5mb"):
		return ErrPhotoTooLarge.Error(), fiber.StatusBadRequest
	case strings.Contains(msgLower, "date is required"):
		return sendValidationMsg("date", "Date is required"), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "shift is required"):
		return sendValidationMsg("shift", "Shift is required"), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "unit code is required"):
		return sendValidationMsg("code", ErrUnitCodeRequired.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "digger code is required"):
		return sendValidationMsg("digger", ErrDiggerCodeRequired.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "status must be one of"):
		return sendValidationMsg("status", ErrStatusInvalid.Error()), fiber.StatusUnprocessableEntity
	case strings.Contains(msgLower, "breakdown reason is required"):
		return sendValidationMsg("reason", ErrBreakdownReasonRequired.Error()), fiber.StatusUnprocessableEntity
	}

	// ── generic fallback with internal error ──────────────────
	return ErrInternal.Error(), fiber.StatusInternalServerError
}

// sendTypedError writes an error/validation-error response depending on whether a field was supplied.
func sendTypedError(c fiber.Ctx, status int, message string, field string) error {
	if field != "" {
		return ValidationError(c, []ErrorDetail{
			{Field: field, Message: message},
		})
	}
	return Error(c, status, message)
}

// sendValidationMsg prefixes the message with field info for internal use.
func sendValidationMsg(field, msg string) string {
	return fmt.Sprintf("%s: %s", field, msg)
}

// ExtractUserMsg extracts the user-friendly part from a validation-style message.
// For example "nik: NIK is required" becomes "NIK is required".
func ExtractUserMsg(msg string) string {
	idx := strings.Index(msg, ": ")
	if idx > 0 && idx+2 < len(msg) {
		return msg[idx+2:]
	}
	return msg
}

// ExtractField extracts the field name from a validation-style message.
// For example "nik: NIK is required" returns "nik".
func ExtractField(msg string) string {
	idx := strings.Index(msg, ": ")
	if idx > 0 {
		return msg[:idx]
	}
	return ""
}

// String-based helpers for backward compatibility with existing services.
func IsValidationMsg(msg string) bool {
	return strings.Contains(msg, ": ")
}
