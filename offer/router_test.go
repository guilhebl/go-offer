package offer

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

	assert.True(t, strings.Contains(pathRegex, regex))
	assert.Equal(t, methods[0], methodName)
}
