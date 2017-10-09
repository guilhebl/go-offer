package main

import (
	"log"
	"net/http"
)

func main() {
	log.Printf("%s","Server starting on port 8080...")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}