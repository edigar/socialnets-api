package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Run("Should return 200 with an ok status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		HealthCheck(rec, newRequest(http.MethodGet, ""))

		assertStatus(t, rec.Code, http.StatusOK)

		var decoded map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("response should be valid JSON: %v", err)
		}
		if decoded["status"] != "ok" {
			t.Errorf("HealthCheck should return status ok. Got: %v", decoded)
		}
	})
}
