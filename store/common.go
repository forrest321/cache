package store

import (
	"github.com/forrest321/cache/common/interfaces"
	"sync"
	"time"
)

// Entry represents each entry in the cache
type Entry struct {
	value     []byte
	expiresAt time.Time
}

// Cache struct for the caching system
type Cache struct {
	entries     map[string]Entry
	lock        sync.RWMutex
	logger      interfaces.Logger
	cleanupType string
}
