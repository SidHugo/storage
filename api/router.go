package api

import (
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		textRoutes = append(textRoutes,
			fmt.Sprintf("%s. Method: %s. URL: %s", route.Name, route.Method, route.Pattern))
	}

	return router
}
