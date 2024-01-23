package usecase

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/entity"
	"strings"
)

type MockUserRepository struct{}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

const NEW_USER_ID = "6705a6cd-eb7b-488b-9e94-7685c95f2707"

var MockUsers = []entity.User{
	{
		Id:       "93226a19-86d6-4ad7-a215-d5999c2870c4",
		Name:     "Fulano",
		Nick:     "fulano",
		Email:    "fulano@mail",
		Password: "$2a$10$IAAZ5kkLrbp.Fczn8Q1M/O8f1Nsjf0sWqFpJwwwsRK.7CkWCwv/sC", //123
	},
	{
		Id:       "d9b56fd4-31b7-4bd5-958f-99028ca5e79a",
		Name:     "Beltrano",
		Nick:     "beltrano",
		Email:    "beltrano@mail",
		Password: "$2a$10$hSBo5G8bQYHjktHk/g4Hz.B0EhogAlz6DYaGkyWqExTTahFAq8vn.", //321
	},
}

func (mr MockUserRepository) Create(user entity.User) (string, error) {
	return NEW_USER_ID, nil
}

func (mr MockUserRepository) FetchByNameOrNick(nameOrNick string) ([]entity.User, error) {
	var users []entity.User
	for _, user := range MockUsers {
		if strings.Contains(user.Name, nameOrNick) || strings.Contains(user.Nick, nameOrNick) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (mr MockUserRepository) FetchById(userId string) (entity.User, error) {
	for _, user := range MockUsers {
		if user.Id == userId {
			return user, nil
		}
	}

	return entity.User{}, sql.ErrNoRows
}

func (mr MockUserRepository) FetchByEmail(email string) (entity.User, error) {
	for _, user := range MockUsers {
		if user.Email == email {
			return user, nil
		}
	}

	return entity.User{}, sql.ErrNoRows
}

func (mr MockUserRepository) Update(userId string, user entity.User) error {
	for i, mockUser := range MockUsers {
		if mockUser.Id == userId {
			MockUsers[i].Name = user.Name
			MockUsers[i].Nick = user.Nick
			MockUsers[i].Email = user.Email
		}
	}

	return nil
}

func (mr MockUserRepository) Delete(userId string) error {
	return nil
}

func (mr MockUserRepository) Follow(userId, follower string) error {
	return nil
}

func (mr MockUserRepository) Unfollow(userId, follower string) error {
	return nil
}

func (mr MockUserRepository) FetchFollowers(userId string) ([]entity.User, error) {
	return nil, nil
}

func (mr MockUserRepository) FetchFollowing(userId string) ([]entity.User, error) {
	return nil, nil
}

func (mr MockUserRepository) FetchPasswordById(userId string) (string, error) {
	for _, user := range MockUsers {
		if user.Id == userId {
			return user.Password, nil
		}
	}

	return "", sql.ErrNoRows
}

func (mr MockUserRepository) UpdatePassword(userId string, passwordHash string) error {
	for i, user := range MockUsers {
		if user.Id == userId {
			MockUsers[i].Password = passwordHash
		}
	}

	return nil
}
