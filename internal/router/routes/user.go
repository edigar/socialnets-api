package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

func userRoutes(c *controller.UserController) []Route {
	return []Route{
		{
			URI:                    "/api/user",
			Method:                 http.MethodPost,
			Function:               c.PostUser,
			AuthenticationRequired: false,
		},
		{
			URI:                    "/api/user",
			Method:                 http.MethodGet,
			Function:               c.GetUsers,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}",
			Method:                 http.MethodGet,
			Function:               c.GetUser,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}",
			Method:                 http.MethodPut,
			Function:               c.UpdateUser,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}",
			Method:                 http.MethodDelete,
			Function:               c.DeleteUser,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/follow",
			Method:                 http.MethodPost,
			Function:               c.Follow,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/unfollow",
			Method:                 http.MethodPost,
			Function:               c.Unfollow,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/followers",
			Method:                 http.MethodGet,
			Function:               c.GetFollowers,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/following",
			Method:                 http.MethodGet,
			Function:               c.GetFollowing,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/update-password",
			Method:                 http.MethodPost,
			Function:               c.UpdatePassword,
			AuthenticationRequired: true,
		},
	}
}
