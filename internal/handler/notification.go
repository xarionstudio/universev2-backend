package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev/internal/pkg"
	"universev/internal/repository"
	"universev/pkg/response"
)

type NotificationHandler struct {
	repo *repository.NotificationRepo
}

func NewNotificationHandler(repo *repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) GetNotifications(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	notifs, err := h.repo.GetByUser(claims.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch notifications: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch notifications", notifs)
}

func (h *NotificationHandler) MarkRead(c fiber.Ctx) error {
	id := c.Params("id")
	if isTrimmedEmpty(id) {
		return response.Error(c, fiber.StatusBadRequest, "Notification ID is required")
	}

	if err := h.repo.MarkRead(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to mark notification as read: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllRead(c fiber.Ctx) error {
	claims, ok := c.Locals("user").(*pkg.JWTCustomClaims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if err := h.repo.MarkAllRead(claims.UserID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to mark all notifications as read: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "All notifications marked as read", nil)
}
