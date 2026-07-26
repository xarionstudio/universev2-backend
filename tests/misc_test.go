package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiscEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/weather/current", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/weather/current", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		if err != nil {
			t.Errorf("Request failed: %v", err)
		} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadGateway {
			t.Errorf("Expected 200 or 502, got %d", resp.StatusCode)
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
