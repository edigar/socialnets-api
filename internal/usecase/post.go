package usecase

import (
	"errors"
	"github.com/edigar/socialnets-api/internal/entity"
	"github.com/edigar/socialnets-api/internal/repository"
)

var ErrAccessDenied = errors.New("access denied")

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

func (p *PostUseCase) GetById(postId uint64) (entity.Post, error) {
	post, err := p.postRepository.FetchById(postId)
	if err != nil {
		return entity.Post{}, nil
	}

	return post, nil
}
func (p *PostUseCase) Update(authorId string, postId uint64, post entity.Post) error {
	if err := post.Prepare(); err != nil {
		return err
	}
	postDb, err := p.postRepository.FetchById(postId)
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
	postDb, err := p.postRepository.FetchById(postId)
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

func (p *PostUseCase) GetUserPosts(userId string) ([]entity.Post, error) {
	posts, err := p.postRepository.FetchUserPosts(userId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostUseCase) LikePost(postId uint64) error {
	if err := p.postRepository.LikePost(postId); err != nil {
		return err
	}

	return nil
}

func (p *PostUseCase) UnLikePost(postId uint64) error {
	if err := p.postRepository.UnlikePost(postId); err != nil {
		return err
	}

	return nil
}
