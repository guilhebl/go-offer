package offer

import (
	"github.com/gorilla/mux"
	"strings"
	"testing"
)

// tests if router is built correctly and routes to the right paths
func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Error("Error while creating Router")
	}
	testRoute(t, router, "Index", "/", "GET")
	testRoute(t, router, "Search", "/offers", "POST")
	testRoute(t, router, "Show", "/offers/", "GET")
}

func testRoute(t *testing.T, router *mux.Router, name, regex, methodName string) {
	routeIndex := router.Get(name)
	pathRegex, _ := routeIndex.GetPathRegexp()
	methods, _ := routeIndex.GetMethods()

	if !strings.Contains(pathRegex, regex) {
		t.Error("wrong route path regex")
	}

	if methods[0] != methodName {
		t.Error("wrong route METHOD")
	}
}
