package entity

import (
	errorType "github.com/edigar/socialnets-api/internal/error_type"
	"strings"
	"time"
)

type Post struct {
	Id         uint64    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorId   string    `json:"authorId,omitempty"`
	AuthorNick string    `json:"authorNick,omitempty"`
	Likes      uint64    `json:"likes"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

func (post *Post) Prepare() error {
	if err := post.validate(); err != nil {
		return err
	}
	post.format()

	return nil
}

func (post *Post) validate() error {
	if post.Title == "" {
		return errorType.NewErrorPostValidation("title is required")
	}
	if post.Content == "" {
		return errorType.NewErrorPostValidation("content is required")
	}

	return nil
}

func (post *Post) format() {
	post.Title = strings.TrimSpace(post.Title)
	post.Content = strings.TrimSpace(post.Content)
}
