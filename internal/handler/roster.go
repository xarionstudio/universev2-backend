package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/export"
	"universev2-backend/internal/model"
	"universev2-backend/internal/repository"
	"universev2-backend/pkg/response"
)

type RosterHandler struct {
	repo *repository.RosterRepo
}

func NewRosterHandler(repo *repository.RosterRepo) *RosterHandler {
	return &RosterHandler{repo: repo}
}

func (h *RosterHandler) GetRosters(c fiber.Ctx) error {
	dept := c.Query("dept")
	metas, err := h.repo.GetRosters(dept)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch rosters: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch rosters", metas)
}

func (h *RosterHandler) UploadRoster(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Roster file is required")
	}

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(file.Filename), ".xls") &&
		!strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		return response.Error(c, fiber.StatusBadRequest, "Only .xlsx, .xls, and .csv files are accepted")
	}

	// Save file to uploads directory
	uploadDir := "uploads/rosters"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create upload directory")
	}

	savePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to save roster file: "+err.Error())
	}

	// Parse month and dept from form fields
	month := c.FormValue("month", time.Now().Format("2006-01"))
	dept := c.FormValue("dept", "Operation")
	label := c.FormValue("label", month)
	createdBy := c.FormValue("createdBy", "System")

	// Generate a unique key for the roster
	key := fmt.Sprintf("%s-%s", strings.ReplaceAll(month, "-", ""), strings.ToLower(dept))

	// Create roster meta entry
	meta := &model.RosterMeta{
		ID:      key,
		Label:   label,
		Month:   month,
		Dept:    dept,
		File:    file.Filename,
		Emp:     "0",
		Rows:    "0",
		By:      createdBy,
		Date:    time.Now().Format("2006-01-02"),
		DateISO: time.Now().Format("2006-01-02"),
		Status:  "aktif",
	}

	if err := h.repo.CreateRoster(meta); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to save roster metadata: "+err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Roster uploaded and processed successfully", meta)
}

func (h *RosterHandler) ExportRoster(c fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		key = c.Query("key")
	}

	month := c.Query("month")
	dept := c.Query("dept")

	if key != "" {
		meta, err := h.repo.GetRosterByID(key)
		if err == nil && meta != nil {
			if month == "" {
				month = meta.Month
			}
			if dept == "" {
				dept = meta.Dept
			}
		}
	}

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	exportRows, err := h.repo.GetExportRosterData(key, dept)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch roster data for export: "+err.Error())
	}

	xlsxData, err := export.GenerateRosterExcel(key, month, dept, exportRows)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to generate roster excel: "+err.Error())
	}

	fileName := fmt.Sprintf("roster_%s_export.xlsx", month)
	if key != "" {
		fileName = fmt.Sprintf("roster_%s_export.xlsx", key)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	return c.Send(xlsxData)
}

func (h *RosterHandler) GetRevisions(c fiber.Ctx) error {
	status := c.Query("status")
	revisions, err := h.repo.GetRevisions(status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch revisions: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch revisions", revisions)
}

func (h *RosterHandler) GetRevisionCodes(c fiber.Ctx) error {
	codes := []fiber.Map{
		{"id": "D", "label": "Day Shift"},
		{"id": "N", "label": "Night Shift"},
		{"id": "OFF", "label": "Off / Libur"},
		{"id": "C", "label": "Cuti Tahunan"},
		{"id": "S", "label": "Sakit"},
	}
	return response.Success(c, fiber.StatusOK, "Success fetch revision codes", codes)
}

func (h *RosterHandler) SubmitBatchRevision(c fiber.Ctx) error {
	var rev model.RosterRevision
	if err := c.Bind().JSON(&rev); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if isTrimmedEmpty(rev.SubmissionID) {
		return sendValidationError(c, "sid", "Submission ID is required")
	}

	rev.Status = "pending"
	if err := h.repo.CreateRevision(&rev); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to submit revision: "+err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Batch revisions submitted successfully", rev)
}

func (h *RosterHandler) DeleteRevision(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid revision ID")
	}

	if err := h.repo.DeleteRevision(id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete revision: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Revision deleted successfully", nil)
}

func (h *RosterHandler) ApproveRevision(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid revision ID")
	}

	if err := h.repo.ApproveRevision(id, "Disetujui oleh Supervisor", "Approved by Supervisor"); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to approve revision: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Revision approved", nil)
}

func (h *RosterHandler) ApproveRevisionWithNote(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid revision ID")
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	note := req.Note
	if note == "" {
		note = "Disetujui dengan catatan"
	}
	if err := h.repo.ApproveRevision(id, note, note); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to approve revision: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Revision approved with specific note", nil)
}

func (h *RosterHandler) RejectRevision(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid revision ID")
	}

	if err := h.repo.RejectRevision(id, "Ditolak", "Rejected"); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to reject revision: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Revision rejected", nil)
}

func (h *RosterHandler) GetRosterDetail(c fiber.Ctx) error {
	key := c.Params("key")
	if key == "" {
		return response.Error(c, fiber.StatusBadRequest, "Roster key is required")
	}

	detail, err := h.repo.GetRosterDetail(key)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch roster detail: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch roster detail", detail)
}

func (h *RosterHandler) GetAttendance(c fiber.Ctx) error {
	date := c.Query("date")
	rows, err := h.repo.GetAttendance(date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch attendance", rows)
}
