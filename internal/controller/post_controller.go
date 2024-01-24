package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/database"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/repository"
	"github.com/edigar/socialnets-api/internal/response"
	"github.com/edigar/socialnets-api/internal/usecase"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func PostPost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	var post entity.Post
	if err = json.Unmarshal(body, &post); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post.AuthorId = userId

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	if err = postUseCase.CreatePost(&post); err != nil {
		var epv *errorType.ErrorPostValidation
		if errors.As(err, &epv) {
			response.Error(w, http.StatusBadRequest, epv.Err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusCreated, post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	posts, err := postUseCase.GetByUser(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	post, err := postUseCase.GetById(postId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	var post entity.Post
	if err = json.Unmarshal(body, &post); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	if err = postUseCase.Update(userId, postId, post); err != nil {
		var epv *errorType.ErrorPostValidation
		if errors.As(err, &epv) {
			response.Error(w, http.StatusBadRequest, epv.Err)
			return
		}
		if errors.Is(err, usecase.ErrAccessDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	if err = postUseCase.Delete(postId, userId); err != nil {
		if errors.Is(err, usecase.ErrAccessDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	posts, err := postUseCase.GetUserPosts(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
	}

	response.JSON(w, http.StatusOK, posts)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	//userId, err := authentication.ExtractUserId(r)
	//if err != nil {
	//	responses.Error(w, http.StatusUnauthorized, err)
	//	return
	//}
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	err = postUseCase.LikePost(postId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	db, err := database.Connect()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	postUseCase := usecase.NewPostUseCase(repository.NewPostRepository(db))
	err = postUseCase.UnLikePost(postId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
