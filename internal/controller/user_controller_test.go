package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

const (
	testUserID  = "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"
	otherUserID = "11111111-1111-1111-1111-111111111111"
)

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("expected status %d, got %d", want, got)
	}
}

func TestPostUser(t *testing.T) {
	validBody := `{"name":"User Name","nick":"user","email":"user@mail.com","password":"123"}`

	t.Run("Should return 201 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{registerFn: func(*entity.User) error { return nil }})
		rec := httptest.NewRecorder()
		c.PostUser(rec, newRequest(http.MethodPost, validBody))
		assertStatus(t, rec.Code, http.StatusCreated)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		c.PostUser(rec, newRequest(http.MethodPost, `{`))
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 400 for a validation error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{registerFn: func(*entity.User) error {
			return errorType.NewErrorUserValidation("name is required")
		}})
		rec := httptest.NewRecorder()
		c.PostUser(rec, newRequest(http.MethodPost, validBody))
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 500 for an unexpected error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{registerFn: func(*entity.User) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		c.PostUser(rec, newRequest(http.MethodPost, validBody))
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestLogin(t *testing.T) {
	body := `{"email":"user@mail.com","password":"123"}`

	t.Run("Should return 200 and a token on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{loginFn: func(string, string) (string, error) {
			return testUserID, nil
		}})
		rec := httptest.NewRecorder()
		c.Login(rec, newRequest(http.MethodPost, body))
		assertStatus(t, rec.Code, http.StatusOK)

		var decoded struct {
			Id    string `json:"id"`
			Token string `json:"token"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("response should be valid JSON: %v", err)
		}
		if decoded.Id != testUserID || decoded.Token == "" {
			t.Errorf("Login should return the user id and a non-empty token. Got: %+v", decoded)
		}
	})

	t.Run("Should return 401 for wrong credentials", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{loginFn: func(string, string) (string, error) {
			return "", bcrypt.ErrMismatchedHashAndPassword
		}})
		rec := httptest.NewRecorder()
		c.Login(rec, newRequest(http.MethodPost, body))
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 500 for an unexpected error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{loginFn: func(string, string) (string, error) {
			return "", errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		c.Login(rec, newRequest(http.MethodPost, body))
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		c.Login(rec, newRequest(http.MethodPost, `{`))
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})
}

func TestGetUsers(t *testing.T) {
	t.Run("Should return 200 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getByNameOrNickFn: func(string) ([]entity.User, error) {
			return []entity.User{{Id: testUserID}}, nil
		}})
		rec := httptest.NewRecorder()
		c.GetUsers(rec, newRequest(http.MethodGet, ""))
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getByNameOrNickFn: func(string) ([]entity.User, error) {
			return nil, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		c.GetUsers(rec, newRequest(http.MethodGet, ""))
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Should return 200 when the user exists", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getByIdFn: func(string) (entity.User, error) {
			return entity.User{Id: testUserID, Name: "User"}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetUser(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("Should return 404 when the user is empty", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getByIdFn: func(string) (entity.User, error) {
			return entity.User{}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetUser(rec, req)
		assertStatus(t, rec.Code, http.StatusNotFound)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getByIdFn: func(string) (entity.User, error) {
			return entity.User{}, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetUser(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestUpdateUser(t *testing.T) {
	validBody := `{"name":"New Name","nick":"new","email":"new@mail.com"}`

	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"userId": testUserID})
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 403 when token user differs from path user", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, `{`), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 400 for a validation error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updateFn: func(string, entity.User) error {
			return errorType.NewErrorUserValidation("invalid")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updateFn: func(string, entity.User) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updateFn: func(string, entity.User) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPut, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdateUser(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"userId": testUserID})
		c.DeleteUser(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 403 when token user differs from path user", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.DeleteUser(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{deleteFn: func(string) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.DeleteUser(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{deleteFn: func(string) error { return errors.New("db down") }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodDelete, ""), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.DeleteUser(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestFollow(t *testing.T) {
	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": otherUserID})
		c.Follow(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 403 when the operation is denied", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{followFn: func(string, string) error {
			return usecase.ErrOperationDenied
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.Follow(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{followFn: func(string, string) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.Follow(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{followFn: func(string, string) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.Follow(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestUnfollow(t *testing.T) {
	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": otherUserID})
		c.Unfollow(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 403 when the operation is denied", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{unfollowFn: func(string, string) error {
			return usecase.ErrOperationDenied
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.Unfollow(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{unfollowFn: func(string, string) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, ""), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.Unfollow(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})
}

func TestGetFollowersAndFollowing(t *testing.T) {
	t.Run("GetFollowers should return 200 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getFollowersFn: func(string) ([]entity.User, error) {
			return []entity.User{}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetFollowers(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("GetFollowers should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getFollowersFn: func(string) ([]entity.User, error) {
			return nil, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetFollowers(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("GetFollowing should return 200 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getFollowingFn: func(string) ([]entity.User, error) {
			return []entity.User{}, nil
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetFollowing(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	})

	t.Run("GetFollowing should return 500 on error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{getFollowingFn: func(string) ([]entity.User, error) {
			return nil, errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodGet, ""), map[string]string{"userId": testUserID})
		c.GetFollowing(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestUpdatePassword(t *testing.T) {
	validBody := `{"current":"123","new":"456"}`

	t.Run("Should return 401 without a token", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": testUserID})
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 403 when token user differs from path user", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": otherUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusForbidden)
	})

	t.Run("Should return 400 for an invalid JSON body", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, `{`), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 401 for a wrong current password", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updatePasswordFn: func(string, dto.Password) error {
			return usecase.ErrWrongPassword
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusUnauthorized)
	})

	t.Run("Should return 400 when the new password is too long", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updatePasswordFn: func(string, dto.Password) error {
			return bcrypt.ErrPasswordTooLong
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("Should return 204 on success", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updatePasswordFn: func(string, dto.Password) error { return nil }})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusNoContent)
	})

	t.Run("Should return 500 on unexpected error", func(t *testing.T) {
		c := NewUserController(&mockUserUseCase{updatePasswordFn: func(string, dto.Password) error {
			return errors.New("db down")
		}})
		rec := httptest.NewRecorder()
		req := withVars(newRequest(http.MethodPost, validBody), map[string]string{"userId": testUserID})
		req = withAuth(req, testUserID)
		c.UpdatePassword(rec, req)
		assertStatus(t, rec.Code, http.StatusInternalServerError)
	})
}
