package usecase

import (
	"database/sql"
	"errors"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	errorType "github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/usecase/mock"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("Should login user with correct e-mail and password", func(t *testing.T) {
		userPassword := "123"

		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		userId, err := userUseCase.Login(usecase.MockUsers[0].Email, userPassword)

		if err != nil {
			t.Errorf("Login should not return an error for a valid email and password: %v. User: %v",
				err,
				usecase.MockUsers[0],
			)
		} else if userId != usecase.MockUsers[0].Id {
			t.Errorf("Login should return correct user id. user id: %v. User id expected: %v",
				userId,
				usecase.MockUsers[0].Id,
			)
		}
	})

	t.Run("Should not login user with incorrect e-mail", func(t *testing.T) {
		userPassword := "123"

		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		userId, err := userUseCase.Login("x", userPassword)

		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Login should return an error for a wrong email. Returned: %v. Error expected: %v",
				err,
				sql.ErrNoRows,
			)
		} else if userId != "" {
			t.Errorf("Login should return empty user id for wrong e-mail. Got: %v.", userId)
		}
	})

	t.Run("Should not login user with incorrect password", func(t *testing.T) {
		userPassword := "1"

		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		userId, err := userUseCase.Login(usecase.MockUsers[0].Email, userPassword)

		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			t.Errorf("Login should return an error for a wrong email. Returned: %v. Error expected: %v",
				err,
				bcrypt.ErrMismatchedHashAndPassword,
			)
		} else if userId != "" {
			t.Errorf("Login should return empty user id for wrong password. Got: %v.", userId)
		}
	})
}

func TestRegister(t *testing.T) {
	t.Run("Should register user with validated data", func(t *testing.T) {
		user := entity.User{
			Name:     "user",
			Nick:     "user",
			Email:    "user@mail",
			Password: "1",
		}

		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Register(&user)
		if err != nil {
			t.Errorf("Register should not return an error for a valid user data. User: %v. Error: %v", user, err)
		} else if user.Id != usecase.NEW_USER_ID {
			t.Errorf("Register should set an id for user. user id: %v. User id expected: %v",
				user.Id,
				usecase.NEW_USER_ID,
			)
		}
	})

	t.Run("Should not register user with invalid data", func(t *testing.T) {
		scenarios := []entity.User{
			{
				Nick:     "user",
				Email:    "user@mail",
				Password: "1",
			},
			{
				Name:     "a",
				Nick:     "user",
				Email:    "user@mail",
				Password: "1",
			},
			{
				Name:     "user",
				Email:    "user@mail",
				Password: "1",
			},
			{
				Name:     "user",
				Nick:     "user",
				Password: "1",
			},
			{
				Name:     "user",
				Nick:     "user",
				Email:    "usermail",
				Password: "1",
			},
			{
				Name:  "user",
				Nick:  "user",
				Email: "user@mail",
			},
		}

		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		for _, scenario := range scenarios {
			err := userUseCase.Register(&scenario)
			var euv *errorType.ErrorUserValidation
			if !errors.As(err, &euv) {
				t.Errorf("Register should return an ErrorUserValidation error for a non-valid user data. User: %v Returned: %v. Error expected: %T",
					scenario,
					err,
					euv.Error(),
				)
			} else if scenario.Id != "" {
				t.Errorf("Register should set empty user id for invalid user data. Got: %v.", scenario.Id)
			}
		}
	})

	//t.Run("Should get an DB error", func(t *testing.T) {
	//	user := entity.User{Name: "Test", Nick: "test", Email: "test@test"}
	//	userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
	//	err := userUseCase.Update(usecase.USER_ERROR, user)
	//	if err.Error() != "driver: bad connection" {
	//		t.Errorf("Update should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
	//	}
	//})
}

