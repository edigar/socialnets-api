package usecase

import (
	"errors"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/repository"
)

var (
	ErrAccessDenied = errors.New("access denied")
	ErrPostNotFound = errors.New("post not found")
)

type PostUseCase struct {
	postRepository repository.Post
}

func NewPostUseCase(postRepository repository.Post) *PostUseCase {
	return &PostUseCase{
		postRepository: postRepository,
	}
}

func (p *PostUseCase) CreatePost(post *entity.Post) error {
	err := post.Prepare()
	if err != nil {
		return err
	}

	post.Id, err = p.postRepository.Create(*post)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostUseCase) GetByUser(userId string) ([]entity.Post, error) {
	posts, err := p.postRepository.FetchByUser(userId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostUseCase) GetById(postId uint64, currentUserId string) (entity.Post, error) {
	post, err := p.postRepository.FetchById(postId, currentUserId)
	if err != nil {
		return entity.Post{}, nil
	}

	return post, nil
}
func (p *PostUseCase) Update(authorId string, postId uint64, post entity.Post) error {
	if err := post.Prepare(); err != nil {
		return err
	}
	postDb, err := p.postRepository.FetchById(postId, authorId)
	if err != nil {
		return err
	}
	if postDb.AuthorId != authorId {
		return ErrAccessDenied
	}

	if err = p.postRepository.Update(postId, post); err != nil {
		return err
	}

	return nil
}

func (p *PostUseCase) Delete(postId uint64, authorId string) error {
	postDb, err := p.postRepository.FetchById(postId, authorId)
	if err != nil {
		return err
	}
	if postDb.AuthorId != authorId {
		return ErrAccessDenied
	}
	if err = p.postRepository.Delete(postId); err != nil {
		return err
	}

	return nil
}

func (p *PostUseCase) GetUserPosts(userId string, currentUserId string) ([]entity.Post, error) {
	posts, err := p.postRepository.FetchUserPosts(userId, currentUserId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostUseCase) LikePost(postId uint64, userId string) error {
	exists, err := p.postRepository.Exists(postId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPostNotFound
	}

	return p.postRepository.Like(postId, userId)
}

func (p *PostUseCase) UnLikePost(postId uint64, userId string) error {
	return p.postRepository.Unlike(postId, userId)
}

func (p *PostUseCase) GetLikers(postId uint64) ([]entity.User, error) {
	return p.postRepository.FetchLikers(postId)
}
