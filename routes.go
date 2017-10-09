package main

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
		Index,
	},
	Route{
		"OfferIndex",
		"GET",
		"/offers",
		OfferIndex,
	},
	Route{
		"OfferCreate",
		"POST",
		"/offers",
		OfferCreate,
	},
	Route{
		"OfferShow",
		"GET",
		"/offers/{id}",
		OfferShow,
	},
}