func TestGetUserById(t *testing.T) {
	t.Run("Should return valid user by his id", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		user, err := userUseCase.GetById(usecase.MockUsers[0].Id)
		if err != nil {
			t.Errorf("GetById should not return an error for a valid Id. User Id: %v. User returned: %v Error: %v",
				usecase.MockUsers[0].Id,
				user,
				err,
			)
		} else if user.Id != usecase.MockUsers[0].Id {
			t.Errorf("Invalid user returned. Got user: %v. Expected: %v",
				user,
				usecase.MockUsers[0],
			)
		}
	})

	t.Run("should return an error for a non-existent id", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		user, err := userUseCase.GetById("wrong-id")

		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("GetById should return ErrNoRows error for a wrong id. Got: %v. Error expected: %v",
				err,
				sql.ErrNoRows,
			)
		} else if user != (entity.User{}) {
			t.Errorf("GetById should return empty user for wrong id. Got: %v.", user)
		}
	})

	t.Run("should return an error for an empty id", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		user, err := userUseCase.GetById("")

		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("GetById should return ErrNoRows error for an empty id. Got: %v. Error expected: %v",
				err,
				sql.ErrNoRows,
			)
		} else if user != (entity.User{}) {
			t.Errorf("GetById should return empty user for empty id. Got: %v.", user)
		}
	})
}

