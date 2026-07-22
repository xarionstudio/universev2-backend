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

	t.Run("POST /api/v1/auth/login", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "angel@unggul.co.id", "password": "admin123"})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/auth/register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": "newuser@unggul.co.id", "name": "New User"})
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/auth/refresh", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/refresh", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/auth/logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}

func TestProfileEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/profile", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/profile/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/profile", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Name"})
		req := httptest.NewRequest("PUT", "/api/v1/profile/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/profile/password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"oldPassword": "admin123", "newPassword": "newpassword123"})
		req := httptest.NewRequest("PUT", "/api/v1/profile/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
