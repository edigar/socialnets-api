package routes

import (
	"github.com/edigar/socialnets-api/src/controllers"
	"net/http"
)

var postRoutes = []Route{
	{
		URI:                    "/api/post",
		Method:                 http.MethodPost,
		Function:               controllers.PostPost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post",
		Method:                 http.MethodGet,
		Function:               controllers.GetPosts,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodGet,
		Function:               controllers.GetPost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodPut,
		Function:               controllers.UpdatePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}",
		Method:                 http.MethodDelete,
		Function:               controllers.DeletePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/user/{userId}/posts",
		Method:                 http.MethodGet,
		Function:               controllers.GetUserPosts,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}/like",
		Method:                 http.MethodPost,
		Function:               controllers.LikePost,
		AuthenticationRequired: true,
	},
	{
		URI:                    "/api/post/{postId}/unlike",
		Method:                 http.MethodPost,
		Function:               controllers.UnlikePost,
		AuthenticationRequired: true,
	},
}
