package routes

import (
	"github.com/edigar/socialnets-api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:                    "/api/user",
		Method:                 http.MethodPost,
		Function:               controllers.PostUser,
		AuthenticationRequired: false,
	},
	{
		URI:                    "/api/user",
		Method:                 http.MethodGet,
		Function:               controllers.GetUsers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodGet,
		Function:               controllers.GetUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdateUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeleteUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/follow",
		Method:                 http.MethodPost,
		Function:               controllers.Follow,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/unfollow",
		Method:                 http.MethodPost,
		Function:               controllers.Unfollow,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/followers",
		Method:                 http.MethodGet,
		Function:               controllers.GetFollowers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/following",
		Method:                 http.MethodGet,
		Function:               controllers.GetFollowing,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/update-password",
		Method:                 http.MethodPost,
		Function:               controllers.UpdatePassword,
		AuthenticationRequired: true,
	},
}
