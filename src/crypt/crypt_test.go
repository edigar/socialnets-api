package crypt

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	t.Run("Should generate and verify a correct password hash", func(t *testing.T) {
		password := "secret"

		hashedPassword, err := Hash(password)
		if err != nil {
			t.Errorf("Hash should generate hashed password, but got an error: %v", err)
		}

		err = Verify(string(hashedPassword), password)
		if err != nil {
			t.Errorf("Verify should return nil for correct password, but got an error: %v", err)
		}
	})

	t.Run("Should return an error if verify a wrong password", func(t *testing.T) {
		password := "secret"

		hashedPassword, err := Hash(password)
		incorrectPassword := "wrong_password"
		err = Verify(string(hashedPassword), incorrectPassword)
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			t.Errorf("Verify should return an error for incorrect password")
		}
	})
}
