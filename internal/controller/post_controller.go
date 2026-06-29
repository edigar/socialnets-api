package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/response"
	"github.com/edigar/socialnets-api/internal/usecase"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func (c *PostController) PostPost(w http.ResponseWriter, r *http.Request) {
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

	if err = c.postUseCase.CreatePost(&post); err != nil {
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

func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	userId, err := authentication.ExtractUserId(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err)
		return
	}

	posts, err := c.postUseCase.GetByUser(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (c *PostController) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	post, err := c.postUseCase.GetById(postId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, post)
}

func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	if err = c.postUseCase.Update(userId, postId, post); err != nil {
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

func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	if err = c.postUseCase.Delete(postId, userId); err != nil {
		if errors.Is(err, usecase.ErrAccessDenied) {
			response.Error(w, http.StatusForbidden, err)
			return
		}

		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *PostController) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := fmt.Sprintf("%s", params["userId"])

	posts, err := c.postUseCase.GetUserPosts(userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

func (c *PostController) LikePost(w http.ResponseWriter, r *http.Request) {
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

	if err = c.postUseCase.LikePost(postId); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *PostController) UnlikePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId, err := strconv.ParseUint(params["postId"], 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	if err = c.postUseCase.UnLikePost(postId); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
