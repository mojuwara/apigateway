package main

import (
	"io"
	"time"
)

type Instance struct {
	Service string `json:"service"` // The service it is an instance of(Ex: Msgbroker instance)
	Addr    string `json:"addr"`    // host:port or domain.com format

	// Used internally
	Time   time.Time // Time the ping was receuved
	Expire time.Time // Time the service will expire
}

type CachedResp struct {
	Url        string // URL in the request made by clients
	Body       []byte // Body of response
	Status     string
	StatusCode int
	ExpireTime time.Time

	From    string
	Service string
	Time    time.Time
}

// type Service struct {
// 	Name      string     `json:"service"`
// 	Endpoints []Endpoint `json:"endpoints"`
// }

// type Endpoint struct {
// 	Path     string `json:"path"`      // Path clients will use
// 	AuthPath string `json:"auth_path"` // Path used to verify authorization for this call, if any
// }

type RequestInfo struct {
	By       string
	For      string
	Url      string
	Service  string
	Endpoint string
	Method   string
	Proto    string
	Host     string
	Body     io.ReadCloser
}
