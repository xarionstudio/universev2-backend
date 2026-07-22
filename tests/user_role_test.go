package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAndRoleEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/users", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "new.user@unggul.co.id", "name": "New User"})
		req := httptest.NewRequest("POST", "/api/v1/users/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/users/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated User"})
		req := httptest.NewRequest("PUT", "/api/v1/users/u1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/users/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/users/u1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/roles", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/roles/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/roles", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "New Role"})
		req := httptest.NewRequest("POST", "/api/v1/roles/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/roles/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Role"})
		req := httptest.NewRequest("PUT", "/api/v1/roles/r2", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/roles/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/roles/r2", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
