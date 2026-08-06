package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

type RosterHandler struct {
	repo      *repository.RosterRepo
	uploadDir string
}

func NewRosterHandler(repo *repository.RosterRepo, uploadDir string) *RosterHandler {
	return &RosterHandler{repo: repo, uploadDir: uploadDir}
}

// GetRosters godoc
// GET /api/rosters
// Query params: page, perPage, dept, status, month, date_from, date_to, logic
func (h *RosterHandler) GetRosters(c fiber.Ctx) error {
	f := filter.ParseFromCtx(c)
	p := pagination.Parse(c.Query("page"), c.Query("perPage"))

	metas, total, err := h.repo.GetRostersPaginated(f, p)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch rosters: "+err.Error())
	}

	meta := pagination.BuildMeta(p, total)
	return response.SuccessPaged(c, fiber.StatusOK, "Success fetch rosters", response.PagedData{
		Items:      metas,
		Pagination: meta,
	})
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
	rosterDir := h.uploadDir + "/rosters"
	if err := os.MkdirAll(rosterDir, 0755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create upload directory")
	}

	savePath := filepath.Join(rosterDir, file.Filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to save roster file: "+err.Error())
	}

	// Parse month and dept from form fields
	month := c.FormValue("month", time.Now().Format("2006-01"))
	dept := c.FormValue("dept", "Operation")
	label := c.FormValue("label", month)
	createdBy := c.FormValue("createdBy", "System")

	// Create roster meta entry (ID auto-increment from DB)
	meta := &model.RosterMeta{
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

	// Parse the roster file and populate roster_schedules
	// Read file content for parsing
	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to open roster file")
	}
	defer f.Close()

	// Parse roster schedules from the file
	schedules, parseErr := export.ParseRosterExcel(f, month)
	if parseErr != nil {
		// File saved but parsing failed — return warning
		return response.Success(c, fiber.StatusCreated, "Roster uploaded but parsing failed: "+parseErr.Error(), meta)
	}

	// Save schedules to DB with roster_file_id
	savedCount := 0
	for _, sched := range schedules {
		sched.RosterFileID = meta.ID
		if err := h.repo.CreateSchedule(&sched); err == nil {
			savedCount++
		}
	}

	// Update meta with counts
	meta.Emp = fmt.Sprintf("%d", savedCount)
	meta.Rows = fmt.Sprintf("%d", len(schedules))
	_ = h.repo.UpdateRosterMeta(meta.ID, meta)

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

// GetRevisions godoc
// GET /api/rosters/revisions
// Query params: page, perPage, status, nik, search, date_from, date_to, logic
func (h *RosterHandler) GetRevisions(c fiber.Ctx) error {
	f := filter.ParseFromCtx(c)
	p := pagination.Parse(c.Query("page"), c.Query("perPage"))

	revisions, total, err := h.repo.GetRevisionsPaginated(f, p)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch revisions: "+err.Error())
	}

	meta := pagination.BuildMeta(p, total)
	return response.SuccessPaged(c, fiber.StatusOK, "Success fetch revisions", response.PagedData{
		Items:      revisions,
		Pagination: meta,
	})
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
	var req struct {
		Revisions []model.RosterRevision `json:"revisions"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	created := make([]model.RosterRevision, 0, len(req.Revisions))

	// Case 1: Batch format { revisions: [...] }
	if len(req.Revisions) > 0 {
		for _, rev := range req.Revisions {
			if isTrimmedEmpty(rev.SubmissionID) {
				return sendValidationError(c, "sid", "Submission ID is required")
			}
			rev.Status = "pending"
			if err := h.repo.CreateRevision(&rev); err != nil {
				return response.Error(c, fiber.StatusInternalServerError, "Failed to submit revision: "+err.Error())
			}
			created = append(created, rev)
		}
		return response.Success(c, fiber.StatusCreated, "Batch revisions submitted successfully", created)
	}

	// Case 2: Single revision format { sid, nik, whatId, ... } (backward compatible)
	var single model.RosterRevision
	if err := c.Bind().JSON(&single); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if isTrimmedEmpty(single.SubmissionID) {
		return sendValidationError(c, "sid", "Submission ID is required")
	}
	single.Status = "pending"
	if err := h.repo.CreateRevision(&single); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to submit revision: "+err.Error())
	}
	created = append(created, single)
	return response.Success(c, fiber.StatusCreated, "Batch revisions submitted successfully", created)
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

	var req dto.ApproveRevisionNoteRequest
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
