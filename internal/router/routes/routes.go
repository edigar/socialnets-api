package routes

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/controller"
	"github.com/edigar/socialnets-api/internal/middleware"
	"github.com/edigar/socialnets-api/internal/repository"
	"github.com/edigar/socialnets-api/internal/usecase"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	URI                    string
	Method                 string
	Function               func(w http.ResponseWriter, r *http.Request)
	AuthenticationRequired bool
}

func Setup(r *mux.Router, db *sql.DB) *mux.Router {
	userController := controller.NewUserController(usecase.NewUserUseCase(repository.NewUserRepository(db)))
	postController := controller.NewPostController(usecase.NewPostUseCase(repository.NewPostRepository(db)))

	routes := userRoutes(userController)
	routes = append(routes, loginRoute(userController))
	routes = append(routes, postRoutes(postController)...)
	routes = append(routes, healthRoute)

	for _, route := range routes {
		if route.AuthenticationRequired {
			r.HandleFunc(route.URI, middleware.Logger(middleware.Authenticate(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middleware.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}
