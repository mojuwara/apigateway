package main

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Logger
var logger = log.Default()

// Redis client shared by goroutines
var rdb *redis.Client

// How frequently, in seconds, instnaces should ping us to be considered "alive"
const TTL_INSTANCES = 30 * time.Second

// Prefix used by all instances
const PREFIX_INSTANCES = "service_name:"
