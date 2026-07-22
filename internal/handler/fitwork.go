package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type FitworkHandler struct {
	repo *repository.FTWRepo
}

func NewFitworkHandler(repo *repository.FTWRepo) *FitworkHandler {
	return &FitworkHandler{repo: repo}
}

func (h *FitworkHandler) GetTodayLog(c fiber.Ctx) error {
	date := c.Query("date", time.Now().Format("2006-01-02"))

	logs, err := h.repo.GetTodayLogs(date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch FTW logs: "+err.Error())
	}

	meta := &response.Meta{Page: 1, Limit: 50, Total: len(logs), TotalPage: 1}
	return response.SuccessWithMeta(c, fiber.StatusOK, "Success fetch FTW logs", logs, meta)
}

func (h *FitworkHandler) SubmitLog(c fiber.Ctx) error {
	var req struct {
		NIK      string `json:"nik"`
		Shift    string `json:"shift"`
		SleepMin *int   `json:"sleepMin"`
		Sleep    string `json:"sleep"`
		SendTime string `json:"sendTime"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}
	if isTrimmedEmpty(req.Shift) {
		return sendValidationError(c, "shift", "Shift is required")
	}

	rec := &model.FTWRecord{
		NIK: req.NIK, Shift: req.Shift, SleepMin: req.SleepMin,
		Sleep: req.Sleep, SendTime: req.SendTime,
		Date: time.Now().Format("2006-01-02"),
	}
	if err := h.repo.Submit(rec); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to submit FTW log: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "FTW log submitted", nil)
}

func (h *FitworkHandler) GetHistory(c fiber.Ctx) error {
	nik := c.Query("nik")
	if isTrimmedEmpty(nik) {
		return response.Error(c, fiber.StatusBadRequest, "Query parameter 'nik' is required")
	}

	logs, err := h.repo.GetHistory(nik, 30)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch history: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch history", logs)
}
