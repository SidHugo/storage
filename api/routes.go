package api

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Ping,
	},
	Route{
		"Signs",
		"POST",
		"/signs",
		CreateSign,
	},
	Route{
		"Sign",
		"GET",
		"/signs/{signName}",
		GetSign,
	},
}
