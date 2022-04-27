package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the Cache for instances
	initInstanceCache()

	// Initialize the cache for responses
	initResponseCache()

	// Create a handler for the configured endpoints
	// initConfigEndpoints()

	http.HandleFunc("/register", register)
	http.HandleFunc("/unregister", unregister)
	http.HandleFunc("/discover", discoveryHandler)

	// Handles all other requests made by clients and other services
	http.HandleFunc("/", genericHandler)

	log.Fatal(http.ListenAndServe(":5678", nil))
}
