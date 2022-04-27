package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Lock around the respCache
var respLock sync.RWMutex

// Cache for Responses, mapping request strings to CachedResponse type
var respCache map[string]CachedResp

func initResponseCache() {
	respCache = make(map[string]CachedResp)
	logger.Println("Initialized the cache for responses")
}

func getCachedResponse(request string) (CachedResp, bool) {
	///////////////////////////////////////////////////////////////////
	respLock.RLock()
	defer respLock.RUnlock()

	resp, ok := respCache[request]
	if !ok {
		logger.Printf("Cache does not contain: '%s'\n", request)
		return CachedResp{}, false
	}

	if resp.ExpireTime.After(time.Now()) {
		logger.Printf("Cached response for %s is expired\n", request)
		return resp, false
	}

	logger.Printf("Found cached response for: '%s'\n", request)
	return resp, true
	///////////////////////////////////////////////////////////////////
}

// For responses to be cached, services must set the "max-age" directive
// in the "Cache-Control" header of the response. Ex: max-age=100000
// The value is milliseconds
func cacheResponse(request string, response *http.Response) bool {
	cacheControl := response.Header.Get(CACHE_CONTROL)
	if cacheControl == "" {
		logger.Printf("Cache-Control header is not set for request to: '%s'\n", request)
		return true
	}

	expTime, ok := getExpTime(cacheControl)
	if !ok {
		logger.Printf("Unable to determine max-age of response to request: '%s'", request)
		return false
	}

	logger.Printf("Attempting to cache '%s'\n", request)
	// Read the body for storage
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Printf("Failed to cache response. Error while reading response body made by request '%s'. Err: '%s'", request, err)
		return false
	}

	// Write the response back to the response
	response.Body = ioutil.NopCloser(strings.NewReader(string(body)))

	// Response that is cached
	resp := CachedResp{
		Url:        request,
		Body:       body,
		Status:     response.Status,
		StatusCode: response.StatusCode,
		ExpireTime: expTime,
		// From: response.,
	}

	///////////////////////////////////////////////////////////////////
	respLock.Lock()
	respCache[request] = resp
	respLock.Unlock()
	///////////////////////////////////////////////////////////////////

	logger.Println("Successfully cached response for:", response)
	return true
}

// Given a string with fmt max-age=N, will return the time it will be, N seconds from now
func getExpTime(str string) (time.Time, bool) {
	if !strings.HasPrefix(str, MAX_AGE) {
		return time.Time{}, false
	}

	maxAgeStr := str[len(MAX_AGE)+1:]
	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		logger.Printf("Error while parsing given max-age: str(%s)\n", maxAgeStr)
		return time.Time{}, false
	}

	return time.Now().Add(time.Second * time.Duration(maxAge)), true
}

// Return the key and value of a string that is in the form <key>=<val>
// func splitKeyValue(str string, sep string) (string, string) {
// 	ndx := strings.Index(str, sep)
// 	if ndx == -1 {
// 		return "", ""
// 	}

// 	return str[:ndx], str[ndx+1:]
// }

func copyResponse(resp *http.Response, w http.ResponseWriter) {

}

func copyCachedResponse(resp CachedResp, w http.ResponseWriter) {

}
