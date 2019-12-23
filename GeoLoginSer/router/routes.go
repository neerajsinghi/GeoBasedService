package router

import (
	"net/http"

	handler "GeoLoginSer/handlers"
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
		"register",
		"POST",
		"/register",
		handler.Register,
	},
	Route{
		"login",
		"POST",
		"/login",
		handler.Login,
	},
	Route{
		"checkuser",
		"POST",
		"/checkuser",
		handler.CheckUser,
	},
	Route{
		"change",
		"POST",
		"/change",
		handler.ChangePassword,
	},
	Route{
		"forgot",
		"POST",
		"/forgot",
		handler.ForgotPassword,
	},
	Route{
		"logout",
		"GET",
		"/logout",
		handler.Logout,
	},
}
