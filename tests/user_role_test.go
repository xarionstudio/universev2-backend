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

	t.Run("GET /api/users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/users", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"email":    "new.user@unggul.co.id",
			"name":     "New User",
			"password": "password1",
			"nik":      "503264111",
			"roles":    []string{"1"},
		})
		req := httptest.NewRequest("POST", "/api/users/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/users/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"name":  "Updated User",
			"email": "angel@unggul.co.id",
			"roles": []string{"1"},
		})
		req := httptest.NewRequest("PUT", "/api/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/users/:id - non-existent user returns 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/users/999", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/roles", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/roles/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/roles", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "New Role"})
		req := httptest.NewRequest("POST", "/api/roles/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/roles/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Role"})
		req := httptest.NewRequest("PUT", "/api/roles/2", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/roles/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/roles/2", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
