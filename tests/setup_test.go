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

	"universev2-backend/internal/config"
	"universev2-backend/internal/model"
	"universev2-backend/internal/router"
)

// AppSettingsDB matches the DB model used by settings_repo.go
type AppSettingsDB struct {
	ID          string `gorm:"column:id;primaryKey"`
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
		&model.MdEntry{},
		&model.User{},
		&model.Role{},
		&model.RolePermission{},
		&model.UserRole{},
		&model.Notification{},
		&model.AudioSchedule{},
		&model.AudioScheduleDisplay{},
		&model.DisplayDevice{},
		&AppSettingsDB{},
		&model.PrestasiScore{},
		&model.PrestasiBadge{},
		&model.PrestasiHistoryEntry{},
	)
	if err != nil {
		panic("Failed to auto-migrate: " + err.Error())
	}

	// Seed test user (angel@unggul.co.id / admin123)
	now := time.Now()
	nik := "503264133"
	testUser := &model.User{
		ID:           "u-test-1",
		Email:        "angel@unggul.co.id",
		Name:         "Angel Test",
		NIK:          &nik,
		PasswordHash: "f2e323ca0d2ae623c55f2659554084c267b6c01c10b5acfe86e8120df80e0d3f",
		PasswordSalt: "testsalt1234567890",
		IsActive:     true,
		Roles:        []string{"r1"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	_ = db.Create(testUser)

	// Seed user role
	_ = db.Create(&model.UserRole{
		UserID: "u-test-1",
		RoleID: "r1",
	})

	// Seed a test role
	_ = db.Create(&model.Role{
		ID:   "r1",
		Name: "Super Admin",
	})

	// Also seed a second role for update/delete tests
	_ = db.Create(&model.Role{
		ID:   "r2",
		Name: "Operator",
	})

	// Seed role permissions for r1
	modules := []string{"employees", "ftw", "roster", "asset", "master", "users", "settings", "prestasi"}
	for _, m := range modules {
		_ = db.Create(&model.RolePermission{
			RoleID:          "r1",
			ModuleName:      m,
			PermissionLevel: "manage",
		})
	}

	// Seed test employee
	_ = db.Create(&model.Employee{
		NIK: "503264133", Name: "Angel Test", Dept: "Operation", Pos: "Operator",
		Status: "active", Company: "PT Unggul",
	})

	// Seed test roster meta
	_ = db.Create(&model.RosterMeta{
		ID: "jul", Label: "July 2026", Month: "2026-07", Dept: "Operation",
		File: "roster_jul.xlsx", Emp: "1", Rows: "1", By: "System",
		Date: "2026-07-01", DateISO: "2026-07-01", Status: "aktif",
	})

	// Seed app settings
	_ = db.Exec("INSERT INTO app_settings (id, app_name) VALUES ('default', 'UniverseV2')")

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

	router.SetupRoutes(app, cfg, db)

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
