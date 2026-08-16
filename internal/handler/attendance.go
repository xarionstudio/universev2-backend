package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	"universev/internal/repository"
	"universev/pkg/pagination"
	"universev/pkg/response"
)

type AttendanceHandler struct {
	repo   *repository.AttendanceRepo
	fpRepo *repository.FingerprintRepo
}

func NewAttendanceHandler(repo *repository.AttendanceRepo, fpRepo *repository.FingerprintRepo) *AttendanceHandler {
	return &AttendanceHandler{repo: repo, fpRepo: fpRepo}
}

func (h *AttendanceHandler) resolveDefaultMachine() string {
	if h.fpRepo != nil {
		devices, err := h.fpRepo.GetActiveDevices()
		if err == nil && len(devices) > 0 && devices[0].Code != "" {
			return devices[0].Code
		}
	}
	return "WEB_MANUAL"
}

// GetAttendanceToday — GET /api/attendance/today
func (h *AttendanceHandler) GetAttendanceToday(c fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")
	_ = h.repo.SyncAttendanceBoard(today)
	rows, err := h.repo.GetLogsByDate(today)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch today attendance: "+err.Error())
	}

	return response.SuccessPaged(c, fiber.StatusOK, "Success fetch today attendance", response.PagedData{
		Items: rows,
		Pagination: pagination.Meta{
			Page:       1,
			PerPage:    50,
			Total:      int64(len(rows)),
			TotalPages: 1,
		},
	})
}

// GetAttendanceByDate — GET /api/attendance/date?date=YYYY-MM-DD
func (h *AttendanceHandler) GetAttendanceByDate(c fiber.Ctx) error {
	date := c.Query("date", time.Now().Format("2006-01-02"))
	_ = h.repo.SyncAttendanceBoard(date)
	rows, err := h.repo.GetLogsByDate(date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance by date: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Success fetch attendance by date", rows)
}

// GetAttendanceRange — GET /api/attendance/range?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *AttendanceHandler) GetAttendanceRange(c fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	if isTrimmedEmpty(from) || isTrimmedEmpty(to) {
		return sendValidationError(c, "from/to", "Date range 'from' and 'to' are required")
	}

	_ = h.repo.SyncAttendanceRange(from, to)
	rows, err := h.repo.GetLogsRange(from, to)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance range: "+err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Success fetch attendance range", rows)
}

// RecordCheckIn — POST /api/attendance/checkin
func (h *AttendanceHandler) RecordCheckIn(c fiber.Ctx) error {
	var req dto.CheckInOutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}

	if isTrimmedEmpty(req.Machine) {
		req.Machine = h.resolveDefaultMachine()
	}

	row, err := h.repo.RecordCheckIn(req.NIK, req.Machine)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to record check-in: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Check-in recorded successfully", row)
}

// RecordCheckOut — POST /api/attendance/checkout
func (h *AttendanceHandler) RecordCheckOut(c fiber.Ctx) error {
	var req dto.CheckInOutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(req.NIK) {
		return sendValidationError(c, "nik", "NIK is required")
	}

	if isTrimmedEmpty(req.Machine) {
		req.Machine = h.resolveDefaultMachine()
	}

	row, err := h.repo.RecordCheckOut(req.NIK, req.Machine)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to record check-out: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Check-out recorded successfully", row)
}
