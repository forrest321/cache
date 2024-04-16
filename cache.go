package cache

import (
	"log"
	"os"
	"sync"
	"time"
)

// Logger interface for logging messages
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// Entry represents each entry in the cache
type Entry struct {
	value     []byte
	expiresAt time.Time
}

// Cache struct for the caching system
type Cache struct {
	entries map[string]Entry
	lock    sync.RWMutex
	config  *Config
	logger  Logger
}

// NewCache initializes a new cache with configurations and a logger
func NewCache(config *Config, logger Logger) *Cache {
	if logger == nil {
		logger = log.New(os.Stdout, "cache: ", log.LstdFlags)
	}

	return &Cache{
		entries: make(map[string]Entry),
		config:  config,
		logger:  logger,
	}
}

// Set adds a new item to the cache or updates an existing one
func (c *Cache) Set(key, value []byte, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.entries[string(key)] = Entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache if it hasn't expired
func (c *Cache) Get(key []byte) (value []byte, ttl time.Duration) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if entry, found := c.entries[string(key)]; found {
		if entry.expiresAt.After(time.Now()) {
			return entry.value, entry.expiresAt.Sub(time.Now())
		}
		if c.config.CleanupType == "lazy" {
			c.lock.RUnlock()
			c.lock.Lock()
			delete(c.entries, string(key))
			c.lock.Unlock()
			c.logger.Printf("Lazy cleanup: expired item removed: %s", key)
			c.lock.RLock()
		}
		return nil, 0
	}
	return nil, 0
}

// cleanupExpiredItems handles the periodic cleanup of expired items for active cleanup
func (c *Cache) cleanupExpiredItems() {
	if c.config.CleanupType == "active" {
		ticker := time.NewTicker(c.config.TickerTime)
		for range ticker.C {
			c.lock.Lock()
			for key, entry := range c.entries {
				if entry.expiresAt.Before(time.Now()) {
					delete(c.entries, key)
					c.logger.Printf("Active cleanup: expired item removed: %s", key)
				}
			}
			c.lock.Unlock()
		}
	}
}
