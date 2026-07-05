package router

import (
	"database/sql"
	"github.com/edigar/socialnets-api/internal/router/routes"
	"github.com/gorilla/mux"
)

func Generate(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	return routes.Setup(r, db)
}
