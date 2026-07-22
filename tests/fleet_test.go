package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFleetAndUnitsEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/fleet/settings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/fleet/settings", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/fleet/settings", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"digger": "EX5002", "loc": "Panel East"})
		req := httptest.NewRequest("POST", "/api/v1/fleet/settings", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/fleet/settings/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"loc": "Panel West"})
		req := httptest.NewRequest("PUT", "/api/v1/fleet/settings/fl-EX5002", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/fleet/settings/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/fleet/settings/fl-EX5002", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/fleet/allocations", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/fleet/allocations", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/fleet/allocations/auto", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/fleet/allocations/auto", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/units/status", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/units/status", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/units/:code/status", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"status": "ready"})
		req := httptest.NewRequest("PUT", "/api/v1/units/DT5108/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/units/:code/status-report", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"reason": "Hose hidrolik bocor"})
		req := httptest.NewRequest("POST", "/api/v1/units/DT5108/status-report", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/units/:code/history", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/units/DT5108/history", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/units/db", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/units/db", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/units/db", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"code": "DT9999", "type": "DT777"})
		req := httptest.NewRequest("POST", "/api/v1/units/db", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/units/db", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"code": "DT9999", "type": "DT777"})
		req := httptest.NewRequest("PUT", "/api/v1/units/db", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/units/db", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/units/db", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
