package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

func loginRoute(c *controller.UserController) Route {
	return Route{
		URI:                    "/api/login",
		Method:                 http.MethodPost,
		Function:               c.Login,
		AuthenticationRequired: false,
	}
}
