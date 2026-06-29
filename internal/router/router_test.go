package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	// nil DB is safe: the repository/usecase/controller constructors only store
	// the handle, and none of the routes exercised below reach a DB query.
	r := Generate(nil)

	t.Run("Should respond 200 on the public health route", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("GET /health should return 200, got %d", rec.Code)
		}
	})

	t.Run("Should respond 401 on protected routes without a token", func(t *testing.T) {
		scenarios := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/api/user"},
			{http.MethodPost, "/api/post"},
		}

		for _, scenario := range scenarios {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(scenario.method, scenario.path, nil)
			r.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("%s %s without a token should return 401, got %d", scenario.method, scenario.path, rec.Code)
			}
		}
	})

	t.Run("Should respond 404 on an unknown path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET /does-not-exist should return 404, got %d", rec.Code)
		}
	})

	t.Run("Should respond 405 for a wrong method on an existing path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/health", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("DELETE /health should return 405, got %d", rec.Code)
		}
	})
}
