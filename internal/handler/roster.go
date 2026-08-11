package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"

	"universev/internal/dto"
	"universev/internal/export"
	"universev/internal/model"
	internalpkg "universev/internal/pkg"
	"universev/internal/repository"
	"universev/pkg/filter"
	"universev/pkg/pagination"
	"universev/pkg/response"
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
		Emp:     0,
		Rows:    "0",
		By:      createdBy,
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

	empMap, _ := h.repo.GetEmployeeNIKMap()

	// Parse roster schedules and validation from the file
	schedules, valResult, parseErr := export.ParseRosterWithValidation(f, month, empMap)
	if parseErr != nil {
		// File saved but parsing failed — return warning
		return response.Success(c, fiber.StatusCreated, "Roster uploaded but parsing failed: "+parseErr.Error(), fiber.Map{
			"meta":       meta,
			"validation": valResult,
		})
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
	meta.Emp = valResult.ValidCount
	meta.Rows = fmt.Sprintf("%d", len(valResult.Preview))
	_ = h.repo.UpdateRosterMeta(meta.ID, meta)

	return response.Success(c, fiber.StatusCreated, "Roster uploaded and processed successfully", fiber.Map{
		"meta":       meta,
		"validation": valResult,
	})
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

func (h *RosterHandler) GetShiftCodes(c fiber.Ctx) error {
	// Try to fetch from master_shift_codes table
	if masterCodes, err := h.repo.GetMasterByCategory("shift_codes"); err == nil && len(masterCodes) > 0 {
		// Convert master data to expected format
		groups := []fiber.Map{}
		groupMap := make(map[string][]fiber.Map)

		for _, code := range masterCodes {
			group := code.GroupID
			if group == "" {
				group = "Other"
			}
			groupMap[group] = append(groupMap[group], fiber.Map{
				"k":   code.Code,
				"v":   code.Name,
				"vEn": code.NameEn,
			})
		}

		for group, codes := range groupMap {
			groups = append(groups, fiber.Map{
				"group":   group,
				"groupEn": group,
				"codes":   codes,
			})
		}

		return response.Success(c, fiber.StatusOK, "Success fetch shift codes", groups)
	}

	// Fallback to hardcoded defaults if master data not available
	groups := []fiber.Map{
		{
			"group":   "Shift & kehadiran",
			"groupEn": "Shifts & attendance",
			"codes": []fiber.Map{
				{"k": "D", "v": "Day shift", "vEn": "Day shift"},
				{"k": "N", "v": "Night shift", "vEn": "Night shift"},
				{"k": "R", "v": "Reguler", "vEn": "Regular"},
				{"k": "STB", "v": "Standby", "vEn": "Standby"},
				{"k": "OFF", "v": "OFF", "vEn": "OFF"},
			},
		},
		{
			"group":   "Cuti & izin",
			"groupEn": "Leave & permits",
			"codes": []fiber.Map{
				{"k": "CR", "v": "Cuti roster", "vEn": "Roster leave"},
				{"k": "AL", "v": "Annual leave", "vEn": "Annual leave"},
				{"k": "LWP", "v": "Izin dengan upah", "vEn": "Paid leave"},
				{"k": "LWOP", "v": "Izin tanpa upah", "vEn": "Unpaid leave"},
				{"k": "PH", "v": "Public holiday", "vEn": "Public holiday"},
				{"k": "PHD", "v": "Public holiday siang", "vEn": "Public holiday (day)"},
			},
		},
		{
			"group":   "Sakit & ketidakhadiran",
			"groupEn": "Sickness & absence",
			"codes": []fiber.Map{
				{"k": "S", "v": "Sakit", "vEn": "Sick"},
				{"k": "A", "v": "Alpha", "vEn": "Alpha / no notice"},
			},
		},
		{
			"group":   "Medis & karantina",
			"groupEn": "Medical & quarantine",
			"codes": []fiber.Map{
				{"k": "MCU", "v": "Medical check up", "vEn": "Medical check up"},
				{"k": "MCR", "v": "Reguler MCU", "vEn": "Regular MCU"},
				{"k": "MCUF", "v": "Follow up MCU", "vEn": "MCU follow-up"},
				{"k": "ISM", "v": "Isolasi mandiri", "vEn": "Self-isolation"},
				{"k": "OBC", "v": "Observasi COVID", "vEn": "COVID observation"},
				{"k": "KRT", "v": "Karantina", "vEn": "Quarantine"},
			},
		},
		{
			"group":   "Tugas & training",
			"groupEn": "Assignment & training",
			"codes": []fiber.Map{
				{"k": "TGS", "v": "Tugas", "vEn": "Assignment"},
				{"k": "DNS", "v": "Dinas", "vEn": "Official duty"},
				{"k": "TRV", "v": "Travel", "vEn": "Travel"},
				{"k": "TR", "v": "Training di luar site", "vEn": "Off-site training"},
				{"k": "TRS", "v": "Training onsite", "vEn": "On-site training"},
				{"k": "IN", "v": "Induksi", "vEn": "Induction"},
			},
		},
		{
			"group":   "Status kepegawaian",
			"groupEn": "Employment status",
			"codes": []fiber.Map{
				{"k": "TERM", "v": "Termination", "vEn": "Termination"},
				{"k": "EOC", "v": "Kontrak berakhir", "vEn": "Contract ended"},
				{"k": "RSG", "v": "Resign", "vEn": "Resign"},
			},
		},
	}
	return response.Success(c, fiber.StatusOK, "Success fetch shift codes", groups)
}

func (h *RosterHandler) GetRevisionCodes(c fiber.Ctx) error {
	codes := []fiber.Map{
		{"id": "D", "label": "Day Shift"},
		{"id": "N", "label": "Night Shift"},
		{"id": "R", "label": "Reguler"},
		{"id": "STB", "label": "Standby"},
		{"id": "OFF", "label": "Off / Libur"},
		{"id": "CR", "label": "Cuti Roster"},
		{"id": "AL", "label": "Annual Leave"},
		{"id": "S", "label": "Sakit"},
		{"id": "A", "label": "Alpha / Tanpa Keterangan"},
		{"id": "LWP", "label": "Izin Dengan Upah"},
		{"id": "LWOP", "label": "Izin Tanpa Upah"},
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

	// Use the authenticated user's name as approver
	approver := "System"
	if u := c.Locals("user"); u != nil {
		if claims, ok := u.(*internalpkg.JWTCustomClaims); ok && claims != nil {
			// Look up user name from DB using claims.UserID
			if user, err := h.repo.GetUserByID(claims.UserID); err == nil && user != nil && user.Name != "" {
				approver = user.Name
			}
		}
	}
	byId := "Disetujui oleh " + approver
	byEn := "Approved by " + approver

	if err := h.repo.ApproveRevision(id, byId, byEn); err != nil {
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

	// Use the authenticated user's name as approver
	approver := "System"
	if u := c.Locals("user"); u != nil {
		if claims, ok := u.(*internalpkg.JWTCustomClaims); ok && claims != nil {
			// Look up user name from DB using claims.UserID
			if user, err := h.repo.GetUserByID(claims.UserID); err == nil && user != nil && user.Name != "" {
				approver = user.Name
			}
		}
	}

	note := req.Note
	if note == "" {
		note = "Disetujui dengan catatan"
	}
	byId := note + " — " + approver
	byEn := note + " — " + approver
	if err := h.repo.ApproveRevision(id, byId, byEn); err != nil {
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

	// Use the authenticated user's name as approver
	approver := "System"
	if u := c.Locals("user"); u != nil {
		if claims, ok := u.(*internalpkg.JWTCustomClaims); ok && claims != nil {
			// Look up user name from DB using claims.UserID
			if user, err := h.repo.GetUserByID(claims.UserID); err == nil && user != nil && user.Name != "" {
				approver = user.Name
			}
		}
	}
	byId := "Ditolak oleh " + approver
	byEn := "Rejected by " + approver

	if err := h.repo.RejectRevision(id, byId, byEn); err != nil {
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
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	rows, err := h.repo.GetAttendance(date)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch attendance: "+err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Success fetch attendance", rows)
}
