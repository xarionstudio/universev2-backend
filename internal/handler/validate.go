package handler

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/pkg/response"
)

func splitAuthHeader(authHeader string) []string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}
	return parts
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nikRegex   = regexp.MustCompile(`^\d{9}$`)
)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

func isValidNIK(nik string) bool {
	return nikRegex.MatchString(strings.TrimSpace(nik))
}

func isTrimmedEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func isPasswordStrong(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	hasDigit := false
	hasLetter := false
	for _, r := range pw {
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
	}
	return hasDigit && hasLetter
}

func sendValidationError(c fiber.Ctx, field, msg string) error {
	return response.ValidationError(c, []response.ErrorDetail{
		{Field: field, Message: msg},
	})
}
