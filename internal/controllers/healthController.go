package controllers

import (
	"github.com/edigar/socialnets-api/internal/responses"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
