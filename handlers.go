package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Handles POST requests from servers making themselves known
func pingHandler(w http.ResponseWriter, r *http.Request) {
	// Read body of message
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Println(err)
		return
	}

	// Unmarshal body of message into Instance object
	var ping Instance
	err = json.Unmarshal(body, &ping)
	if err != nil {
		logger.Println("Error", err, "while unmarshalling ping", string(body))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not unmarshal ping:" + err.Error()))
		return
	}

	// Update cache with this instance
	UpdateInstance(&ping)
	w.Write([]byte("Message received"))
}

// Should be used internally, clients shouldn't have to be aware of the services
// Returns an addr(string) of an instance of the given service
// URL should be "service/?name=<service_name>"
func serviceHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	service := query.Get("name")
	if service == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Requests for service discovery must have the service name in the 'name' query param"))
		return
	}

	// Get a random instance of this service
	host := GetInstance(service)
	if host == "" {
		w.Write([]byte("There are no instances currently running for this service."))
	}

	w.Write([]byte(host))
}
