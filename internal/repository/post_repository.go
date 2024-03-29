package repository

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/entity"
)

type Post interface {
	Create(post entity.Post) (uint64, error)
	FetchById(postId uint64) (entity.Post, error)
	FetchByUser(userId string) ([]entity.Post, error)
	Update(postId uint64, post entity.Post) error
	Delete(postId uint64) error
	FetchUserPosts(userId string) ([]entity.Post, error)
	LikePost(postId uint64) error
	UnlikePost(postId uint64) error
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db}
}

func (r PostRepository) Create(post entity.Post) (uint64, error) {
	var postId uint64
	insertStmt := `INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(insertStmt, post.Title, post.Content, post.AuthorId).Scan(&postId)
	if err != nil {
		return 0, err
	}

	return postId, nil
}

func (r PostRepository) FetchById(postId uint64) (entity.Post, error) {
	row, err := r.db.Query(
		"SELECT p.*, u.nick FROM posts p INNER JOIN users u ON u.id = p.author WHERE p.id = $1",
		postId,
	)
	if err != nil {
		return entity.Post{}, err
	}
	defer row.Close()

	var post entity.Post
	if row.Next() {
		if err := row.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return entity.Post{}, err
		}
	}

	return post, nil
}

func (r PostRepository) FetchByUser(userId string) ([]entity.Post, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT p.*, u.nick FROM posts p
		LEFT JOIN users u ON u.id = p.author
		LEFT JOIN followers f ON p.author = f.user_id WHERE u.id = $1 OR f.follower = $1
		ORDER BY 1 desc`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		var post entity.Post
		if err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r PostRepository) Update(postId uint64, post entity.Post) error {
	updateStmt := "UPDATE posts SET title=$1, content=$2 WHERE id=$3"
	_, err := r.db.Exec(updateStmt, post.Title, post.Content, postId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) Delete(postId uint64) error {
	deleteStmt := "DELETE FROM posts WHERE id=$1"
	_, err := r.db.Exec(deleteStmt, postId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) FetchUserPosts(userId string) ([]entity.Post, error) {
	rows, err := r.db.Query(
		"SELECT p.*, u.nick FROM posts p JOIN users u ON u.id = p.author WHERE p.author = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post

	for rows.Next() {
		var post entity.Post
		if err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.AuthorId,
			&post.Likes,
			&post.CreatedAt,
			&post.AuthorNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r PostRepository) LikePost(postId uint64) error {
	updateStmt := "UPDATE posts SET likes = likes + 1 WHERE id=$1"
	_, err := r.db.Exec(updateStmt, postId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) UnlikePost(postId uint64) error {
	//updateStmt := "UPDATE posts SET likes = CASE WHEN likes > 0 THEN likes - 1 ELSE 0 END WHERE id=$1"
	updateStmt := "UPDATE posts SET likes = likes - 1 WHERE id=$1 AND likes > 0"
	_, err := r.db.Exec(updateStmt, postId)
	if err != nil {
		return err
	}

	return nil
}
