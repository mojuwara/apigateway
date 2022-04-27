package main

import (
	"sync"
	"time"
)

// Lock around the instCache
var instLock sync.RWMutex

/*
Cache for Instances, mapping service manes to their instances

instances = {
	"service_name": {
		Instance: true,...
	},
}
*/
var instanceCache map[string]map[Instance]bool

func initInstanceCache() {
	instanceCache = make(map[string]map[Instance]bool)
	logger.Println("Initialized the cache for instances")
}

func getInstance(service string) (Instance, bool) {
	var instance Instance

	//////////////////////////////////////////////////////////////////
	instLock.RLock()

	instanceGroup, ok := instanceCache[service]
	if ok {
		instance, ok = getRandomInstance(instanceGroup)
	}

	instLock.RUnlock()
	//////////////////////////////////////////////////////////////////

	if ok {
		logger.Printf("Instance '%s' randomly chosen for service: '%s'\n", instance.Addr, service)
		return instance, true
	} else {
		logger.Printf("No available instances for service: '%s'\n", service)
		return instance, false
	}
}

// Find a random instance that is not expired yet
func getRandomInstance(instanceGroup map[Instance]bool) (Instance, bool) {
	for instance := range instanceGroup {
		if instance.Expire.After(time.Now()) {
			return instance, true
		} else {
			delete(instanceGroup, instance)
		}
	}
	return Instance{}, false
}

func registerInstance(instance Instance) bool {
	//////////////////////////////////////////////////////////////////
	instLock.Lock()

	instanceGroup, ok := instanceCache[instance.Service]
	if !ok {
		instanceCache[instance.Service] = make(map[Instance]bool)
	}
	instanceGroup[instance] = true

	instLock.Unlock()
	//////////////////////////////////////////////////////////////////

	logger.Printf("Successfully cached instance '%s' for service '%s':", instance.Addr, instance.Service)
	return true
}

// TODO: If service says its not available for a certain service
func unregisterInstance(instance Instance) bool {
	//////////////////////////////////////////////////////////////////
	instLock.Lock()

	instanceGroup, ok := instanceCache[instance.Service]
	if ok {
		delete(instanceGroup, instance)
	}

	instLock.Unlock()
	//////////////////////////////////////////////////////////////////

	logger.Printf("Successfully cached instance '%s' for service '%s':", instance.Addr, instance.Service)
	return true
}
