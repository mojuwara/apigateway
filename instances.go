package main

import (
	"context"
	"fmt"
	"time"
)

func CreateServiceKey(service string) string {
	return fmt.Sprintf("service:%s", service)
}

// Each service will have its own map/hash "service:<service_name>"
// The fields in the map will be <host_name> and the values will be <expire_time>
func UpdateInstance(inst *Instance) bool {
	ctx := context.Background()

	key := CreateServiceKey(inst.Service)
	field := inst.Addr
	value, _ := inst.Expire.MarshalText()

	_, err := rdb.HSet(ctx, key, field, value).Result()
	if err != nil {
		logger.Printf("Error while storing instance: '%s', expire: '%s' in key: '%s'. Err: %s\n", field, inst.Expire, key, err)
		return false
	} else {
		logger.Printf("Successfully stored instance: '%s', expire: '%s' in key: '%s'\n", field, inst.Expire, key)
		return true
	}
}

// Given a service name, will return a random instance/host if any are known
func GetInstance(service string) string {
	// Number of random hosts to try to fetch from Redis
	const NumHostsToFetch = 5

	ctx := context.Background()
	key := CreateServiceKey(service)

	// While there are still fields in the hash
	for rdb.HLen(ctx, key).Val() > 0 {
		fields, _ := rdb.HRandField(ctx, key, NumHostsToFetch, true).Result()

		for i := 0; i < len(fields); {
			field := fields[i]   // Current element is the field
			value := fields[i+1] // Next element is its value

			// Return current host if hasn't expired
			var expireTime time.Time
			expireTime.UnmarshalText([]byte(value))
			if expireTime.After(time.Now()) {
				logger.Printf("Returning '%s' as random instance for service: '%s'", field, service)
				return field
			}

			// Remove the current host from hash since it has expired and try next pair
			rdb.HDel(ctx, key, field)
			i += 2
		}
	}

	logger.Println("No instances for requested service:", service)
	return ""
}
