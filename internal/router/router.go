package router

import (
	"github.com/edigar/socialnets-api/internal/router/routes"
	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	return routes.Setup(r)
}
