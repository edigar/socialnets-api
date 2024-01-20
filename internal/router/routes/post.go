package routes

import (
	"github.com/edigar/socialnets-api/internal/controller"
	"net/http"
)

var postRoutes = []Route{
	{
		URI:                    "/api/post",
		Method:                 http.MethodPost,
		Function:               controller.PostPost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post",
		Method:                 http.MethodGet,
		Function:               controller.GetPosts,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodGet,
		Function:               controller.GetPost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodPut,
		Function:               controller.UpdatePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodDelete,
		Function:               controller.DeletePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/posts",
		Method:                 http.MethodGet,
		Function:               controller.GetUserPosts,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}/like",
		Method:                 http.MethodPost,
		Function:               controller.LikePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}/unlike",
		Method:                 http.MethodPost,
		Function:               controller.UnlikePost,
		AuthenticationRequired: true,
	},
}
