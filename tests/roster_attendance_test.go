package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRostersAndAttendanceEndpoints(t *testing.T) {
	app, token := setupTestApp()

	t.Run("GET /api/v1/rosters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rosters/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/rosters/:key/export", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rosters/jul/export", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/rosters/upload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/rosters/upload", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/rosters/revisions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rosters/revisions", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/rosters/revisions/codes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rosters/revisions/codes", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/rosters/revisions/batch", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"revisions": []interface{}{}})
		req := httptest.NewRequest("POST", "/api/v1/rosters/revisions/batch", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/v1/rosters/revisions/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/rosters/revisions/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/rosters/approvals/:id/approve", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/rosters/approvals/1/approve", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /api/v1/rosters/approvals/:id/note", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"note": "Approved with note"})
		req := httptest.NewRequest("PATCH", "/api/v1/rosters/approvals/1/note", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/v1/rosters/approvals/:id/reject", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/v1/rosters/approvals/1/reject", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/rosters/attendance", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/rosters/attendance", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/attendance/today", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/attendance/today", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/attendance/date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/attendance/date?date=2026-07-22", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/attendance/range", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/attendance/range?from=2026-07-01&to=2026-07-31", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/attendance/checkin", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264133", "machine": "FP-01"})
		req := httptest.NewRequest("POST", "/api/v1/attendance/checkin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/attendance/checkout", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264133", "machine": "FP-01"})
		req := httptest.NewRequest("POST", "/api/v1/attendance/checkout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})
}
