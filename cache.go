package cache

import (
	"github.com/forrest321/cache/common"
	"github.com/forrest321/cache/common/interfaces"
	"github.com/forrest321/cache/store"
	"log"
	"os"
	"time"
)

type Cache struct {
	store  interfaces.Store
	logger interfaces.Logger
	config common.Config
}

func New(configPath string, logger *log.Logger) (*Cache, error) {
	// Load Config
	config, err := common.LoadConfig(configPath)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return nil, err
	}

	if logger == nil {
		logger = log.New(os.Stdout, "cache: ", log.LstdFlags)
	}

	// Initialize InMemoryStore
	memStore := store.NewInMemoryStore(config, logger)

	return &Cache{
		store:  memStore,
		logger: logger,
		config: *config,
	}, nil
}

// Get passes through to the InMemoryStore Get method.
func (c *Cache) Get(key string) ([]byte, bool) {
	return c.store.Get(key)
}

// Set passes through to the InMemoryStore Set method.
func (c *Cache) Set(key string, value []byte, expiresAt time.Time) {
	c.store.Set(key, value, expiresAt)
}

// Delete passes through to the InMemoryStore Delete method.
func (c *Cache) Delete(key string) {
	c.store.Delete(key)
}

// StartCleanup passes through to the InMemoryStore StartCleanup method.
func (c *Cache) StartCleanup(cleanupInterval time.Duration) {
	c.store.StartCleanup(cleanupInterval)
}
