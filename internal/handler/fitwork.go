package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/dto"
	"universev2-backend/internal/export"
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/filter"
	"universev2-backend/pkg/pagination"
	"universev2-backend/pkg/response"
)

type FitworkHandler struct {
	repo *repository.FTWRepo
}

func NewFitworkHandler(repo *repository.FTWRepo) *FitworkHandler {
	return &FitworkHandler{repo: repo}
}

// GetTodayLog godoc
// GET /api/ftw/today
// Query params: date, page, perPage, search, status, dept, nik, date_from, date_to, logic
func (h *FitworkHandler) GetTodayLog(c fiber.Ctx) error {
	// If explicit date_from/date_to or search/status filters are set, use paginated query
	f := filter.ParseFromCtx(c)

	// Default date: today, applied as DateFrom == DateTo when no range given
	if f.DateFrom == "" && f.DateTo == "" {
		date := c.Query("date", time.Now().Format("2006-01-02"))
		f.DateFrom = date
		f.DateTo = date
	}

	p := pagination.Parse(c.Query("page"), c.Query("perPage"))

	logs, total, err := h.repo.GetLogsPaginated(f, p)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch FTW logs: "+err.Error())
	}

	meta := pagination.BuildMeta(p, total)
	return response.SuccessPaged(c, fiber.StatusOK, "Success fetch FTW logs", response.PagedData{
		Items:      logs,
		Pagination: meta,
	})
}

// SubmitLog godoc
// POST /api/ftw/submit
func (h *FitworkHandler) SubmitLog(c fiber.Ctx) error {
	var req dto.SubmitFTWLogRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}
	if isTrimmedEmpty(req.Shift) {
		return sendValidationError(c, "shift", "Shift is required")
	}

	eval := model.EvaluateFTW(req.SleepMin)
	rec := &model.FTWRecord{
		NIK: req.NIK, Shift: req.Shift, SleepMin: req.SleepMin,
		Sleep: req.Sleep, SendTime: req.SendTime,
		Date:      time.Now().Format("2006-01-02"),
		St:        eval.Status, RestHours: eval.RestHours, CanWork: eval.CanWork,
	}
	if err := h.repo.Submit(rec); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to submit FTW log: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "FTW log submitted", rec)
}

// GetHistory godoc
// GET /api/ftw/history?nik=xxx
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

// ExportFTW godoc
// GET /api/ftw/export
// Accepts same filter params as GetTodayLog; returns an xlsx download.
func (h *FitworkHandler) ExportFTW(c fiber.Ctx) error {
	f := filter.ParseFromCtx(c)

	// Default to today if no date range provided
	if f.DateFrom == "" && f.DateTo == "" {
		today := c.Query("date", time.Now().Format("2006-01-02"))
		f.DateFrom = today
		f.DateTo = today
	}

	logs, err := h.repo.GetAllFiltered(f)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch FTW data: "+err.Error())
	}

	xlsxData, err := export.GenerateFTWExcel(logs)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate export: "+err.Error())
	}

	fileName := fmt.Sprintf("ftw_export_%s.xlsx", time.Now().Format("20060102"))
	if f.DateFrom != "" {
		fileName = fmt.Sprintf("ftw_export_%s.xlsx", f.DateFrom)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	return c.Send(xlsxData)
}
