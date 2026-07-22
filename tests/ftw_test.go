package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFTWEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/ftw/today", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ftw/today", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/ftw/history", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/ftw/history", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/ftw/submit", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"sleepMin": 440, "shift": "siang"})
		req := httptest.NewRequest("POST", "/api/v1/ftw/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})
}
