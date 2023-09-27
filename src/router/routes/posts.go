package routes

import (
	"devbook_api/src/controllers"
	"net/http"
)

var routesPosts = []Route{
	{
		Uri:                   "/posts",
		Method:                http.MethodPost,
		Function:              controllers.CreatePost,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts",
		Method:                http.MethodGet,
		Function:              controllers.GetPosts,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts/{id}",
		Method:                http.MethodGet,
		Function:              controllers.GetPost,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts/{id}",
		Method:                http.MethodPut,
		Function:              controllers.UpdatePost,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts/{id}",
		Method:                http.MethodDelete,
		Function:              controllers.DeletePost,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/posts",
		Method:                http.MethodGet,
		Function:              controllers.GetPostsByUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts/{id}/like",
		Method:                http.MethodPost,
		Function:              controllers.LikePost,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/posts/{id}/dislike",
		Method:                http.MethodPost,
		Function:              controllers.DislikePost,
		RequireAuthentication: true,
	},
}
