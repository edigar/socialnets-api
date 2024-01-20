package controller

import (
	"github.com/edigar/socialnets-api/internal/response"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
