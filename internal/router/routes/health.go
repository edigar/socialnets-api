package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

var healthRoute = Route{
	URI:                    "/health",
	Method:                 http.MethodGet,
	Function:               controller.HealthCheck,
	AuthenticationRequired: false,
}
