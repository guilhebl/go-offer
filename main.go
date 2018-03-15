package main

import (
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/offer"
)

// runs app in PROD mode
func main() {
	run("prod")
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

	// inits app module setting up worker pool and other global scoped objects
	offer.BuildInstance(mode)

	log.Fatal(http.ListenAndServe(port, offer.GetInstance().Router))
}
