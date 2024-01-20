package middleware

import (
	"github.com/edigar/socialnets-api/internal/authentication"
	"github.com/edigar/socialnets-api/internal/response"
	"log"
	"net/http"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" %s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authentication.TokenValidate(r); err != nil {
			response.Error(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r)
	}
}
