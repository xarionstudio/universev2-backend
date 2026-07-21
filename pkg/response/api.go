package response

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp"`
}

type Meta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success   bool          `json:"success"`
	Message   string        `json:"message"`
	Errors    []ErrorDetail `json:"errors,omitempty"`
	Timestamp string        `json:"timestamp"`
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func Success(c fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: now(),
	})
}

func SuccessWithMeta(c fiber.Ctx, status int, message string, data interface{}, meta *Meta) error {
	return c.Status(status).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: now(),
	})
}

func Error(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{
		Success:   false,
		Message:   message,
		Timestamp: now(),
	})
}

func ErrorWithDetails(c fiber.Ctx, status int, message string, errors []ErrorDetail) error {
	return c.Status(status).JSON(ErrorResponse{
		Success:   false,
		Message:   message,
		Errors:    errors,
		Timestamp: now(),
	})
}

func ValidationError(c fiber.Ctx, errors []ErrorDetail) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
		Success:   false,
		Message:   "Validation failed",
		Errors:    errors,
		Timestamp: now(),
	})
}