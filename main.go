package main

import (
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer"
)

// runs app in PROD mode
func main() {
	run(model.Prod)
}

// run starts the app
// mode - PROD or TEST modes will use different config values depending on mode.
func run(mode string) {
	const defaultPort = ":8080"
	log.Printf("Server starting - port %s - mode %s", defaultPort, mode)

	// build module and setup server to listen at default port
	startServer(defaultPort, mode)
}

// starts a new server instance using mode config and port
func startServer(port, mode string) {
	router := offer.NewRouter()
	// inits app module setting up worker pool and other global scoped objects
	offer.BuildInstance(router, mode)

	log.Fatal(http.ListenAndServe(port, router))
}
