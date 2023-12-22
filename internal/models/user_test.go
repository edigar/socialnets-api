package models

import (
	"testing"
	"time"
)

type UserScenarios struct {
	user     User
	expected User
}

func TestUserPrepareWithValidUser(t *testing.T) {
	t.Run("Should format and validate data with valid user", func(t *testing.T) {
		createdAt := time.Now()
		scenarios := []UserScenarios{
			{
				User{
					Id:        1,
					Name:      "teste name",
					Nick:      "nick",
					Email:     "name@mail.com",
					Password:  "123",
					CreatedAt: createdAt,
				},
				User{
					Id:        1,
					Name:      "teste name",
					Nick:      "nick",
					Email:     "name@mail.com",
					Password:  "123",
					CreatedAt: createdAt,
				},
			},
			{
				User{
					Id:        1,
					Name:      "  teste name    ",
					Nick:      "   nick  ",
					Email:     "    name@mail.com   ",
					Password:  "123",
					CreatedAt: createdAt,
				},
				User{
					Id:        1,
					Name:      "teste name",
					Nick:      "nick",
					Email:     "name@mail.com",
					Password:  "123",
					CreatedAt: createdAt,
				},
			},
			{
				User{
					Id:        1,
					Name:      "  teste  name    ",
					Nick:      "   nick  ",
					Email:     "    name@mail.com   ",
					Password:  "123",
					CreatedAt: createdAt,
				},
				User{
					Id:        1,
					Name:      "teste  name",
					Nick:      "nick",
					Email:     "name@mail.com",
					Password:  "123",
					CreatedAt: createdAt,
				},
			},
		}

		for _, scenario := range scenarios {
			err := scenario.user.Prepare("no register")

			if err != nil {
				t.Errorf("User prepare should not return an error for a valid post: %v. Scenario: %v", err, scenario.user)
			}
			if scenario.user != scenario.expected {
				t.Errorf(
					"User prepare should correctly format the user data. User: %v. User expected: %v",
					scenario.user,
					scenario.expected,
				)
			}
		}
	})

	t.Run("Should generate password hash with valid user on register step", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "Name", Nick: "nick", Email: "name@mail", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("register")

		if err != nil {
			t.Errorf("User prepare should not return error on register step")
		}
		if user.Password == "123" {
			t.Errorf("User prepare should generate password hash on register step")
		}
	})

	t.Run("Should return an error if name has less 3 characters", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "ab", Nick: "nick", Email: "name@mail.com", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("no register")

		if err.Error() != "username must have three or more characters" {
			t.Errorf("User prepare should return error if name has less than three letters. Error: %v", err)
		}
	})

	t.Run("Should return an error if name has a non character", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "name1", Nick: "nick", Email: "name@mail.com", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("no register")

		if err.Error() != "username must have three or more characters" {
			t.Errorf("User prepare should return error if name has a non character. Error: %v", err)
		}
	})

	t.Run("Should return an 'username is required' error if name is empty", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "", Nick: "nick", Email: "name@mail.com", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("no register")

		if err.Error() != "username is required" {
			t.Errorf("User prepare should return error if name is empty. Error: %v", err)
		}
	})

	t.Run("Should return an 'nick is required' error if nick is empty", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "Name", Nick: "", Email: "name@mail.com", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("no register")

		if err.Error() != "nick is required" {
			t.Errorf("User prepare should return error if nick is empty. Error: %v: ", err)
		}
	})

	t.Run("Should return an 'email is required' error if email is empty", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "Name", Nick: "nick", Email: "", Password: "123", CreatedAt: createdAt}
		err := user.Prepare("no register")

		if err.Error() != "email is required" {
			t.Errorf("User prepare should return error if email is empty. Error: %v", err)
		}
	})

	t.Run("Should return an error if email is invalid", func(t *testing.T) {
		createdAt := time.Now()
		users := []User{
			{
				Id:        1,
				Name:      "name",
				Nick:      "nick",
				Email:     "a",
				Password:  "123",
				CreatedAt: createdAt,
			},
			{
				Id:        1,
				Name:      "name",
				Nick:      "nick",
				Email:     "a@",
				Password:  "123",
				CreatedAt: createdAt,
			},
			{
				Id:        1,
				Name:      "name",
				Nick:      "nick",
				Email:     "@a",
				Password:  "123",
				CreatedAt: createdAt,
			},
			{
				Id:        1,
				Name:      "name",
				Nick:      "nick",
				Email:     "abc",
				Password:  "123",
				CreatedAt: createdAt,
			},
		}

		for _, user := range users {
			err := user.Prepare("no register")

			if err == nil {
				t.Errorf("User prepare should return error if email is invalid")
			}
		}
	})

	t.Run("Should return an error if password is empty on register step", func(t *testing.T) {
		createdAt := time.Now()
		user := User{Id: 1, Name: "Name", Nick: "nick", Email: "name@mail", Password: "", CreatedAt: createdAt}
		err := user.Prepare("register")

		if err.Error() != "password is required" {
			t.Errorf("User prepare should return error if password is empty on register step. Error %v", err)
		}
	})
}
