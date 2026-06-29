package errorType

import (
	"errors"
	"testing"
)

func TestErrorUserValidation(t *testing.T) {
	err := NewErrorUserValidation("name is required")

	if err.Error() != "name is required" {
		t.Errorf("Error() should return the message, got %q", err.Error())
	}

	var target *ErrorUserValidation
	if !errors.As(err, &target) {
		t.Errorf("errors.As should detect *ErrorUserValidation")
	}
}

func TestErrorPostValidation(t *testing.T) {
	err := NewErrorPostValidation("title is required")

	if err.Error() != "title is required" {
		t.Errorf("Error() should return the message, got %q", err.Error())
	}

	var target *ErrorPostValidation
	if !errors.As(err, &target) {
		t.Errorf("errors.As should detect *ErrorPostValidation")
	}
}
