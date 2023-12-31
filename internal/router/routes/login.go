package routes

import (
	"github.com/edigar/socialnets-api/internal/controllers"
	"net/http"
)

var loginRoute = Route{
	URI:                    "/api/login",
	Method:                 http.MethodPost,
	Function:               controllers.Login,
	AuthenticationRequired: false,
}
