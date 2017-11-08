package main

import (
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/offer"
)

func main() {
	log.Printf("%s","Server starting on port 8080...")

	router := offer.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}