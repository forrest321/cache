package interfaces

import "time"

// Store represents a cache store.
type Store interface {
	// Set a new item in the cache that will expire after the default expiration time.
	Set(key string, value interface{})

	// Get an item from the cache. Returns the item or nil, and a found flag.
	Get(key string) (interface{}, bool)

	// Delete an item from the cache. Does nothing if the key is not in the cache.
	Delete(key string)

	// StartCleanup starts the background expiration cleanup goroutine.
	// cleanupInterval determines how frequently expired cache items are removed.
	StartCleanup(cleanupInterval time.Duration)
}
