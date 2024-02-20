package controllers

import (
	"devbook_api/src/authentication"
	"devbook_api/src/database"
	"devbook_api/src/models"
	"devbook_api/src/repositories"
	"devbook_api/src/responses"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}
	bodyRequest, err := io.ReadAll(r.Body)

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post

	err = json.Unmarshal(bodyRequest, &post)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorId = userId

	if err := post.Prepare(); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	repository := repositories.NewPostRepository(db)

	post.Id, err = repository.Create(userId, post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, post)

}
func GetPosts(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	repository := repositories.NewPostRepository(db)

	posts, err := repository.GetPosts(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)

}

func GetPost(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
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

	repository := repositories.NewPostRepository(db)

	post, err := repository.GetPostById(postId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	postId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
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

	repository := repositories.NewPostRepository(db)

	postSaved, err := repository.GetPostById(postId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postSaved.AuthorId != userId {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to update a post that is not yours"))
		return
	}
	bodyRequest, err := io.ReadAll(r.Body)

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post models.Post

	err = json.Unmarshal(bodyRequest, &post)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	err = post.Prepare()
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	err = repository.Update(postId, post)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)

	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	postId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
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

	repository := repositories.NewPostRepository(db)

	postSaved, err := repository.GetPostById(postId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	if postSaved.AuthorId != userId {
		responses.Error(w, http.StatusForbidden, errors.New("it is not possible to delete a post that is not yours"))
		return
	}

	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = repository.Delete(postId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewPostRepository(db)

	posts, err := repository.GetPostsByUser(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)

}

func LikePost(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewPostRepository(db)

	err = repository.Like(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
func DislikePost(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewPostRepository(db)

	err = repository.Dislike(userId)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
