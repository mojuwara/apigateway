package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Key is in the format "<service>:<addr>"
// All instances of a service will have the same "<service>" prefix for scanning
func CreateInstanceKey(inst *Instance) string {
	return fmt.Sprintf("%s:%s", inst.Service, inst.Addr)
}

// Waits for data to be pushed in InstanceChannel adn then updates the cache
// TODO: How to handle key not being inserted
func UpdateInstance(inst *Instance) bool {
	ctx := context.Background()

	key := CreateInstanceKey(inst)
	exists := rdb.Exists(ctx, key).Val()

	// Insert key if new, otherwise reset its TTL
	if exists == 1 {
		err := rdb.Expire(ctx, key, TTL_INSTANCES).Err()
		if err != nil {
			logger.Println("Error extending the lifetime of key:", key, ". Err:", err)
			return false
		} else {
			logger.Println("Successfully extended the lifetime of key:", key)
			return true
		}
	} else {
		err := rdb.Set(ctx, key, inst.Addr, TTL_INSTANCES).Err()
		if err != nil {
			logger.Println("Error inserting key", key, "into Redis:", err)
			return false
		} else {
			logger.Println("Successfully inserted key", key, "into Redis")
			return true
		}
	}
}

// Returns the addr of a random instance of this service
// TODO: Speed up by adding a set in Redis that stores these keys. Check from the set first and update occassionally, but may be overly-complicated
func GetInstance(service string) string {
	logger.Println("Getting a random instance of service:", service)

	// First check if there's a cached list of instances for this service
	instance := getCachedInstances(service)
	if instance != "" {
		return removePrefix(service, instance)
	}

	// Scan Redis for instances if none cached, and then cache them
	allInstances := scanInstances(service)
	if len(allInstances) == 0 {
		logger.Println("Unable to find instance of service:", service)
		return ""
	}

	// Store instances in a Redis set for fast retrieval later
	ctx := context.Background()
	err := rdb.SAdd(ctx, service, allInstances).Err()
	if err != nil {
		logger.Println("Error caching instances of service:", service, err)
		return getRandomInstance(service, allInstances)
	}

	// Have this cache expire in 1min so it doesn't get too stale
	err = rdb.Expire(ctx, service, time.Minute).Err()
	if err != nil {
		logger.Println("Error while trying to set expire time for key:", service, err)
		rdb.Del(ctx, service)
	}

	// Finally return a random instance
	instance = getRandomInstance(service, allInstances)
	logger.Println("Found random instance for service:", service, "instance:", instance)
	return instance
}

// Scan Redis keys and find all instances for given service
func scanInstances(service string) []string {
	ctx := context.Background()
	prefix := service + ":*" // "service_name:*"
	allKeys := []string{}

	for {
		var cursor uint64
		keys, cursor, err := rdb.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			logger.Println("Error while scanning Redis for prefix:", prefix, err)
			break
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}
	}
	return allKeys
}

// Find instances of given service from previous scans, if any are active
func getCachedInstances(service string) string {
	// Return if it isn't cached
	ctx := context.Background()
	if rdb.Exists(ctx, service).Val() == 0 {
		return ""
	}

	// Return first random instance of this service in our cache
	// Remove it from the cache if its no longer in Redis(hasn't pings)
	// Or return empty string when the length of the cache is 0
	for {
		instance, _ := rdb.SRandMember(ctx, service).Result()
		if rdb.Exists(ctx, instance).Val() == 1 {
			return instance
		} else {
			rdb.SRem(ctx, service, instance)
		}

		if rdb.SCard(ctx, instance).Val() == 0 {
			return ""
		}
	}
}

// Return a random instance
func getRandomInstance(service string, allInstances []string) string {
	instance := allInstances[rand.Intn(len(allInstances))]
	return removePrefix(service, instance)
}

// Remove the "<service_name>:" prefix
func removePrefix(service string, instance string) string {
	return instance[len(service)+1:]
}
