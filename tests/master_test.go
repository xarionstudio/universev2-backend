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

	t.Run("GET /api/master/egi", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/master/egi", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/master/product", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/master/product", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/master/egi", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"code": "NEW-EGI", "name": "New EGI Type"})
		req := httptest.NewRequest("POST", "/api/master/egi", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/master/egi/:id", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated EGI"})
		req := httptest.NewRequest("PUT", "/api/master/egi/DT777", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/master/egi/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/master/egi/NEW-EGI", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})
}
