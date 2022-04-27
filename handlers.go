package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Handles POST requests from servers making themselves known
func register(w http.ResponseWriter, r *http.Request) {
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

	ping.Time = time.Now()
	ping.Expire = ping.Time.Add(TTL_INSTANCES)

	// Keep track of this instance
	if registerInstance(ping) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Handles POST requests from servers making themselves known
func unregister(w http.ResponseWriter, r *http.Request) {
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

	ping.Time = time.Now()
	ping.Expire = ping.Time.Add(TTL_INSTANCES)

	// Keep track of this instance
	if unregisterInstance(ping) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Should be used internally, clients shouldn't have to be aware of the services
// Returns an addr(string) of an instance of the given service
// URL should be "service/?name=<service_name>"
func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	service := query.Get("name")
	if service == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Requests for service discovery must have the service name in the 'name' query param"))
		return
	}

	// Get a random instance of this service
	instance, ok := getInstance(service)
	if ok {
		w.Write([]byte(instance.Addr))
	} else {
		w.Write([]byte("There are no instances currently running for this service."))
	}
}

// All client requests configured in endpoints.json will be handled here
// Expected the request URL will be <domain>.com/<service_name>/<endpoint>
// Ex: www.google.com/user/login => Call the "login" endpoint for the "user" service
func genericHandler(w http.ResponseWriter, r *http.Request) {
	if val, ok := getCachedResponse(r.URL.Path); ok {
		copyCachedResponse(val, w)
		return
	}

	// HTTP Client to send forwarded request
	client := &http.Client{Timeout: REQUEST_TIMEOUT}

	// Make request on behalf of user
	reqInfo := getRequestInfo(r)
	req, err := http.NewRequest(reqInfo.Method, reqInfo.Endpoint, reqInfo.Body)
	if err != nil {
		logger.Printf("Error while creating forwarded request. Err: '%s'\n", err)
		return
	}
	req.Header.Add(X_FORWARDED_FOR, reqInfo.By)
	req.Header.Add(X_FORWARDED_HOST, reqInfo.Host)
	req.Header.Add(X_FORWARDED_PROTO, reqInfo.Proto)

	logger.Printf("Forwarding '%s' '%s' request from '%s' to instance '%s' of service '%s' with endpoint '%s'\n",
		reqInfo.Proto,
		reqInfo.Method,
		reqInfo.By,
		reqInfo.For,
		reqInfo.Service,
		reqInfo.Endpoint,
	)

	resp, err := client.Do(req)
	if err != nil {
		logger.Println("Error while sending forwarded request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cacheResponse(reqInfo.Url, resp)
	copyResponse(resp, w)
}

func getRequestInfo(r *http.Request) RequestInfo {
	service, endpoint := getServiceAndEndpoint(r.URL)
	instance, _ := getInstance(service)

	return RequestInfo{
		Body:     r.Body,
		Proto:    r.Proto,
		Service:  service,
		Endpoint: endpoint,
		Method:   r.Method,
		Url:      r.URL.Path, // The entire URL path, including the service
		By:       r.RemoteAddr,
		For:      instance.Addr,
		Host:     r.Header.Get(HOST), // Location of this API Gateway
	}
}

// Return the service & endpoint of the URL, possibly empty
func getServiceAndEndpoint(url *url.URL) (string, string) {
	var service string
	var endpoint string

	ndx := strings.Index(url.Path, "/")
	if ndx == -1 {
		logger.Println("Could not find service in request:", url.Path)
		return service, endpoint
	}

	service = url.Path[:ndx]
	endpoint = url.Path[ndx:]
	return service, endpoint
}
