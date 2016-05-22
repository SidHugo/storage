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
		"Create sign",
		"POST",
		"/signs",
		CreateSign,
	},
	Route{
		"Get sign",
		"GET",
		"/signs/{signName}",
		GetSign,
	},
	Route{
		"Get signs",
		"GET",
		"/signs",
		GetSigns,
	},
	Route{
		"Delete sign",
		"DELETE",
		"/signs/{signName}",
		DeleteSign,
	},
	Route{
		"Create user",
		"POST",
		"/users",
		CreateUser,
	},
}
