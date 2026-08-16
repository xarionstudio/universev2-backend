package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"universev/internal/config"
	"universev/internal/model"
	"universev/internal/router"
)

// AppSettingsDB matches the DB model used by settings_repo.go
type AppSettingsDB struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement"`
	AppName     string `gorm:"column:app_name"`
	AppEnv      string `gorm:"column:app_env"`
	CompanyLogo string `gorm:"column:company_logo"`
	Theme       string `gorm:"column:theme"`
	Lang        string `gorm:"column:lang"`
	MenuVisJSON string `gorm:"column:menu_vis_json"`
}

func (AppSettingsDB) TableName() string { return "app_settings" }

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Auto-migrate all models (including repo-specific types)
	err = db.AutoMigrate(
		&model.Employee{},
		&model.Competency{},
		&model.FTWRecord{},
		&model.RosterMeta{},
		&model.RosterSchedule{},
		&model.RosterRevision{},
		&model.AttendanceRow{},
		&model.FleetSetting{},
		&model.FleetSettingUnit{},
		&model.FleetAlloc{},
		&model.UnitDb{},
		&model.UnitStatusRow{},
		&model.UnitHistoryRow{},
		// MdEntry removed — master_entries table dropped in migration 000006
		&model.User{},
		&model.Role{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.Notification{},
		&model.AudioSchedule{},
		&model.AudioScheduleDisplay{},
		&model.DisplayDevice{},
		&model.DisplayFleet{},
		&AppSettingsDB{},
		&model.PrestasiScore{},
		&model.PrestasiBadge{},
		&model.PrestasiHistoryEntry{},
		&model.MasterEGIType{},
		&model.MasterProduct{},
		&model.MasterEqClass{},
		&model.MasterArea{},
		&model.MasterTempudo{},
		&model.MasterBus{},
		&model.MasterLocationEx{},
		&model.MasterMess{},
		&model.MasterRunningText{},
		&model.FingerprintDevice{},
		&model.BusinessRule{},
		&model.MasterShiftCode{},
	)
	if err != nil {
		panic("Failed to auto-migrate: " + err.Error())
	}

	// Seed test user (angel@unggul.co.id / admin123)
	now := time.Now()
	nik := "503264133"
	testUser := &model.User{
		Email:        "angel@unggul.co.id",
		Name:         "Angel Test",
		NIK:          &nik,
		PasswordHash: "f2e323ca0d2ae623c55f2659554084c267b6c01c10b5acfe86e8120df80e0d3f",
		PasswordSalt: "testsalt1234567890",
		IsActive:     true,
		Roles:        []string{"1"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_ = db.Create(testUser)

	// Seed user role
	_ = db.Create(&model.UserRole{
		UserID: testUser.ID,
		RoleID: 1,
	})

	// Seed a test role
	_ = db.Create(&model.Role{
		ID:   1,
		Name: "Super Admin",
	})

	// Also seed a second role for update/delete tests
	_ = db.Create(&model.Role{
		ID:   2,
		Name: "Operator",
	})

	// Seed role permissions for r1
	modules := []string{"dashboard", "display", "employees", "ftw", "roster", "asset", "master", "users", "settings", "prestasi"}
	for _, m := range modules {
		_ = db.Create(&model.RolePermission{
			RoleID:          1,
			ModuleName:      m,
			PermissionLevel: "manage",
		})
	}

	// Seed test employee
	_ = db.Create(&model.Employee{
		NIK: "503264133", Name: "Angel Test", Dept: "Operation", Pos: "Operator",
		Status: "aktif", Company: "PT Unggul",
	})

	// Seed test roster meta & schedules
	meta := &model.RosterMeta{
		Label: "July 2026", Month: "2026-07", Dept: "Operation",
		File: "roster_jul.xlsx", Emp: 1, Rows: "1", By: "System",
		DateISO: "2026-07-01", Status: "aktif",
	}
	_ = db.Create(meta)

	_ = db.Create(&model.RosterSchedule{
		RosterFileID: meta.ID, EmployeeNIK: "503264133", ScheduleDate: "2026-07-15", ShiftCode: "D",
	})
	_ = db.Create(&model.RosterSchedule{
		RosterFileID: meta.ID, EmployeeNIK: "503264133", ScheduleDate: "2026-07-16", ShiftCode: "D",
	})
	_ = db.Create(&model.RosterSchedule{
		RosterFileID: meta.ID, EmployeeNIK: "503264133", ScheduleDate: "2026-07-17", ShiftCode: "D",
	})
	_ = db.Create(&model.RosterSchedule{
		RosterFileID: meta.ID, EmployeeNIK: "503264133", ScheduleDate: "2026-07-18", ShiftCode: "D",
	})

	// Seed test roster revisions
	_ = db.Create(&model.RosterRevision{
		ID: 1, SubmissionID: "SID-001", NIK: "503264133", Name: "Angel Test",
		WhatId: "Day Shift", WhatEn: "Day Shift", WhenId: "2026-07-15", WhenEn: "2026-07-15",
		Status: "pending", TargetDate: "2026-07-15",
	})
	_ = db.Create(&model.RosterRevision{
		ID: 2, SubmissionID: "SID-002", NIK: "503264133", Name: "Angel Test",
		WhatId: "Night Shift", WhatEn: "Night Shift", WhenId: "2026-07-16", WhenEn: "2026-07-16",
		Status: "pending", TargetDate: "2026-07-16",
	})
	_ = db.Create(&model.RosterRevision{
		ID: 3, SubmissionID: "SID-003", NIK: "503264133", Name: "Angel Test",
		WhatId: "Night Shift", WhatEn: "Night Shift", WhenId: "2026-07-17", WhenEn: "2026-07-17",
		Status: "pending", TargetDate: "2026-07-17",
	})
	_ = db.Create(&model.RosterRevision{
		ID: 4, SubmissionID: "SID-004", NIK: "503264133", Name: "Angel Test",
		WhatId: "Off", WhatEn: "Off", WhenId: "2026-07-18", WhenEn: "2026-07-18",
		Status: "pending", TargetDate: "2026-07-18",
	})

	// Seed app settings (omit id so it uses auto-increment)
	_ = db.Exec("INSERT INTO app_settings (app_name) VALUES ('UniverseV2')")

	return db
}

// setupTestApp initializes Fiber app with test DB and returns valid JWT token
func setupTestApp() (*fiber.App, string) {
	cfg := &config.Config{
		AppName:            "UniverseV2 API Test",
		AppEnv:             "test",
		AppPort:            "8080",
		JWTSecret:          "testsecretkey12345678901234567890",
		JWTExpiration:      1 * 3600 * 1000000000, // 1 hour
		CORSAllowedOrigins: "*",
	}

	db := setupTestDB()

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": cfg.AppName,
			"env":     cfg.AppEnv,
		})
	})

	router.SetupRoutes(app, cfg, db, nil)

	// Obtain JWT Token for protected endpoints
	loginReq := map[string]string{
		"email":    "angel@unggul.co.id",
		"password": "admin123",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return app, ""
	}

	respBody, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(respBody, &loginResp)

	return app, loginResp.Data.Token
}

func TestHealthCheck(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
