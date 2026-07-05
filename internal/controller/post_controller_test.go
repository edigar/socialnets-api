package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/usecase"
)

func TestPostPost(t *testing.T) {
	validBody := `{"title":"Title","content":"Content"}`

	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		c.PostPost(rec, newRequest(http.MethodPost, validBody))
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		c.PostPost(rec, withAuth(newRequest(http.MethodPost, `{`), testUserID))
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 400 for a validation error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{createPostFn: func(*entity.Post) error {
			return errorType.NewErrorPostValidation("title is required")
		}})
		rec := httptest.NewRecorder()
		c.PostPost(rec, withAuth(newRequest(http.MethodPost, validBody), testUserID))
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 201 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{createPostFn: func(*entity.Post) error { return nil }})
		rec := httptest.NewRecorder()
		c.PostPost(rec, withAuth(newRequest(http.MethodPost, validBody), testUserID))
		assertStatus(t, rec.Code, http.StatusCreated)
	})

	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{createPostFn: func(*entity.Post) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		c.PostPost(rec, withAuth(newRequest(http.MethodPost, validBody), testUserID))
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetPosts(t *testing.T) {
	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		c.GetPosts(rec, newRequest(http.MethodGet, ""))
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 200 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getByUserFn: func(string) ([]entity.Post, error) {
			return []entity.Post{}, nil
		}})
		rec := httptest.NewRecorder()
		c.GetPosts(rec, withAuth(newRequest(http.MethodGet, ""), testUserID))
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getByUserFn: func(string) ([]entity.Post, error) {
			return nil, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		c.GetPosts(rec, withAuth(newRequest(http.MethodGet, ""), testUserID))
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetPost(t *testing.T) {
	t.Run("Should return 400 for a non-numeric post id", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"postId": "abc"})
		c.GetPost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 200 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getByIdFn: func(uint64) (entity.Post, error) {
			return entity.Post{Id: 1}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"postId": "1"})
		c.GetPost(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getByIdFn: func(uint64) (entity.Post, error) {
			return entity.Post{}, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"postId": "1"})
		c.GetPost(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestUpdatePost(t *testing.T) {
	validBody := `{"title":"Title","content":"Content"}`

	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "1"})
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 400 for a non-numeric post id", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "abc"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, `{`), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 400 for a validation error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{updateFn: func(string, uint64, entity.Post) error {
			return errorType.NewErrorPostValidation("invalid")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 403 when access is denied", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{updateFn: func(string, uint64, entity.Post) error {
			return usecase.ErrAccessDenied
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{updateFn: func(string, uint64, entity.Post) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{updateFn: func(string, uint64, entity.Post) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.UpdatePost(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestDeletePost(t *testing.T) {
	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"postId": "1"})
		c.DeletePost(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 400 for a non-numeric post id", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"postId": "abc"})
		req = withAuth(req, testUserID)
		c.DeletePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 403 when access is denied", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{deleteFn: func(uint64, string) error {
			return usecase.ErrAccessDenied
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.DeletePost(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{deleteFn: func(uint64, string) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.DeletePost(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{deleteFn: func(uint64, string) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"postId": "1"})
		req = withAuth(req, testUserID)
		c.DeletePost(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetUserPosts(t *testing.T) {
	t.Run("Should return 200 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getUserPostsFn: func(string) ([]entity.Post, error) {
			return []entity.Post{}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetUserPosts(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{getUserPostsFn: func(string) ([]entity.Post, error) {
			return nil, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetUserPosts(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestLikePost(t *testing.T) {
	t.Run("Should return 400 for a non-numeric post id", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "abc"})
		c.LikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{likePostFn: func(uint64) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "1"})
		c.LikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{likePostFn: func(uint64) error { return errors.New("db down") }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "1"})
		c.LikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestUnlikePost(t *testing.T) {
	t.Run("Should return 400 for a non-numeric post id", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "abc"})
		c.UnlikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{unLikePostFn: func(uint64) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "1"})
		c.UnlikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewPostController(&mockPostUseCase{unLikePostFn: func(uint64) error { return errors.New("db down") }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"postId": "1"})
		c.UnlikePost(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}
