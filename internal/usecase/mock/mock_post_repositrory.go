package usecase

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/entity"
)

type MockPostRepository struct{}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{}
}

const NEW_POST_ID = 3

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
	return nil, nil
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
	return nil
}

func (mr MockPostRepository) FetchUserPosts(userId string) ([]entity.Post, error) {
	return nil, nil
}

func (mr MockPostRepository) LikePost(postId uint64) error {
	return nil
}

func (mr MockPostRepository) UnlikePost(postId uint64) error {
	return nil
}
