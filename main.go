package main

import (
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/offer"
)

func main() {
	log.Printf("%s", "Server starting on port 8080...")

	// start module
	initModule()

	router := offer.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initModule() {

	// inits app module setting up worker pool and other global scoped objects
	offer.GetInstance()
}
