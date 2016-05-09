package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"log"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Route descriptors
type Routes []Route
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"BusinessesAll",
		"GET",
		"/businesses",
		BusinessesAll,
	},
	Route{
		"BusinessById",
		"GET",
		"/business/{id}",
		BusinessById,
	},
}

// Aggregate routes to pass to server
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(Logger(route.HandlerFunc, route.Name))
	}
	return router
}

// Basic logger routine
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

