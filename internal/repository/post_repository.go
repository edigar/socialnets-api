package repository

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/entity"
)

type Post interface {
	Create(post entity.Post) (uint64, error)
	FetchById(postId uint64, currentUserId string) (entity.Post, error)
	FetchByUser(userId string) ([]entity.Post, error)
	Update(postId uint64, post entity.Post) error
	Delete(postId uint64) error
	FetchUserPosts(userId string, currentUserId string) ([]entity.Post, error)
	Exists(postId uint64) (bool, error)
	Like(postId uint64, userId string) error
	Unlike(postId uint64, userId string) error
	FetchLikers(postId uint64) ([]entity.User, error)
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

func (r PostRepository) FetchById(postId uint64, currentUserId string) (entity.Post, error) {
	row, err := r.db.Query(
		`SELECT p.id, p.title, p.content, p.author,
			(SELECT COUNT(*) FROM post_likes pl WHERE pl.post_id = p.id) AS likes,
			p.created_at, u.nick,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $2) AS liked_by_me
		FROM posts p INNER JOIN users u ON u.id = p.author WHERE p.id = $1`,
		postId, currentUserId,
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
			&post.LikedByMe,
		); err != nil {
			return entity.Post{}, err
		}
	}

	return post, nil
}

func (r PostRepository) FetchByUser(userId string) ([]entity.Post, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT p.id, p.title, p.content, p.author,
			(SELECT COUNT(*) FROM post_likes pl WHERE pl.post_id = p.id) AS likes,
			p.created_at, u.nick,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $1) AS liked_by_me
		FROM posts p
		LEFT JOIN users u ON u.id = p.author
		LEFT JOIN followers f ON p.author = f.user_id WHERE u.id = $1 OR f.follower = $1
		ORDER BY 6 desc`,
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
			&post.LikedByMe,
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

func (r PostRepository) FetchUserPosts(userId string, currentUserId string) ([]entity.Post, error) {
	rows, err := r.db.Query(
		`SELECT p.id, p.title, p.content, p.author,
			(SELECT COUNT(*) FROM post_likes pl WHERE pl.post_id = p.id) AS likes,
			p.created_at, u.nick,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = $2) AS liked_by_me
		FROM posts p JOIN users u ON u.id = p.author WHERE p.author = $1`,
		userId, currentUserId,
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
			&post.LikedByMe,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r PostRepository) Exists(postId uint64) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r PostRepository) Like(postId uint64, userId string) error {
	stmt := "INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	_, err := r.db.Exec(stmt, postId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) Unlike(postId uint64, userId string) error {
	stmt := "DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2"
	_, err := r.db.Exec(stmt, postId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostRepository) FetchLikers(postId uint64) ([]entity.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.name, u.nick FROM post_likes pl
		JOIN users u ON u.id = pl.user_id
		WHERE pl.post_id = $1 ORDER BY pl.created_at`,
		postId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(&user.Id, &user.Name, &user.Nick); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
