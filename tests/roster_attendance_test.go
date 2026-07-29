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

	t.Run("GET /api/rosters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/rosters/:key/export", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/1/export", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/rosters/:key/detail", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/1/detail", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/rosters/upload - missing file returns 400", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/rosters/upload", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 (missing file), got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/rosters/revisions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/revisions", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/rosters/revisions/codes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/revisions/codes", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/rosters/revisions/batch", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"sid": "SID-001", "revisions": []interface{}{}})
		req := httptest.NewRequest("POST", "/api/rosters/revisions/batch", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /api/rosters/revisions/:id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/rosters/revisions/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/rosters/approvals/:id/approve", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/rosters/approvals/1/approve", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /api/rosters/approvals/:id/note", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"note": "Approved with note"})
		req := httptest.NewRequest("PATCH", "/api/rosters/approvals/1/note", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /api/rosters/approvals/:id/reject", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/rosters/approvals/1/reject", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/rosters/attendance", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/rosters/attendance", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/attendance/today", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/attendance/today", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/attendance/date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/attendance/date?date=2026-07-22", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/attendance/range", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/attendance/range?from=2026-07-01&to=2026-07-31", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/attendance/checkin", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264133", "machine": "FP-01"})
		req := httptest.NewRequest("POST", "/api/attendance/checkin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/attendance/checkout", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"nik": "503264133", "machine": "FP-01"})
		req := httptest.NewRequest("POST", "/api/attendance/checkout", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected 201, got %d", resp.StatusCode)
		}
	})
}
