package main

import (
	"log"
	"time"
)

// Logger
var logger = log.Default()

// How frequently, instnaces should ping us to be considered "alive"
const TTL_INSTANCES = time.Minute

// Default request timeout is 10 seconds, unless configured otherwise
const REQUEST_TIMEOUT = time.Second * 10

// HTTP HEADERS
const (
	HOST      = "Host"
	FORWARDED = "Forwarded"

	// Regarding proxy
	X_FORWARDED_FOR   = "X-Forwarded-For"
	X_FORWARDED_HOST  = "X-Forwarded-Host"
	X_FORWARDED_PROTO = "X-Forwarded-Proto"

	// Regarding Caching
	CACHE_CONTROL = "Cache-Control"
	MAX_AGE       = "max-age"
)
