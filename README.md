# Socialnets-API

A social network API with the functionalities (so far) of user registration, following users, search, posting, and liking posts.

## Requirements

- [PostgreSQL](https://www.postgresql.org/) 16.0
- [Golang](https://go.dev/) 1.21.1

If you use [docker](https://www.docker.com/), here you will find a complete environment container.

## Installation

Clone or download this repository.

Copy `.env.example` to `.env`.

> [!IMPORTANT]
> You must set `SECRET_KEY` in this file, which can literally be any string.

However, you can generate a `SECRET_KEY` with the following Go code, which will print it in the prompt. After that, simply copy the generated hash to the `.env` file.

```Go
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}

	secret := base64.StdEncoding.EncodeToString(key)
	fmt.Println(secret)
}
```

After that, choose one of the following options (Docker or Local).

### Docker

At the root of the project there is a Dockerfile and docker-compose.yml files with requirements to run the project.
You can run with docker by building the image and running it.

Execute:

```bash
docker-compose up
```

You'll have API running on http://localhost

### Local

Set database configuration on `.env` file and create the database. In docker folder, you will found a SQL script that will create database and tables.

```bash
# Install dependencies
go mod tidy

# Start application
go run ./cmd/api/main.go
```

If you want to generate the executable, use:
```bash
# Install dependencies
go mod tidy

# Generate executable
go build ./cmd/api

# Start application
./api
```

You'll have API running on http://localhost:8000 (if you don't change `API_PORT` on `.env`).

## Usage

After starting application, you'll have access to following routes:

| Method | URI                                | Authentication | Description                             |
|:------:|------------------------------------|:--------------:|-----------------------------------------|
|  GET   | /health                            |       No       | Application status health check         |
|  POST  | /api/login                         |       No       | User login                              |
|  POST  | /api/user                          |      Yes       | Create an user                          |
|  GET   | /api/user                          |      Yes       | Search for users                        |
|  GET   | /api/user/{userId}                 |      Yes       | Get an user data                        |
|  PUT   | /api/user/{userId}                 |      Yes       | Update an user data                     |
| DELETE | /api/user/{userId}                 |      Yes       | Delete an user                          |
|  POST  | /api/user/{userId}/follow          |      Yes       | Logged user follows an user             |
|  POST  | /api/user/{userId}/unfollow        |      Yes       | Logged user unfollows an user           |
|  GET   | /api/user/{userId}/followers       |      Yes       | Get all followers by an user            |
|  GET   | /api/user/{userId}/following       |      Yes       | Gets all users who are following a user |
|  POST  | /api/user/{userId}/update-password |      Yes       | Update password user                    |
|  POST  | /api/post                          |      Yes       | Create a post                           |
|  GET   | /api/post                          |      Yes       | Get all post for a logged user          |
|  GET   | /api/post/{postId}                 |      Yes       | Get a post                              |
|  PUT   | /api/post/{postId}                 |      Yes       | Update a post                           |
| DELETE | /api/post/{postId}                 |      Yes       | Delete a post                           |
|  GET   | /api/user/{userId}/posts           |      Yes       | Gets all posts from a user              |
|  POST  | /api/post/{postId}/like            |      Yes       | Like a user post                        |
|  POST  | /api/post/{postId}/unlike          |      Yes       | Unlike a user post                      |

Authentication, once performed with login (email) and password on `/api/login`, is maintained via a JWT Bearer Token.

You can find more information, like payload and responses in the application swagger (coming soon).

So API is already to use.

## Tests

To run tests run on root of the project:

```bash
go test ./...
```

### Coverage

Run at the root of the project to see tests coverage:
```bash
go test ./... --cover
```

## API documentation

:construction:

Soon a [swagger](https://swagger.io/docs/specification/2-0/what-is-swagger/) doc will be added.

## TODO

- [ ] Add [GORM](https://gorm.io) to the project
- [ ] Improve tests
- [ ] Add [Swagger](https://swagger.io) documentation
- [x] Implements UUID for user's Id
- [ ] Register who liked a post, so that each user can only like each post once, in addition to having information on who liked each post. 
    + [ ] Add likes table
    + [ ] Control who like or unlike a post.
- [ ] Implements migrations
- [ ] Password recovery

## License

[MIT](./LICENSE)