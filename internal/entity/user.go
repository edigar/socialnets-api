package entity

import (
	"fmt"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/pkg/crypt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type User struct {
	Id        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Nick      string     `json:"nick,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}
	if err := user.format(step); err != nil {
		return err
	}

	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errorType.NewErrorUserValidation("username is required")
	}

	if !regexp.MustCompile("^[A-Za-z\\s]{3,}$").MatchString(user.Name) {
		return errorType.NewErrorUserValidation("username must have three or more characters")
	}

	if user.Nick == "" {
		return errorType.NewErrorUserValidation("nick is required")
	}

	if user.Email == "" {
		return errorType.NewErrorUserValidation("email is required")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errorType.NewErrorUserValidation(fmt.Sprintf("invalid email. %s", err))
	}

	if step == "register" && user.Password == "" {
		return errorType.NewErrorUserValidation("password is required")
	}

	return nil
}

func (user *User) format(step string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if step == "register" {
		passwordHash, err := crypt.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(passwordHash)
	}

	return nil
}
