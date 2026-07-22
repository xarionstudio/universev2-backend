package handler

import (
	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/pkg"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type NotificationHandler struct {
	repo *repository.NotificationRepo
}

func NewNotificationHandler(repo *repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

func (h *NotificationHandler) GetNotifications(c fiber.Ctx) error {
	userID := "u1"
	if claims, ok := c.Locals("user").(*pkg.JWTCustomClaims); ok && claims != nil {
		userID = claims.UserID
	}

	notifs, err := h.repo.GetByUser(userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch notifications: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch notifications", notifs)
}

func (h *NotificationHandler) MarkRead(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.MarkRead(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to mark notification as read: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllRead(c fiber.Ctx) error {
	userID := "u1"
	if claims, ok := c.Locals("user").(*pkg.JWTCustomClaims); ok && claims != nil {
		userID = claims.UserID
	}

	if err := h.repo.MarkAllRead(userID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to mark all notifications as read: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "All notifications marked as read", nil)
}
