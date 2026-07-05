package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	t.Run("Should write content-type, status code and JSON body", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		data := struct {
			Message string `json:"message"`
		}{Message: "hello"}

		JSON(recorder, http.StatusCreated, data)

		if recorder.Code != http.StatusCreated {
			t.Errorf("JSON should write status %d. Got: %d", http.StatusCreated, recorder.Code)
		}
		if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("JSON should set Content-Type to application/json. Got: %q", contentType)
		}

		var decoded struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(recorder.Body).Decode(&decoded); err != nil {
			t.Fatalf("JSON body should be valid JSON: %v", err)
		}
		if decoded.Message != data.Message {
			t.Errorf("JSON should encode the given data. Expected message %q. Got: %q", data.Message, decoded.Message)
		}
	})

	t.Run("Should not write a body when data is nil", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		JSON(recorder, http.StatusNoContent, nil)

		if recorder.Code != http.StatusNoContent {
			t.Errorf("JSON should write status %d. Got: %d", http.StatusNoContent, recorder.Code)
		}
		if recorder.Body.Len() != 0 {
			t.Errorf("JSON should not write a body when data is nil. Got: %q", recorder.Body.String())
		}
	})
}

func TestError(t *testing.T) {
	t.Run("Should write an error body with the given status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		Error(recorder, http.StatusBadRequest, errors.New("boom"))

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Error should write status %d. Got: %d", http.StatusBadRequest, recorder.Code)
		}

		var decoded struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(recorder.Body).Decode(&decoded); err != nil {
			t.Fatalf("Error body should be valid JSON: %v", err)
		}
		if decoded.Error != "boom" {
			t.Errorf("Error should render the error message. Expected %q. Got: %q", "boom", decoded.Error)
		}
	})
}
