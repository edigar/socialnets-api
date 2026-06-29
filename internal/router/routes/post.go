package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

func postRoutes(c *controller.PostController) []Route {
	return []Route{
		{
			URI:                    "/api/post",
			Method:                 http.MethodPost,
			Function:               c.PostPost,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post",
			Method:                 http.MethodGet,
			Function:               c.GetPosts,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post/{postId}",
			Method:                 http.MethodGet,
			Function:               c.GetPost,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post/{postId}",
			Method:                 http.MethodPut,
			Function:               c.UpdatePost,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post/{postId}",
			Method:                 http.MethodDelete,
			Function:               c.DeletePost,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/user/{userId}/posts",
			Method:                 http.MethodGet,
			Function:               c.GetUserPosts,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post/{postId}/like",
			Method:                 http.MethodPost,
			Function:               c.LikePost,
			AuthenticationRequired: true,
		},
		{
			URI:                    "/api/post/{postId}/unlike",
			Method:                 http.MethodPost,
			Function:               c.UnlikePost,
			AuthenticationRequired: true,
		},
	}
}
