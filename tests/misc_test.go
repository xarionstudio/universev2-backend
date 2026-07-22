package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiscEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/weather/current", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/weather/current", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/prestasi/leaderboard", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/prestasi/leaderboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
