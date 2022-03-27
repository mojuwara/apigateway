package main

import "time"

type Instance struct {
	Service string `json:"service"` // The service it is an instance of(Ex: Msgbroker instance)
	Addr    string `json:"addr"`    // host:port or domain.com format

	// Used internally
	Time   time.Time // Time the ping was receuved
	Expire time.Time // Time the service will expire
}

// Maps the Addr of an instance to its object
type AddrMap map[string]Instance

// Maps names of services to their instances
type ServicesMap map[string]AddrMap
