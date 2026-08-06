package handler

import (
	"github.com/gofiber/fiber/v3"

	internalpkg "universev/internal/pkg"
	"universev/pkg/response"
)

// Wrapper functions that delegate to internal/pkg for backwards compatibility

func splitAuthHeader(authHeader string) []string {
	return internalpkg.SplitAuthHeader(authHeader)
}

func isValidEmail(email string) bool {
	return internalpkg.IsValidEmail(email)
}

func isValidNIK(nik string) bool {
	return internalpkg.IsValidNIK(nik)
}

func isTrimmedEmpty(s string) bool {
	return internalpkg.IsTrimmedEmpty(s)
}

func isPasswordStrong(pw string) bool {
	return internalpkg.IsPasswordStrong(pw)
}

func sendValidationError(c fiber.Ctx, field, msg string) error {
	return response.ValidationError(c, []response.ErrorDetail{
		{Field: field, Message: msg},
	})
}
