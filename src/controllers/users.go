package controllers

import (
	"devbook_api/src/authentication"
	"devbook_api/src/database"
	"devbook_api/src/models"
	"devbook_api/src/repositories"
	"devbook_api/src/responses"
	"devbook_api/src/security"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err = user.Prepare(true); err != nil {
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
	userId, err := repository.Create(user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	user.Id = userId

	responses.JSON(w, http.StatusCreated, user)

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	nameOrNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	users, err := repository.Get(nameOrNick)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
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

	user, err := repository.GetById(userId)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdFromToken, err := authentication.ExtractUserId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userId != userIdFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to update a user if not yours"))
		return
	}

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

	if err = user.Prepare(false); err != nil {
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

	err = repository.Update(userId, user)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdFromToken, err := authentication.ExtractUserId(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userId != userIdFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to delete a user if not yours"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	err = repository.Delete(userId)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdFromToken, err := authentication.ExtractUserId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if followerId == userIdFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to follow your own user"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	err = repository.FollowUser(followerId, userIdFromToken)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	followerId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdFromToken, err := authentication.ExtractUserId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if followerId == userIdFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to unfollow your own user"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	err = repository.UnfollowUser(followerId, userIdFromToken)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
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

	followers, err := repository.GetFollowers(userId)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func GetFollowings(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
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

	followers, err := repository.GetFollowings(userId)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, followers)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userIdFromToken, err := authentication.ExtractUserId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	if userId != userIdFromToken {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to change other user's password"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password models.Password

	if err := json.Unmarshal(body, &password); err != nil {
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

	savedPassword, err := repository.SearchPassword(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if err := security.VerifyPassword(password.Current, savedPassword); err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("it is not possible to change other user's password"))
		return
	}

	hashedPassword, err := security.Hash(password.New)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	err = repository.UpdatePassword(userId, hashedPassword)

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
