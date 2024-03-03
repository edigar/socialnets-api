package usecase

import (
	"database/sql"
	"errors"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/error_type"
	"github.com/edigar/socialnets-api/internal/usecase/mock"
	"reflect"
	"testing"
)

func TestCreatePost(t *testing.T) {
	t.Run("Should create a post with validated data", func(t *testing.T) {
		post := entity.Post{
			Title:   "Title 1",
			Content: "Content 1",
		}

		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.CreatePost(&post)
		if err != nil {
			t.Errorf("CreatePost should not return an error for a valid post data. Post: %v. Error: %v", post, err)
		} else if post.Id != usecase.NEW_POST_ID {
			t.Errorf("CreatePost should set an id for user. User id: %v. User id expected: %v",
				post.Id,
				usecase.NEW_POST_ID,
			)
		}
	})

	t.Run("Should not create a post with non-validated data", func(t *testing.T) {
		scenarios := []entity.Post{
			{
				Title: "Title a",
			},
			{
				Content: "Content b",
			},
		}

		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		for _, scenario := range scenarios {
			err := postUseCase.CreatePost(&scenario)
			var epv *errorType.ErrorPostValidation
			if !errors.As(err, &epv) {
				t.Errorf("CreatePost should return an ErrorPostValidation error for a non-valid post data. Post: %v Returned: %v. Error expected: %T",
					scenario,
					err,
					epv.Error(),
				)
			} else if scenario.Id != 0 {
				t.Errorf("CreatePost should set 0 on post id for non-valid post data. Got: %v.", scenario.Id)
			}
		}
	})
}

func TestGetByUser(t *testing.T) {
	t.Run("Should get all posts of user", func(t *testing.T) {
		userId := usecase.MockPosts[1].AuthorId
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		posts, err := postUseCase.GetByUser(userId)
		if err != nil {
			t.Errorf("GetByUser should not return an error for a valid user id. User: %v. Error: %v", userId, err)
		}

		if !reflect.DeepEqual(posts, usecase.MockPosts) {
			t.Errorf("GetByUser should return posts for user. Expected: %v. Got: %v", posts, usecase.MockPosts)
		}
	})

	t.Run("Should get a bad connection database error", func(t *testing.T) {
		userId := usecase.POST_ERROR
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		posts, err := postUseCase.GetByUser(userId)
		if err.Error() != "driver: bad connection" {
			t.Errorf("GetByUser should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}

		if posts != nil {
			t.Errorf("GetByUser should not get any post if connection get an error. Posts: %v.", posts)
		}
	})
}

func TestGetById(t *testing.T) {
	t.Run("Should get a post by id", func(t *testing.T) {
		postId := usecase.MockPosts[0].Id
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		post, err := postUseCase.GetById(postId)
		if err != nil {
			t.Errorf("GetById should not return an error for a valid id. Post id: %v. Error: %v", postId, err)
		}

		if post != usecase.MockPosts[0] {
			t.Errorf("GetByUser should return post by id. Expected: %v. Got: %v", usecase.MockPosts[0].Id, post)
		}
	})

	t.Run("Should get empty Post if id is invalid", func(t *testing.T) {
		postId := uint64(99)
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		post, err := postUseCase.GetById(postId)
		if err != nil {
			t.Errorf("GetById should not return an error for a non-valid id. Post id: %v. Error: %v", postId, err)
		}

		if post != (entity.Post{}) {
			t.Errorf("GetByUser should return empty post for a non-valid id. Expected: %v. Got: %v", entity.Post{}, post)
		}
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("Should update post with valid id", func(t *testing.T) {
		post := entity.Post{Title: "Title test", Content: "Content test"}
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Update(usecase.MockPosts[0].AuthorId, usecase.MockPosts[0].Id, post)
		if err != nil {
			t.Errorf("Update should not return an error with valid data. Data sended: %v. Post updated: %v Error: %v",
				post,
				usecase.MockPosts[0],
				err,
			)
		} else if post.Title != usecase.MockPosts[0].Title || post.Content != usecase.MockPosts[0].Content {
			t.Errorf("Update should update title and content of mock post 0. Data sended: %v. User updated: %v",
				post,
				usecase.MockUsers[0],
			)
		}
	})

	t.Run("Should not update user with non-valid post id", func(t *testing.T) {
		post := entity.Post{Title: "Title test", Content: "Content test"}
		originalPosts := usecase.MockPosts
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Update(usecase.MockUsers[0].Id, 0, post)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Update should return sql.ErrNoRows error with non-valid post id. Data sended: %v. User updated: %v Error: %v",
				post,
				usecase.MockPosts[0],
				err,
			)
		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
			t.Errorf("Update should not update none of posts with non-valid id. Data sended: %v. Posts: %v", post, usecase.MockPosts)
		}
	})

	t.Run("Should not update post with non-valid author id", func(t *testing.T) {
		post := entity.Post{Title: "Title test", Content: "Content test"}
		originalPosts := usecase.MockPosts
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Update("wrong-author-id", usecase.MockPosts[0].Id, post)
		if !errors.Is(err, ErrAccessDenied) {
			t.Errorf("Update should return ErrAccessDenied error with non-valid author id. Data sended: %v. Post updated: %v Error: %v",
				post,
				usecase.MockPosts[0],
				err,
			)
		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
			t.Errorf("Update should not update none of posts with non-valid author id. Data sended: %v. Posts: %v",
				post,
				usecase.MockPosts,
			)
		}
	})

	t.Run("Should not update post with non-valid data", func(t *testing.T) {
		scenarios := []entity.Post{
			{
				Title: "Title a",
			},
			{
				Content: "Content a",
			},
		}

		originalPosts := usecase.MockPosts
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		for _, scenario := range scenarios {
			err := postUseCase.Update(usecase.MockUsers[0].Id, usecase.MockPosts[0].Id, scenario)
			var euv *errorType.ErrorPostValidation
			if !errors.As(err, &euv) {
				t.Errorf("Update should return an ErrorPostValidation error for a non-valid post data. Post: %v Error returned: %v. Error expected: %T",
					scenario,
					err,
					euv.Error(),
				)
			} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
				t.Errorf("Update should not update none of posts with non-valid data. Data sended: %v. Original posts: %v", scenario, originalPosts)
			}
		}
	})
}

func TestGetUserPosts(t *testing.T) {
	t.Run("Should get user posts by id", func(t *testing.T) {
		userId := usecase.MockPosts[1].AuthorId
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		posts, err := postUseCase.GetUserPosts(userId)
		if err != nil {
			t.Errorf("GetByUser should not return an error for a valid user id. User: %v. Error: %v", userId, err)
		}

		if len(posts) != 2 {
			t.Errorf("GetByUser should return all posts of user. User: %v. Posts: %v", userId, posts)
		}

		for _, post := range posts {
			if post != usecase.MockPosts[1] && post != usecase.MockPosts[2] {
				t.Errorf("GetByUser should return only posts of user. Post: %v. Posts: %v", post, posts)
			}
		}
	})

	t.Run("Should get a bad connection database error", func(t *testing.T) {
		userId := usecase.POST_ERROR
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		posts, err := postUseCase.GetUserPosts(userId)
		if err.Error() != "driver: bad connection" {
			t.Errorf("GetByUser should get a bad connection error. Expected: %v. Got: %v", "driver: bad connection", err)
		}

		if posts != nil {
			t.Errorf("GetByUser should not get any post if connection get an error. Posts: %v.", posts)
		}
	})
}

func TestLikePost(t *testing.T) {
	t.Run("Should like a post", func(t *testing.T) {
		postId := usecase.MockPosts[0].Id
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.LikePost(postId)
		if err != nil {
			t.Errorf("LikePost should not return error for a valid post. Error: %v", err)
		}

		if usecase.MockPosts[0].Likes != 1 {
			t.Errorf("LikePost should increase like to 1. Got: %v", usecase.MockPosts[0].Likes)
		}

		usecase.MockPosts[0].Likes = 0
	})

	t.Run("Should not like a post with non-valid id", func(t *testing.T) {
		postId := uint64(99)
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.LikePost(postId)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("LikePost should return sql.ErrNoRows error with non-valid id. Error: %v", err)
		}
	})
}

