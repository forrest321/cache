package cache

import (
	"log"
	"os"
	"sync"
	"time"
)

// Logger interface abstracts logging functions to allow for different logging implementations
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// Node represents each entry in the cache
type Node struct {
	key       []byte
	value     []byte
	expiresAt time.Time
	next      *Node
}

// Cache is the main structure representing the cache itself
type Cache struct {
	head   *Node
	lock   sync.Mutex
	pool   sync.Pool
	config *Config
	logger Logger
}

// NewCache creates a new cache instance with the provided configuration and logger
func NewCache(config *Config, logger Logger) *Cache {
	if logger == nil {
		// Default to standard logger if none provided
		logger = log.New(os.Stdout, "cache: ", log.LstdFlags)
	}

	cache := &Cache{
		config: config,
		logger: logger,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, config.ByteSliceSize)
			},
		},
	}

	// Start the cleanup process if cleanup is set to active
	if config.CleanupType == "active" {
		go cache.cleanupExpiredItems()
	}

	return cache
}

// Set adds a new item to the cache or updates an existing one
func (c *Cache) Set(key, value []byte, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	newNode := &Node{key: key, value: value, expiresAt: time.Now().Add(ttl)}
	newNode.next = c.head
	c.head = newNode
}

func (c *Cache) Get(key []byte) (value []byte, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var prev *Node = nil
	current := c.head
	for current != nil {
		if current.expiresAt.Before(time.Now()) {
			if c.config.CleanupType == "lazy" {
				// Remove expired node
				if prev == nil {
					c.head = current.next
				} else {
					prev.next = current.next
				}
				c.logger.Printf("Lazy cleanup: expired item removed: %s", current.key)
				current = current.next
				continue
			} else if c.config.CleanupType == "none" {
				// Do not remove expired items
				if string(current.key) == string(key) {
					return nil, 0 // Expired item is still returned as nil
				}
			}
		}
		if string(current.key) == string(key) {
			return current.value, current.expiresAt.Sub(time.Now())
		}
		prev = current
		current = current.next
	}
	return nil, 0
}

func (c *Cache) cleanupExpiredItems() {
	if c.config.CleanupType != "active" {
		return // Exit if cleanup is not set to active
	}

	ticker := time.NewTicker(c.config.TickerTime)
	for range ticker.C {
		c.lock.Lock()
		var prev *Node = nil
		current := c.head
		for current != nil {
			next := current.next
			if current.expiresAt.Before(time.Now()) {
				if prev == nil {
					c.head = next
				} else {
					prev.next = next
				}
				c.logger.Printf("Active cleanup: expired item removed: %s", current.key)
			} else {
				prev = current
			}
			current = next
		}
		c.lock.Unlock()
	}
}
