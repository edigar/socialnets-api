package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/database"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/repository"
	"github.com/edigar/socialnets-api/internal/response"
	"github.com/edigar/socialnets-api/internal/usecase"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strings"
)

func PostUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user entity.User
	if err = json.Unmarshal(body, &user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.Register(&user); err != nil {
		var uve *errorType.ErrorUserValidation
		if errors.As(err, &uve) {
			response.Error(w, http.StatusBadRequest, uve.Err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("search"))
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	users, err := userUseCase.GetByNameOrNick(nameOrNick)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))

	user, err := userUseCase.GetById(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if (entity.User{}) == user {
		response.JSON(w, http.StatusNotFound, nil)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	tokenUserId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	if userId != tokenUserId {
		response.Error(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user entity.User
	if err = json.Unmarshal(body, &user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.Update(userId, user); err != nil {
		var uve *errorType.ErrorUserValidation
		if errors.As(err, &uve) {
			response.Error(w, http.StatusBadRequest, uve.Err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	tokenUserId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	if userId != tokenUserId {
		response.Error(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.Delete(userId); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func Follow(w http.ResponseWriter, r *http.Request) {
	follower, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.Follow(userId, follower); err != nil {
		if errors.Is(err, usecase.ErrOperationDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	follower, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.Unfollow(userId, follower); err != nil {
		if errors.Is(err, usecase.ErrOperationDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	followers, err := userUseCase.GetFollowers(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	following, err := userUseCase.GetFollowing(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, following)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	tokenUserId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	if userId != tokenUserId {
		response.Error(w, http.StatusForbidden, errors.New("access denied"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password dto.Password
	if err = json.Unmarshal(body, &password); err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	if err = userUseCase.UpdatePassword(userId, password); err != nil {
		if errors.Is(err, usecase.ErrWrongPassword) {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			response.Error(w, http.StatusBadRequest, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
