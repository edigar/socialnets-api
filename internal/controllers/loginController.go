package controllers

import (
	"encoding/json"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/database"
	"github.com/edigar/socialnets-api/internal/models"
	"github.com/edigar/socialnets-api/internal/repositories"
	"github.com/edigar/socialnets-api/internal/responses"
	"github.com/edigar/socialnets-api/pkg/crypt"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var userRequest models.User
	if err = json.Unmarshal(body, &userRequest); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userRepository := repositories.UserRepositoryFactory(db)
	userDb, err := userRepository.GetByEmail(userRequest.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err = crypt.Verify(userDb.Password, userRequest.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(userDb.Id)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, models.AuthenticationDTO{Id: userDb.Id, Token: token})
}
