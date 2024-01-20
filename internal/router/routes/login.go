package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

var loginRoute = Route{
	URI:                    "/api/login",
	Method:                 http.MethodPost,
	Function:               controller.Login,
	AuthenticationRequired: false,
}
