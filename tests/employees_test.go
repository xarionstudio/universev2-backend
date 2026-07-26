package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmployeesEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/employees", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/employees/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/employees/:nik", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/employees/503264133", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/employees", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264999", "name": "Test Employee"})
		req := httptest.NewRequest("POST", "/api/employees/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/employees/:nik", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Test Employee"})
		req := httptest.NewRequest("PUT", "/api/employees/503264999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/employees/:nik", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/employees/503264999", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/employees/:nik/competencies", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/employees/503264133/competencies", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/employees/:nik/competencies", func(t *testing.T) {
		body, _ := json.Marshal([]interface{}{
			map[string]interface{}{
				"type":     "simper",
				"category": "A",
				"expiry":   "2026-12-31",
			},
		})
		req := httptest.NewRequest("PUT", "/api/employees/503264133/competencies", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/employees/:nik/photo - missing file returns 400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/employees/503264133/photo", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 (missing file), got %d", resp.StatusCode)
		}
	})
}
