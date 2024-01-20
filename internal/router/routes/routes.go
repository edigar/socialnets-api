package routes

import (
	"github.com/edigar/socialnets-api/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	URI                    string
	Method                 string
	Function               func(w http.ResponseWriter, r *http.Request)
	AuthenticationRequired bool
}

func Setup(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, loginRoute)
	routes = append(routes, postRoutes...)
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
