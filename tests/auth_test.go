package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("POST /api/auth/login", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "angel@unggul.co.id", "password": "admin123"})
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/auth/register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    "newuser@unggul.co.id",
			"name":     "New User",
			"nik":      "503264111",
			"password": "password1",
			"dept":     "Operation",
			"pos":      "Operator",
		})
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/auth/refresh", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/auth/logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

func TestProfileEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/profile", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/profile/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/profile", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"name":  "Updated Name",
			"email": "angel@unggul.co.id",
		})
		req := httptest.NewRequest("PUT", "/api/profile/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/profile/password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"oldPassword":     "admin123",
			"newPassword":     "newpassword1",
			"confirmPassword": "newpassword1",
		})
		req := httptest.NewRequest("PUT", "/api/profile/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
