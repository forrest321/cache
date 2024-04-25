package store

import (
	"github.com/forrest321/cache"
	"github.com/forrest321/cache/common/interfaces"
	"log"
	"os"
	"sync"
	"time"
)

type InMemoryStore struct {
	Cache
}

func NewInMemoryStore(config *cache.Config, logger interfaces.Logger) *InMemoryStore {
	if logger == nil {
		logger = log.New(os.Stdout, "cache: ", log.LstdFlags)
	}
	return &InMemoryStore{
		Cache: Cache{
			config:  config,
			entries: make(map[string]Entry),
			lock:    sync.RWMutex{},
			logger:  logger,
		},
	}
}

// Get Retrieves a value for a given key; returns nil if no such key exists
func (s *InMemoryStore) Get(key string) []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if entry, found := s.entries[string(key)]; found {
		if entry.expiresAt.After(time.Now()) {
			return entry.value
		}
		if s.config.CleanupType == "lazy" {
			s.lock.RUnlock()
			s.lock.Lock()
			delete(s.entries, string(key))
			s.lock.Unlock()
			s.logger.Printf("Lazy cleanup: expired item removed: %s", key)
			s.lock.RLock()
		}
		return nil
	}
	return nil
}

// Set Inserts a value into the cache using the provided key
func (s *InMemoryStore) Set(key string, value []byte, expiresAt time.Time) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.entries[key] = Entry{
		value:     value,
		expiresAt: expiresAt,
	}
}

// Delete Removes a key-value pair from the cache; does nothing if key does not exist
func (s *InMemoryStore) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.entries, key)
}

// StartCleanup initiates a cleanup routine that periodically removes expired entries from the cache.
// The cleanupInterval parameter specifies the time interval between each cleanup operation.
// Panic is thrown as a placeholder until the implementation is completed.
func (s *InMemoryStore) StartCleanup(cleanupInterval time.Duration) {
	if s.config.CleanupType == "active" {
		ticker := time.NewTicker(s.config.TickerTime)
		for range ticker.C {
			s.lock.Lock()
			for key, entry := range s.entries {
				if entry.expiresAt.Before(time.Now()) {
					delete(s.entries, key)
					s.logger.Printf("Active cleanup: expired item removed: %s", key)
				}
			}
			s.lock.Unlock()
		}
	}
}
