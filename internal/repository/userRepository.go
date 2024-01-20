package repository

import (
	"database/sql"
	"fmt"
	"github.com/edigar/socialnets-api/internal/entity"
	"time"
)

type User interface {
	Create(user entity.User) (string, error)
	FetchBy(nameOrNick string) ([]entity.User, error)
	FetchById(userId string) (entity.User, error)
	FetchByEmail(email string) (entity.User, error)
	Update(userId string, user entity.User) error
	Delete(userId string) error
	Follow(userId, follower string) error
	Unfollow(userId, follower string) error
	FetchFollowers(userId string) ([]entity.User, error)
	FetchFollowing(userId string) ([]entity.User, error)
	FetchPasswordById(userId string) (string, error)
	UpdatePassword(userId string, passwordHash string) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) Create(user entity.User) (string, error) {
	var userId string
	insertStmt := `INSERT INTO users (name, nick, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(insertStmt, user.Name, user.Nick, user.Email, user.Password).Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (r UserRepository) FetchBy(nameOrNick string) ([]entity.User, error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)
	rows, err := r.db.Query(
		"SELECT id, name, nick, email, password, created_at, updated_at FROM users WHERE name LIKE $1 OR nick LIKE $1",
		nameOrNick,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r UserRepository) FetchById(userId string) (entity.User, error) {
	row, err := r.db.Query(
		"SELECT id, name, nick, email, password, created_at, updated_at FROM users WHERE id = $1",
		userId,
	)

	if err != nil {
		return entity.User{}, err
	}
	defer row.Close()

	var user entity.User

	if row.Next() {
		if err := row.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return entity.User{}, err
		}
	}

	return user, nil
}

func (r UserRepository) FetchByEmail(email string) (entity.User, error) {
	row, err := r.db.Query("SELECT id, password FROM users WHERE email = $1", email)
	if err != nil {
		return entity.User{}, err
	}
	defer row.Close()

	var user entity.User
	if row.Next() {
		if err := row.Scan(&user.Id, &user.Password); err != nil {
			return entity.User{}, err
		}
	}

	return user, nil
}

func (r UserRepository) Update(userId string, user entity.User) error {
	updateStmt := "UPDATE users SET name=$1, nick=$2, email=$3, updated_at=$4 WHERE id=$5"
	_, err := r.db.Exec(updateStmt, user.Name, user.Nick, user.Email, time.Now(), userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Delete(userId string) error {
	deleteStmt := "DELETE FROM users WHERE id=$1"
	_, err := r.db.Exec(deleteStmt, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Follow(userId, follower string) error {
	insertStmt := "INSERT INTO followers (user_id, follower) VALUES ($1, $2) ON CONFLICT (user_id, follower) DO NOTHING"
	_, err := r.db.Exec(insertStmt, userId, follower)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) Unfollow(userId, follower string) error {
	deleteStmt := "DELETE FROM followers WHERE user_id=$1 AND follower=$2"
	_, err := r.db.Exec(deleteStmt, userId, follower)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) FetchFollowers(userId string) ([]entity.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.created_at, u.updated_at
		FROM users u INNER JOIN followers f ON u.id = f.follower WHERE f.user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r UserRepository) FetchFollowing(userId string) ([]entity.User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.name, u.nick, u.email, u.created_at, u.updated_at
		FROM users u INNER JOIN followers f ON u.id = f.user_id WHERE f.follower = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r UserRepository) FetchPasswordById(userId string) (string, error) {
	row, err := r.db.Query("SELECT password FROM users WHERE id = $1", userId)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var user entity.User

	if row.Next() {
		if err = row.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (r UserRepository) UpdatePassword(userId string, passwordHash string) error {
	updateStmt := "UPDATE users SET password=$1, updated_at=$2 WHERE id=$3"
	_, err := r.db.Exec(updateStmt, passwordHash, time.Now(), userId)
	if err != nil {
		return err
	}

	return nil
}
