package usecase

import (
	"database/sql"
	"errors"
	"github.com/edigar/socialnets-api/internal/entity"
)

type MockPostRepository struct{}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{}
}

const NEW_POST_ID = 4
const ERROR = "ERROR"

var MockPosts = []entity.Post{
	{
		Id:       1,
		Title:    "Title 1",
		Content:  "Content 1",
		AuthorId: "93226a19-86d6-4ad7-a215-d5999c2870c4",
		Likes:    0,
	},
	{
		Id:       2,
		Title:    "Title 2",
		Content:  "Content 2",
		AuthorId: "d9b56fd4-31b7-4bd5-958f-99028ca5e79a",
		Likes:    0,
	},
	{
		Id:       3,
		Title:    "Title 3",
		Content:  "Content 3",
		AuthorId: "d9b56fd4-31b7-4bd5-958f-99028ca5e79a",
		Likes:    0,
	},
}

func (mr MockPostRepository) Create(post entity.Post) (uint64, error) {
	return NEW_POST_ID, nil
}

func (mr MockPostRepository) FetchById(postId uint64) (entity.Post, error) {
	for _, post := range MockPosts {
		if post.Id == postId {
			return post, nil
		}
	}

	return entity.Post{}, sql.ErrNoRows
}

func (mr MockPostRepository) FetchByUser(userId string) ([]entity.Post, error) {
	if userId == ERROR {
		return nil, errors.New("driver: bad connection")
	}

	return MockPosts, nil
}

func (mr MockPostRepository) Update(postId uint64, post entity.Post) error {
	for i, mockPost := range MockPosts {
		if mockPost.Id == postId {
			MockPosts[i].Title = post.Title
			MockPosts[i].Content = post.Content
		}
	}

	return nil
}

func (r MockPostRepository) Delete(postId uint64) error {
	index := 99
	for i, post := range MockPosts {
		if postId == post.Id {
			index = i
		}
	}

	if index != 99 {
		MockPosts = append(MockPosts[:index], MockPosts[index+1:]...)
		return nil
	}

	return sql.ErrNoRows
}

func (mr MockPostRepository) FetchUserPosts(userId string) ([]entity.Post, error) {
	if userId == ERROR {
		return nil, errors.New("driver: bad connection")
	}

	var posts []entity.Post
	for _, post := range MockPosts {
		if post.AuthorId == userId {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (mr MockPostRepository) LikePost(postId uint64) error {
	for i, post := range MockPosts {
		if postId == post.Id {
			MockPosts[i].Likes++
			return nil
		}
	}

	return sql.ErrNoRows
}

func (mr MockPostRepository) UnlikePost(postId uint64) error {
	for i, post := range MockPosts {
		if postId == post.Id && post.Likes > 0 {
			MockPosts[i].Likes--
			return nil
		}
	}

	return sql.ErrNoRows
}
