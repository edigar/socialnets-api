package controller

import (
	"encoding/json"
	"errors"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/database"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/repository"
	"github.com/edigar/socialnets-api/internal/response"
	"github.com/edigar/socialnets-api/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user entity.User
	if err = json.Unmarshal(body, &user); err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userUseCase := usecase.NewUserUseCase(repository.NewUserRepository(db))
	userId, err := userUseCase.Login(user.Email, user.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	token, err := authentication.CreateToken(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.Authentication{Id: userId, Token: token})
}