func TestGetUserByNameOrNick(t *testing.T) {
	t.Run("Should return one user by his name", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		users, err := userUseCase.GetByNameOrNick(usecase.MockUsers[0].Name)
		if err != nil {
			t.Errorf("GetByNameOrNick should not return an error for a valid name. User name: %v. Users returned: %v Error: %v",
				usecase.MockUsers[0].Name,
				users,
				err,
			)
		} else if len(users) != 1 {
			t.Errorf("GetByNameOrNick should return just one user with valid name. Got %v.", users)
		}
	})

	t.Run("Should return one user by his nickname", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		users, err := userUseCase.GetByNameOrNick(usecase.MockUsers[0].Nick)
		if err != nil {
			t.Errorf("GetByNameOrNick should not return an error for a valid nickname. User nickname: %v. Users returned: %v Error: %v",
				usecase.MockUsers[0].Nick,
				users,
				err,
			)
		} else if len(users) != 1 {
			t.Errorf("GetByNameOrNick should return just one user with valid nickname. Got %v.", users)
		}
	})

	t.Run("Should return some users by his name or nick", func(t *testing.T) {
		str := "ano"
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		users, err := userUseCase.GetByNameOrNick(str)
		if err != nil {
			t.Errorf("GetByNameOrNick should not return an error for a valid string. string: %v. Users returned: %v Error: %v",
				str,
				users,
				err,
			)
		} else if len(users) != 2 {
			t.Errorf("GetByNameOrNick should return both of users with a string. Got %v.", users)
		}
	})

	t.Run("Should return empty list if no find user", func(t *testing.T) {
		str := "x"
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		users, err := userUseCase.GetByNameOrNick(str)
		if err != nil {
			t.Errorf("GetByNameOrNick should not return an error for a valid string. string: %v. Users returned: %v Error: %v",
				str,
				users,
				err,
			)
		} else if len(users) != 0 {
			t.Errorf("GetByNameOrNick should return empty list if no find user. Got %v.", users)
		}
	})

	t.Run("Should get an DB error", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		users, err := userUseCase.GetByNameOrNick(usecase.USER_ERROR)
		if err.Error() != "driver: bad connection" {
			t.Errorf("GetByNameOrNick should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}

		if users != nil {
			t.Errorf("GetByNameOrNick should not get any user if connection get an error. users: %v.", users)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Should update user with valid id", func(t *testing.T) {
		user := entity.User{Name: "Test", Nick: "test", Email: "test@test"}
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Update(usecase.MockUsers[0].Id, user)
		if err != nil {
			t.Errorf("Update should not return an error with valid data. Data sended: %v. User updated: %v Error: %v",
				user,
				usecase.MockUsers[0],
				err,
			)
		} else if user.Name != usecase.MockUsers[0].Name || user.Nick != usecase.MockUsers[0].Nick || user.Email != usecase.MockUsers[0].Email {
			t.Errorf("Update should update name, nick and email of mock user 0. Data sended: %v. User updated: %v", user, usecase.MockUsers[0])
		}
	})

	t.Run("Should not update user with non-valid id", func(t *testing.T) {
		user := entity.User{Name: "Test", Nick: "test", Email: "test@test"}
		originalUsers := usecase.MockUsers
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Update("x", user)
		if err != nil {
			t.Errorf("Update should not return an error with valid data. Data sended: %v. User updated: %v Error: %v",
				user,
				usecase.MockUsers[0],
				err,
			)
		} else if originalUsers[0] != usecase.MockUsers[0] || originalUsers[1] != usecase.MockUsers[1] {
			t.Errorf("Update should not update none of user with non-valid id. Data sended: %v. Users: %v", user, usecase.MockUsers)
		}
	})

	t.Run("Should not update user with non-valid data", func(t *testing.T) {
		scenarios := []entity.User{
			{
				Nick:  "user",
				Email: "user@mail",
			},
			{
				Name:  "a",
				Nick:  "user",
				Email: "user@mail",
			},
			{
				Name:  "user",
				Email: "user@mail",
			},
			{
				Name: "user",
				Nick: "user",
			},
			{
				Name:  "user",
				Nick:  "user",
				Email: "usermail",
			},
		}

		originalUsers := usecase.MockUsers
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		for _, scenario := range scenarios {
			err := userUseCase.Update(usecase.MockUsers[0].Id, scenario)
			var euv *errorType.ErrorUserValidation
			if !errors.As(err, &euv) {
				t.Errorf("Update should return an ErrorUserValidation error for a non-valid user data. User: %v Error returned: %v. Error expected: %T",
					scenario,
					err,
					euv.Error(),
				)
			} else if originalUsers[0] != usecase.MockUsers[0] || originalUsers[1] != usecase.MockUsers[1] {
				t.Errorf("Update should not update none of user with non-valid data. Data sended: %v. Original users: %v", scenario, originalUsers)
			}
		}
	})

	t.Run("Should get an DB error", func(t *testing.T) {
		user := entity.User{Name: "Test", Nick: "test", Email: "test@test"}
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Update(usecase.USER_ERROR, user)
		if err.Error() != "driver: bad connection" {
			t.Errorf("Update should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}
	})
}

func TestFollow(t *testing.T) {
	t.Run("Should return ErrOperationDenied if user id is equal to follower id", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Follow(usecase.MockUsers[0].Id, usecase.MockUsers[0].Id)
		if !errors.Is(err, ErrOperationDenied) {
			t.Errorf("Follow should return ErrOperationDenied if user id is equal to follower id. Got %v.", err)
		}
	})

	t.Run("Should return bad connection error DB", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Follow(usecase.MockUsers[0].Id, usecase.MockUsers[1].Id)
		if err.Error() != "driver: bad connection" {
			t.Errorf("Folow should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}
	})
}

func TestUnfollow(t *testing.T) {
	t.Run("Should return ErrOperationDenied if user id is equal to follower id", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Unfollow(usecase.MockUsers[0].Id, usecase.MockUsers[0].Id)
		if !errors.Is(err, ErrOperationDenied) {
			t.Errorf("Unfollow should return ErrOperationDenied if user id is equal to follower id. Got %v.", err)
		}
	})

	t.Run("Should return bad connection error DB", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Unfollow(usecase.MockUsers[0].Id, usecase.MockUsers[1].Id)
		if err.Error() != "driver: bad connection" {
			t.Errorf("Folow should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Run("Should update user password", func(t *testing.T) {
		passwordDto := dto.Password{New: "abc", Current: "123"}
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		oldPassword := usecase.MockUsers[0].Password
		err := userUseCase.UpdatePassword(usecase.MockUsers[0].Id, passwordDto)
		if err != nil {
			t.Errorf("UpdatePassword should not return an error for a valid user and password. Error: %v", err)
		} else if usecase.MockUsers[0].Password == oldPassword {
			t.Errorf("UpdatePassword should change user password. old password(123): %v Updated(abc): %v",
				oldPassword,
				usecase.MockUsers[0].Password,
			)
		}
	})

	t.Run("Should not update user password if current password is wrong", func(t *testing.T) {
		passwordDto := dto.Password{New: "abc", Current: "wrong-password"}
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		oldPassword := usecase.MockUsers[0].Password
		err := userUseCase.UpdatePassword(usecase.MockUsers[0].Id, passwordDto)
		if !errors.Is(err, ErrWrongPassword) {
			t.Errorf("UpdatePassword should return ErrWrongPassword error for wrong password. Got: %v", err)
		} else if usecase.MockUsers[0].Password != oldPassword {
			t.Errorf("UpdatePassword should not change user password if current password is wrong. old password(123): %v Updated(abc): %v",
				oldPassword,
				usecase.MockUsers[0].Password,
			)
		}
	})

	t.Run("Should get error if userId doesn't exist", func(t *testing.T) {
		passwordDto := dto.Password{New: "abc", Current: "123"}
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.UpdatePassword("wrong-id", passwordDto)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("UpdatePassword should return ErrNoRows error for an Nonexistent id. Got: %v. Error expected: %v",
				err,
				sql.ErrNoRows,
			)
		}
	})
}

func TestGetFollowers(t *testing.T) {
	t.Run("Should get user followers", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		followers, err := userUseCase.GetFollowers(usecase.MockUsers[0].Id)
		if err != nil {
			t.Errorf("GetFollowers should not return an error for a valid user id. Error: %v", err)
		}

		expectedFollowers := []entity.User{usecase.MockUsers[1], usecase.MockUsers[2]}
		if !reflect.DeepEqual(followers, expectedFollowers) {
			t.Errorf("GetFollowers should return user followers. Expected: %v. Got: %v", expectedFollowers, followers)
		}
	})

	t.Run("Should get an DB error", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		followers, err := userUseCase.GetFollowers(usecase.USER_ERROR)
		if err.Error() != "driver: bad connection" {
			t.Errorf("GetFollowers should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}

		if followers != nil {
			t.Errorf("GetFollowers should not get any follower if connection get an error. Posts: %v.", followers)
		}
	})
}

func TestGetFollowing(t *testing.T) {
	t.Run("Should get user following", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		following, err := userUseCase.GetFollowing(usecase.MockUsers[0].Id)
		if err != nil {
			t.Errorf("GetFollowing should not return an error for a valid user id. Error: %v", err)
		}

		expectedFollowing := []entity.User{usecase.MockUsers[1], usecase.MockUsers[2]}
		if !reflect.DeepEqual(following, expectedFollowing) {
			t.Errorf("GetFollowing should return user followers. Expected: %v. Got: %v", expectedFollowing, following)
		}
	})

	t.Run("Should get an DB error", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		following, err := userUseCase.GetFollowing(usecase.USER_ERROR)
		if err.Error() != "driver: bad connection" {
			t.Errorf("GetFollowing should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}

		if following != nil {
			t.Errorf("GetFollowing should not get any following if connection get an error. Posts: %v.", following)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Should get nil if delete an user", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Delete(usecase.MockUsers[2].Id)
		if err != nil {
			t.Errorf("Delete should not return an error if delete user. Error: %v", err)
		}
	})

	t.Run("Should get an DB error", func(t *testing.T) {
		userUseCase := NewUserUseCase(usecase.NewMockUserRepository())
		err := userUseCase.Delete(usecase.USER_ERROR)
		if err.Error() != "driver: bad connection" {
			t.Errorf("Delete should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}
	})
}
