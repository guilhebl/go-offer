package offer

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
		"Search",
		"POST",
		"/offers",
		Search,
	},
	Route{
		"Offers",
		"GET",
		"/offerlist",
		SearchDatastore,
	},
	Route{
		"AddOffer",
		"POST",
		"/offerlist",
		AddOffer,
	},
	Route{
		"Reset",
		"GET",
		"/reset",
		ResetDatastore,
	},
	Route{
		"Show",
		"GET",
		"/offers/{id}",
		Show,
	},
}
