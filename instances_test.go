package main

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestCreateInstanceKey(t *testing.T) {
	var (
		randService = "service_a"
		randHost    = "host_a"
		inst        = Instance{Service: randService, Addr: randHost}
	)

	expected := randService + ":" + randHost // "service_a:host_a"
	actual := CreateInstanceKey(&inst)
	assert.Equal(t, expected, actual, "Key created for instance isn't in expected format")
}

// TODO: Make this test better
func TestGetInstance(t *testing.T) {
	var (
		randService = "service_a"
		randHost    = "host_a"
		inst        = Instance{Service: randService, Addr: randHost}
		key         = CreateInstanceKey(&inst)
	)

	// Set the global rdb for testing
	var mock redismock.ClientMock
	rdb, mock = redismock.NewClientMock()

	// Test on empty DB
	mock.ExpectScan(0, randService+":*", 0)
	expectedHost := ""
	actualHost := GetInstance(randService)
	assert.Equal(t, expectedHost, actualHost, "Empty DB should not have any hosts")

	// Expect to get a random host after insterting this instance
	mock.ExpectSet(key, randHost, 0)
	statusErr := rdb.Set(context.TODO(), key, randHost, 0)
	if statusErr != nil {
		return
	}
	actualHost = GetInstance(randService)
	assert.Equal(t, randHost, actualHost, "Instance returned from GetInstance not expected")
}

func TestUpdateInstance(t *testing.T) {
	var (
		randService = "service_a"
		randHost    = "host_a"
		inst        = Instance{Service: randService, Addr: randHost}
		key         = CreateInstanceKey(&inst)
	)

	// Set the global rdb for testing
	var mock redismock.ClientMock
	rdb, mock = redismock.NewClientMock()

	// Test inserting an instance
	mock.ExpectSet(key, randHost, TTL_INSTANCES).SetVal(randHost)
	retVal := UpdateInstance(&inst)
	assert.True(t, retVal, "Host not found after inserting in Redis")

	// Test re-inserting an instance, should reset the TTL
	mock.ExpectSet(key, randHost, TTL_INSTANCES).SetVal(randHost)
	retVal = UpdateInstance(&inst)
	assert.True(t, retVal, "Host not found after re-inserting in Redis")
}
