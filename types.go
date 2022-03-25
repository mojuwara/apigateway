package main

type Instance struct {
	Service string `json:"service"` // The service it is an instance of(Ex: Msgbroker instance)
	Addr    string `json:"addr"`    // host:port or domain.com format
}
