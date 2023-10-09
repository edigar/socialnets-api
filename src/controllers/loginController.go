package controllers

import (
	"encoding/json"
	"github.com/edigar/socialnets-api/src/authentication"
	"github.com/edigar/socialnets-api/src/crypt"
	"github.com/edigar/socialnets-api/src/database"
	"github.com/edigar/socialnets-api/src/models"
	"github.com/edigar/socialnets-api/src/repositories"
	"github.com/edigar/socialnets-api/src/responses"
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

	responses.JSON(w, http.StatusOK, token)
}
