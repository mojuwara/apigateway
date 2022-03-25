package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize connection with Redis server
	initRedis()

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/service", serviceHandler)

	log.Fatal(http.ListenAndServe(":5678", nil))
}
