package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

var userRoutes = []Route{
	{
		URI:                    "/api/user",
		Method:                 http.MethodPost,
		Function:               controller.PostUser,
		AuthenticationRequired: false,
	},
	{
		URI:                    "/api/user",
		Method:                 http.MethodGet,
		Function:               controller.GetUsers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodGet,
		Function:               controller.GetUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodPut,
		Function:               controller.UpdateUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}",
		Method:                 http.MethodDelete,
		Function:               controller.DeleteUser,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/follow",
		Method:                 http.MethodPost,
		Function:               controller.Follow,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/unfollow",
		Method:                 http.MethodPost,
		Function:               controller.Unfollow,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/followers",
		Method:                 http.MethodGet,
		Function:               controller.GetFollowers,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/following",
		Method:                 http.MethodGet,
		Function:               controller.GetFollowing,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/update-password",
		Method:                 http.MethodPost,
		Function:               controller.UpdatePassword,
		AuthenticationRequired: true,
	},
}
