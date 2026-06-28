package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/gorilla/mux"
)

// mockUserUseCase implements controller.UserUseCase with configurable behavior per test.
type mockUserUseCase struct {
	loginFn           func(email string, password string) (string, error)
	registerFn        func(user *entity.User) error
	getByIdFn         func(id string) (entity.User, error)
	getByNameOrNickFn func(nameOrNick string) ([]entity.User, error)
	updateFn          func(userId string, user entity.User) error
	deleteFn          func(userId string) error
	followFn          func(userId string, follower string) error
	unfollowFn        func(userId string, follower string) error
	getFollowersFn    func(userId string) ([]entity.User, error)
	getFollowingFn    func(userId string) ([]entity.User, error)
	updatePasswordFn  func(userId string, password dto.Password) error
}

func (m *mockUserUseCase) Login(email string, password string) (string, error) {
	return m.loginFn(email, password)
}
func (m *mockUserUseCase) Register(user *entity.User) error       { return m.registerFn(user) }
func (m *mockUserUseCase) GetById(id string) (entity.User, error) { return m.getByIdFn(id) }
func (m *mockUserUseCase) GetByNameOrNick(nameOrNick string) ([]entity.User, error) {
	return m.getByNameOrNickFn(nameOrNick)
}
func (m *mockUserUseCase) Update(userId string, user entity.User) error {
	return m.updateFn(userId, user)
}
func (m *mockUserUseCase) Delete(userId string) error { return m.deleteFn(userId) }
func (m *mockUserUseCase) Follow(userId string, follower string) error {
	return m.followFn(userId, follower)
}
func (m *mockUserUseCase) Unfollow(userId string, follower string) error {
	return m.unfollowFn(userId, follower)
}
func (m *mockUserUseCase) GetFollowers(userId string) ([]entity.User, error) {
	return m.getFollowersFn(userId)
}
func (m *mockUserUseCase) GetFollowing(userId string) ([]entity.User, error) {
	return m.getFollowingFn(userId)
}
func (m *mockUserUseCase) UpdatePassword(userId string, password dto.Password) error {
	return m.updatePasswordFn(userId, password)
}

// mockPostUseCase implements controller.PostUseCase with configurable behavior per test.
type mockPostUseCase struct {
	createPostFn   func(post *entity.Post) error
	getByUserFn    func(userId string) ([]entity.Post, error)
	getByIdFn      func(postId uint64) (entity.Post, error)
	updateFn       func(authorId string, postId uint64, post entity.Post) error
	deleteFn       func(postId uint64, authorId string) error
	getUserPostsFn func(userId string) ([]entity.Post, error)
	likePostFn     func(postId uint64) error
	unLikePostFn   func(postId uint64) error
}

func (m *mockPostUseCase) CreatePost(post *entity.Post) error { return m.createPostFn(post) }
func (m *mockPostUseCase) GetByUser(userId string) ([]entity.Post, error) {
	return m.getByUserFn(userId)
}
func (m *mockPostUseCase) GetById(postId uint64) (entity.Post, error) { return m.getByIdFn(postId) }
func (m *mockPostUseCase) Update(authorId string, postId uint64, post entity.Post) error {
	return m.updateFn(authorId, postId, post)
}
func (m *mockPostUseCase) Delete(postId uint64, authorId string) error {
	return m.deleteFn(postId, authorId)
}
func (m *mockPostUseCase) GetUserPosts(userId string) ([]entity.Post, error) {
	return m.getUserPostsFn(userId)
}
func (m *mockPostUseCase) LikePost(postId uint64) error   { return m.likePostFn(postId) }
func (m *mockPostUseCase) UnLikePost(postId uint64) error { return m.unLikePostFn(postId) }

// --- request helpers ---

func newRequest(method, body string) *http.Request {
	if body == "" {
		return httptest.NewRequest(method, "/", nil)
	}
	return httptest.NewRequest(method, "/", strings.NewReader(body))
}

// withAuth attaches a valid Bearer token whose userId claim is the given id.
func withAuth(r *http.Request, userId string) *http.Request {
	token, _ := authentication.CreateToken(userId)
	r.Header.Set("Authorization", "Bearer "+token)
	return r
}

func withVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}
