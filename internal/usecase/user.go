package usecase

import (
	"errors"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/repository"
	"github.com/edigar/socialnets-api/pkg/crypt"
)

var (
	ErrOperationDenied = errors.New("operation denied")
	ErrWrongPassword   = errors.New("wrong password")
)

type UserUseCase struct {
	userRepository repository.User
}

func NewUserUseCase(userRepository repository.User) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}

func (u *UserUseCase) Login(email string, password string) (string, error) {
	user, err := u.userRepository.FetchByEmail(email)
	if err != nil {
		return "", err
	}

	if err = crypt.Verify(user.Password, password); err != nil {
		return "", err
	}

	return user.Id, nil
}

func (u *UserUseCase) Register(user *entity.User) error {
	err := user.Prepare("register")
	if err != nil {
		return err
	}

	user.Id, err = u.userRepository.Create(*user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetById(id string) (entity.User, error) {
	user, err := u.userRepository.FetchById(id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *UserUseCase) GetByNameOrNick(nameOrNick string) ([]entity.User, error) {
	user, err := u.userRepository.FetchBy(nameOrNick)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) Update(userId string, user entity.User) error {
	err := user.Prepare("edit")
	if err != nil {
		return err
	}

	if err = u.userRepository.Update(userId, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) Delete(userId string) error {
	if err := u.userRepository.Delete(userId); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) Follow(userId string, follower string) error {
	if follower == userId {
		return ErrOperationDenied
	}

	if err := u.userRepository.Follow(userId, follower); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) Unfollow(userId string, follower string) error {
	if follower == userId {
		return ErrOperationDenied
	}

	if err := u.userRepository.Unfollow(userId, follower); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetFollowers(userId string) ([]entity.User, error) {
	users, err := u.userRepository.FetchFollowers(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserUseCase) GetFollowing(userId string) ([]entity.User, error) {
	users, err := u.userRepository.FetchFollowing(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserUseCase) UpdatePassword(userId string, password dto.Password) error {
	passwordDb, err := u.userRepository.FetchPasswordById(userId)
	if err != nil {
		return err
	}

	if err = crypt.Verify(passwordDb, password.Current); err != nil {
		return ErrWrongPassword
	}

	passwordHash, err := crypt.Hash(password.New)
	if err != nil {
		return err
	}

	if err = u.userRepository.UpdatePassword(userId, string(passwordHash)); err != nil {
		return err
	}

	return nil
}