func TestUnLikePost(t *testing.T) {
	t.Run("Should unlike a post", func(t *testing.T) {
		usecase.MockPosts[0].Likes = 2
		postId := usecase.MockPosts[0].Id
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.UnLikePost(postId)
		if err != nil {
			t.Errorf("UnLikePost should not return error for a valid post. Error: %v", err)
		}

		if usecase.MockPosts[0].Likes != 1 {
			t.Errorf("LikePost should decrease like to 1. Got: %v", usecase.MockPosts[0].Likes)
		}

		usecase.MockPosts[0].Likes = 0
	})

	t.Run("Should not unlike a post with non-valid id", func(t *testing.T) {
		postId := uint64(99)
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.UnLikePost(postId)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("UnLikePost should return sql.ErrNoRows error with non-valid id. Error: %v", err)
		}
	})
}

func TestDeletePost(t *testing.T) {
	t.Run("Should delete post by id", func(t *testing.T) {
		originalPosts := usecase.MockPosts
		postId := usecase.MockPosts[0].Id
		expectedPosts := []entity.Post{usecase.MockPosts[1], usecase.MockPosts[2]}
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Delete(postId, usecase.MockPosts[0].AuthorId)
		if err != nil {
			t.Errorf("Delete should not return an error for a valid post id and author id. Post: %v. Error: %v", usecase.MockPosts[0], err)
		}
		if !reflect.DeepEqual(expectedPosts, usecase.MockPosts) {
			t.Errorf("Delete should remove post with id %v. Expected: %v. Got: %v", postId, expectedPosts, usecase.MockPosts)
		}

		usecase.MockPosts = originalPosts
	})

	t.Run("Should return an error if post id doesn't exist", func(t *testing.T) {
		var postId uint64
		postId = 999
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Delete(postId, usecase.MockPosts[0].AuthorId)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Delete should return an error for invalid post id. Expected: %v. Got: %v", sql.ErrNoRows, err)
		}
	})

	t.Run("Should not delete post with non-valid author id", func(t *testing.T) {
		originalPosts := usecase.MockPosts
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Delete(usecase.MockPosts[0].Id, "wrong-author-id")
		if !errors.Is(err, ErrAccessDenied) {
			t.Errorf("Delete should return ErrAccessDenied error with non-valid author id. Post to delete: %v Error: %v",
				usecase.MockPosts[0],
				err,
			)
		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
			t.Errorf("Delete should not delete with non-valid author id. Posts: %v", usecase.MockPosts)
		}
	})

	t.Run("Should not delete post with non-valid id", func(t *testing.T) {
		originalPosts := usecase.MockPosts
		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
		err := postUseCase.Delete(0, usecase.MockPosts[0].AuthorId)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Delete should return sql.ErrNoRows error with non-valid id. Error: %v", err)
		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
			t.Errorf("Delete should not delete with non-valid id. Posts: %v", usecase.MockPosts)
		}
	})
}
