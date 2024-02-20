package routes

import (
	"devbook_api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		Uri:                   "/users",
		Method:                http.MethodPost,
		Function:              controllers.CreateUser,
		RequireAuthentication: false,
	},
	{
		Uri:                   "/users",
		Method:                http.MethodGet,
		Function:              controllers.GetUsers,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}",
		Method:                http.MethodGet,
		Function:              controllers.GetUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}",
		Method:                http.MethodPut,
		Function:              controllers.UpdateUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}",
		Method:                http.MethodDelete,
		Function:              controllers.DeleteUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/follow",
		Method:                http.MethodPost,
		Function:              controllers.FollowUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/unfollow",
		Method:                http.MethodPost,
		Function:              controllers.UnfollowUser,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/followers",
		Method:                http.MethodGet,
		Function:              controllers.GetFollowers,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/followings",
		Method:                http.MethodGet,
		Function:              controllers.GetFollowings,
		RequireAuthentication: true,
	},
	{
		Uri:                   "/users/{id}/update-password",
		Method:                http.MethodPost,
		Function:              controllers.UpdatePassword,
		RequireAuthentication: true,
	},
}
