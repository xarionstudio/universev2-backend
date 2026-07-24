package pkg

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nikRegex   = regexp.MustCompile(`^\d{9}$`)
)

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.TrimSpace(email))
}

// IsValidNIK checks if the NIK is exactly 9 digits
func IsValidNIK(nik string) bool {
	return nikRegex.MatchString(strings.TrimSpace(nik))
}

// IsTrimmedEmpty checks if a string is empty after trimming spaces
func IsTrimmedEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsPasswordStrong checks if password has min 8 chars with both letters and digits
func IsPasswordStrong(pw string) bool {
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

// SplitAuthHeader splits "Bearer <token>" into parts
func SplitAuthHeader(authHeader string) []string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}
	return parts
}
