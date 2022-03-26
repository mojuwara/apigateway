# API Gateway

### Possible Implementations:
- `(Winner) Each instance is a key(<service>:<instance>)`
	- When a service is selected, it scan the key space
	- Pros:
		- Can use Redis built-in expire to auto-remove
	- Cons:
		- Must scan the set of keys that have the <service> prefix
		- Not random
- `Each service is a key(service:<service_name>)` and the instances of that service are members of the set
	- When that service is requested, it selects a random instance
	- Pros:
		- Randomness for load balancing
		- Easy to update TTL for instances when they ping us
	- Cons:
		- Must check if the instance is expired before using it, possibly asking for multiple instances fo that service
		- Uses too much space. Each ping from an instance is stored in the set
- `Implement custom cache`, mapping services to instances
	- Pros:
		- Full control
	- Cons:
		- Locking/synchronization to update the mapping

### Learn
- Pipelines
- Rings