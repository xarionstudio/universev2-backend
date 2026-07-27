package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/pkg/response"
)

type PrestasiHandler struct {
	svc *service.PrestasiService
}

func NewPrestasiHandler(repo *repository.PrestasiRepo) *PrestasiHandler {
	return &PrestasiHandler{svc: service.NewPrestasiService(repo)}
}

func (h *PrestasiHandler) GetLeaderboard(c fiber.Ctx) error {
	daysStr := c.Query("days", "30")
	days, _ := strconv.Atoi(daysStr)

	records, err := h.svc.GetLeaderboard(days)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch leaderboard: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch leaderboard", records)
}

func (h *PrestasiHandler) GetOperatorHistory(c fiber.Ctx) error {
	nik := c.Params("nik")
	daysStr := c.Query("days", "90")
	days, _ := strconv.Atoi(daysStr)

	history, err := h.svc.GetHistory(nik, days)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch operator history: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch operator history", history)
}

// Recalculate godoc
// POST /api/prestasi/recalculate
// Recalculates points & ranking engine
func (h *PrestasiHandler) Recalculate(c fiber.Ctx) error {
	daysStr := c.Query("days", "30")
	days, _ := strconv.Atoi(daysStr)

	if err := h.svc.Recalculate(days); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to recalculate points: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Prestasi points recalculated successfully", nil)
}

