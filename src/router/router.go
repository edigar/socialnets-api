package router

import (
	"github.com/edigar/socialnets-api/src/router/routes"
	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	return routes.Setup(r)
}
