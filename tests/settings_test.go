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

	t.Run("GET /api/notifications", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/notifications/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/notifications/:id/read", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/notifications/1/read", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/notifications/read-all", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/notifications/read-all", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/settings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/settings/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/settings", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"theme": "dark"})
		req := httptest.NewRequest("PUT", "/api/settings/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/settings/audio", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/settings/audio", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/settings/displays", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/settings/displays", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/settings/displays", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "TV Front Office", "loc": "Lobby"})
		req := httptest.NewRequest("POST", "/api/settings/displays", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/settings/displays/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "TV Gate Selatan", "loc": "Gate Selatan"})
		req := httptest.NewRequest("PUT", "/api/settings/displays/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/settings/displays/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/settings/displays/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/displays/:id/heartbeat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/displays/1/heartbeat", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/fingerprint/devices", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/fingerprint/devices", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/settings/audio", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"title": "Morning Bell", "when": "07:00"})
		req := httptest.NewRequest("POST", "/api/settings/audio", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/settings/audio/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"title": "Updated Bell", "when": "08:00"})
		req := httptest.NewRequest("PUT", "/api/settings/audio/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/settings/audio/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/settings/audio/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
