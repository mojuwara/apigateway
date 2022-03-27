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

// How frequently, instnaces should ping us to be considered "alive"
const TTL_INSTANCES = time.Minute
