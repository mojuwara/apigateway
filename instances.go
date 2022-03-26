package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Key is in the format "service:<service_name>"
// All instances of a service will have the same "<service>" prefix for scanning
func CreateServiceKey(service string) string {
	return fmt.Sprintf("service:%s", service)
}

// Create/update set of instances for this service
func UpdateInstance(inst Instance) bool {
	// Encode the instance for storage
	bytes, err := json.Marshal(inst)
	if err != nil {
		logger.Printf("Error encoding instance '%s' for service '%s'. Err: %s", inst.Addr, inst.Service, err)
		return false
	} else {
		logger.Println("Marshalled", bytes)
	}

	ctx := context.Background()
	key := CreateServiceKey(inst.Service)

	err = rdb.SAdd(ctx, key, bytes).Err()
	if err != nil {
		logger.Printf("Error adding instance '%s' for service '%s' into Redis. Err: %s", inst.Addr, inst.Service, err)
		return false
	} else {
		logger.Printf("Successfully added instance '%s' for service '%s' into Redis", inst.Addr, inst.Service)
		return true
	}
}

// Returns the addr of a random instance of this service
// TODO: Speed up by adding a set in Redis that stores these keys. Check from the set first and update occassionally, but may be overly-complicated
func GetInstance(service string) string {
	logger.Println("Getting a random instance of service:", service)

	ctx := context.Background()
	key := CreateServiceKey(service)

	for rdb.SCard(ctx, key).Val() > 0 {
		marshalledInst, _ := rdb.SRandMember(ctx, key).Result()
		instance := unmarshalInstance(marshalledInst)

		// If it is currently past the expire time, remove from set
		if time.Now().After(instance.Expire) {
			logger.Printf("Removing '%s' instance from service '%s'. It expired at %s\n", instance.Addr, instance.Service, instance.Expire)
			rdb.SRem(ctx, key, marshalledInst)
		} else {
			logger.Printf("Chose '%s' as random instance for service '%s'\n", instance.Addr, instance.Service)
			return instance.Addr
		}
	}

	logger.Printf("No instances available for requested service '%s'\n", service)
	return ""
}

func unmarshalInstance(instStr string) Instance {
	var inst Instance
	json.Unmarshal([]byte(instStr), &inst)
	return inst
}
