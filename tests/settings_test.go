package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotificationsAndSettings(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/notifications", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/notifications/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/notifications/:id/read", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/notifications/n1/read", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/notifications/read-all", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/notifications/read-all", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/settings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/settings/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/settings", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"theme": "dark"})
		req := httptest.NewRequest("PUT", "/api/v1/settings/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/settings/audio", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/settings/audio", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/settings/displays", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/settings/displays", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/settings/displays", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "TV Front Office", "loc": "Lobby"})
		req := httptest.NewRequest("POST", "/api/v1/settings/displays", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/settings/displays/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "TV Gate Selatan", "loc": "Gate Selatan"})
		req := httptest.NewRequest("PUT", "/api/v1/settings/displays/DSP-A01", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/settings/displays/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/settings/displays/DSP-A01", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/displays/:id/heartbeat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/displays/DSP-A01/heartbeat", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/fingerprint/devices", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/fingerprint/devices", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/settings/audio", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"title": "Morning Bell", "when": "07:00"})
		req := httptest.NewRequest("POST", "/api/v1/settings/audio", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/settings/audio/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"title": "Updated Bell", "when": "08:00"})
		req := httptest.NewRequest("PUT", "/api/v1/settings/audio/audio-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/settings/audio/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/settings/audio/audio-1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
