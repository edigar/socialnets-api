package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/response"
	"github.com/edigar/socialnets-api/internal/usecase"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strings"
)

func (c *UserController) PostUser(w http.ResponseWriter, r *http.Request) {
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

	if err = c.userUseCase.Register(&user); err != nil {
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

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("search"))

	users, err := c.userUseCase.GetByNameOrNick(nameOrNick)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, users)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	user, err := c.userUseCase.GetById(userId)
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

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = c.userUseCase.Update(userId, user); err != nil {
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

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if err = c.userUseCase.Delete(userId); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *UserController) Follow(w http.ResponseWriter, r *http.Request) {
	follower, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	if err = c.userUseCase.Follow(userId, follower); err != nil {
		if errors.Is(err, usecase.ErrOperationDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *UserController) Unfollow(w http.ResponseWriter, r *http.Request) {
	follower, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	if err = c.userUseCase.Unfollow(userId, follower); err != nil {
		if errors.Is(err, usecase.ErrOperationDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *UserController) GetFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	followers, err := c.userUseCase.GetFollowers(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, followers)
}

func (c *UserController) GetFollowing(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	following, err := c.userUseCase.GetFollowing(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, following)
}

func (c *UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	if err = c.userUseCase.UpdatePassword(userId, password); err != nil {
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
