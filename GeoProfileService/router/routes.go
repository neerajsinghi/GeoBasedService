package router

import (
	"net/http"

	handler "GeoProfileService/handlers"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

var routes = Routes{
	Route{
		"getpeople",
		"POST",
		"/getpeople",
		handler.GetPeople,
	},
	Route{
		"searchpeople",
		"POST",
		"/searchpeople",
		handler.SearchPeople,
	},
	Route{
		"updatelocation",
		"POST",
		"/updatelocation",
		handler.UpdateLocation,
	},
	Route{
		"updateprofile",
		"POST",
		"/updateprofile",
		handler.UpdateProfile,
	},
	Route{
		"uploademergency",
		"POST",
		"/uploademergency",
		handler.UploadEmergency,
	},

	Route{
		"uploadimage",
		"POST",
		"/uploadimage",
		handler.UploadImage,
	},
	Route{
		"getprofile",
		"POST",
		"/getprofile",
		handler.GetProfile,
	},

	Route{
		"getuserdata",
		"POST",
		"/getuserdata",
		handler.GetUserData,
	},
}
