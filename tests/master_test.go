package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMasterDataEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/master/dept", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/master/dept", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/master/pos", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/master/pos", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/master/dept", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "New Department"})
		req := httptest.NewRequest("POST", "/api/v1/master/dept", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/master/dept/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated Department"})
		req := httptest.NewRequest("PUT", "/api/v1/master/dept/dept-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/master/dept/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/master/dept/dept-1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
