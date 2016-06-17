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
		"/signs/{link}",
		GetSign,
	},
	Route{
		"Get sign by JSON parameter map",
		"GET",
		"/signsJson",
		GetSignJson,
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
		"/signs/{link}",
		DeleteSign,
	},
	Route{
		"Create user",
		"POST",
		"/users",
		CreateUser,
	},
	Route{
		"Get user",
		"GET",
		"/users/{key}",
		GetUser,
	},
	Route{
		"Delete user",
		"DELETE",
		"/users/{key}",
		DeleteUser,
	},
	Route{
		"Authorization",
		"GET",
		"/users/authorize/{login}/{password}",
		Authorization,
	},
	Route{
		"Get subscriptions",
		"GET",
		"/users/subscriptions/{login}/{password}",
		GetSubscriptions,
	},
	Route{
		"Get last results",
		"GET",
		"/users/lastresults/{login}/{password}/{requiredLogin}",
		GetLastResults,
	},
	Route{
		"Get users IPs",
		"GET",
		"/users/getips/{login}",
		GetUserIPs,
	},
	Route{
		"Get cluster info",
		"GET",
		"/clusterInfo",
		GetClusterInfo,
	},
	Route{
		"Get database info",
		"GET",
		"/dbInfo/{dbName}",
		GetDbInfo,
	},
	Route{
		"Get queries stats",
		"GET",
		"/queryStats",
		GetQueryStats,
	},
	Route{
		"Test DB speed",
		"GET",
		"/testDbSpeed/{quantity}",
		TestDbSpeed,
	},
	Route{
		"Help",
		"GET",
		"/help",
		GetHelp,
	},
}

// used for /help get method
var textRoutes []string
