package router

import (
	"net/http"

	handler "UploadService/handlers"
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
		"upload",
		"POST",
		"/upload",
		handler.Upload,
	}, Route{
		"download",
		"GET",
		"/download",
		handler.Download,
	},
}
