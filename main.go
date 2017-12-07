package main

import (
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/offer"
)

func main() {
	run()
}

func run() {
	const defaultPort = ":8080"
	log.Printf("Server starting - port %s ...", defaultPort)

	// build module and setup server to listen at default port
	startServer(defaultPort)
}

func startServer(port string) {
	router := offer.NewRouter()
	// inits app module setting up worker pool and other global scoped objects
	offer.BuildInstance(router)

	log.Fatal(http.ListenAndServe(port, router))
}
