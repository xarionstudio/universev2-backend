package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"universev2-backend/internal/config"
	"universev2-backend/internal/router"
)

// setupTestApp initializes Fiber app and returns valid JWT token
func setupTestApp() (*fiber.App, string) {
	cfg := &config.Config{
		AppName:            "UniverseV2 API Test",
		AppEnv:             "test",
		AppPort:            "8080",
		JWTSecret:          "testsecretkey12345678901234567890",
		JWTExpiration:      1 * 3600 * 1000000000, // 1 hour
		CORSAllowedOrigins: "*",
	}

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

	router.SetupRoutes(app, cfg, nil)

	// Obtain JWT Token for protected endpoints
	loginReq := map[string]string{
		"email":    "angel@unggul.co.id",
		"password": "admin123",
	}
	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
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
