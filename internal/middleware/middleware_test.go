package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edigar/socialnets-api/internal/authentication"
)

func TestAuthenticate(t *testing.T) {
	t.Run("Should call the next handler for a valid token", func(t *testing.T) {
		called := false
		next := func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}

		token, err := authentication.CreateToken("eedf21bf-dde8-4c85-b50b-89a1cba87c2e")
		if err != nil {
			t.Fatalf("CreateToken should not return an error: %v", err)
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()

		Authenticate(next)(recorder, request)

		if !called {
			t.Errorf("Authenticate should call the next handler for a valid token")
		}
		if recorder.Code != http.StatusOK {
			t.Errorf("Authenticate should respond with status %d for a valid token. Got: %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("Should respond 401 and not call next for an invalid token", func(t *testing.T) {
		called := false
		next := func(w http.ResponseWriter, r *http.Request) {
			called = true
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Authorization", "Bearer invalid-token")
		recorder := httptest.NewRecorder()

		Authenticate(next)(recorder, request)

		if called {
			t.Errorf("Authenticate should not call the next handler for an invalid token")
		}
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Authenticate should respond with status %d for an invalid token. Got: %d", http.StatusUnauthorized, recorder.Code)
		}
	})

	t.Run("Should respond 401 and not call next when the Authorization header is missing", func(t *testing.T) {
		called := false
		next := func(w http.ResponseWriter, r *http.Request) {
			called = true
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		Authenticate(next)(recorder, request)

		if called {
			t.Errorf("Authenticate should not call the next handler when the Authorization header is missing")
		}
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("Authenticate should respond with status %d when the Authorization header is missing. Got: %d", http.StatusUnauthorized, recorder.Code)
		}
	})
}

func TestLogger(t *testing.T) {
	t.Run("Should call the next handler", func(t *testing.T) {
		called := false
		next := func(w http.ResponseWriter, r *http.Request) {
			called = true
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()

		Logger(next)(recorder, request)

		if !called {
			t.Errorf("Logger should call the next handler")
		}
	})
}
