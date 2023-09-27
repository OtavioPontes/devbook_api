package controllers

import (
	"devbook_api/src/authentication"
	"devbook_api/src/database"
	"devbook_api/src/models"
	"devbook_api/src/repositories"
	"devbook_api/src/responses"
	"devbook_api/src/security"
	"encoding/json"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	repository := repositories.NewUserRepository(db)

	savedUser, err := repository.SearchByEmail(user.Email)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err := security.VerifyPassword(user.Password, savedUser.Password); err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(savedUser.Id)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, token)

}
