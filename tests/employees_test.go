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

	t.Run("GET /api/v1/employees", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/employees/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/employees/:nik", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/employees/503264133", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/employees", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264999", "name": "Test Employee"})
		req := httptest.NewRequest("POST", "/api/v1/employees/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/employees/:nik", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Test Employee"})
		req := httptest.NewRequest("PUT", "/api/v1/employees/503264999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/employees/:nik", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/employees/503264999", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
