package controller

import (
	"github.com/edigar/socialnets-api/internal/dto"
	"github.com/edigar/socialnets-api/internal/entity"
)

// UserUseCase is the set of user operations the controllers depend on.
type UserUseCase interface {
	Login(email string, password string) (string, error)
	Register(user *entity.User) error
	GetById(id string) (entity.User, error)
	GetByNameOrNick(nameOrNick string) ([]entity.User, error)
	Update(userId string, user entity.User) error
	Delete(userId string) error
	Follow(userId string, follower string) error
	Unfollow(userId string, follower string) error
	GetFollowers(userId string) ([]entity.User, error)
	GetFollowing(userId string) ([]entity.User, error)
	UpdatePassword(userId string, password dto.Password) error
}

// PostUseCase is the set of post operations the controllers depend on.
type PostUseCase interface {
	CreatePost(post *entity.Post) error
	GetByUser(userId string) ([]entity.Post, error)
	GetById(postId uint64) (entity.Post, error)
	Update(authorId string, postId uint64, post entity.Post) error
	Delete(postId uint64, authorId string) error
	GetUserPosts(userId string) ([]entity.Post, error)
	LikePost(postId uint64) error
	UnLikePost(postId uint64) error
}

type UserController struct {
	userUseCase UserUseCase
}

func NewUserController(userUseCase UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

type PostController struct {
	postUseCase PostUseCase
}

func NewPostController(postUseCase PostUseCase) *PostController {
	return &PostController{postUseCase: postUseCase}
}
