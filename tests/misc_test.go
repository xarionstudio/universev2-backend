package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
)

func TestMiscEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/weather/current", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/weather/current", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req, fiber.TestConfig{Timeout: 5 * time.Second})
		if err != nil {
			// Network timeout to external Open-Meteo API is acceptable in offline test environment
			t.Logf("Weather external API request timed out or failed (offline env): %v", err)
		} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadGateway && resp.StatusCode != http.StatusGatewayTimeout {
			t.Errorf("Expected 200, 502, or 504, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/prestasi/leaderboard", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/prestasi/leaderboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/prestasi/:nik/history", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/prestasi/503264133/history", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
