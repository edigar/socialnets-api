package routes

import (
	"github.com/edigar/socialnets-api/internal/controllers"
	"net/http"
)

var healthRoute = Route{
	URI:                    "/health",
	Method:                 http.MethodGet,
	Function:               controllers.HealthCheck,
	AuthenticationRequired: false,
}
