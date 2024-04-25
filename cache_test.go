package cache

import (
	"bytes"
	"github.com/forrest321/cache/common"
	"github.com/forrest321/cache/common/interfaces"
	"testing"
	"time"
)

// TestNewCache checks if the cache initializes correctly with the test config
func TestNewCache(t *testing.T) {
	cache := setupTestCache()
	if cache == nil {
		t.Fatal("Failed to create cache with test configuration")
	}
}

// TestCacheOperations tests setting and retrieving items from the cache
func TestCacheOperations(t *testing.T) {
	cache := setupTestCache()
	key := "testKey"
	value := []byte("testValue")
	cache.Set(key, value, time.Now().Add(cache.config.DefaultTTL))

	retrievedValue, _ := cache.Get(key)
	if string(retrievedValue) != string(value) {
		t.Errorf("Expected %s, got %s", value, retrievedValue)
	}

	// Wait for item to expire and check if it's gone
	time.Sleep(cache.config.DefaultTTL + 1*time.Second)
	retrievedValue, _ = cache.Get(key)
	if retrievedValue != nil {
		t.Errorf("Expected item to be nil after expiration, got %s", retrievedValue)
	}
}

// TestBackgroundCleanup verifies that expired items are removed
func TestBackgroundCleanup(t *testing.T) {
	cache := setupTestCache()
	key := "expireKey"
	value := []byte("expireValue")
	shortTTL := 1 * time.Second
	cache.Set(key, value, time.Now().Add(shortTTL))
	time.Sleep(shortTTL + 1*time.Second) // Allow time for the item to expire

	if val, _ := cache.Get(key); val != nil {
		t.Errorf("Expected key to be nil, but got %s", val)
	}
}

func TestLazyCleanup(t *testing.T) {
	config := &common.Config{
		TickerTime:    1 * time.Minute, // Longer ticker time, irrelevant for lazy cleanup
		ByteSliceSize: 512,
		DefaultTTL:    1 * time.Second, // Short TTL to quickly expire items
		CleanupType:   "lazy",
	}
	logger := newTestLogger()
	cache, _ := New("", nil)
	cache.config = *config
	cache.logger = logger

	key := "lazyKey"
	value := []byte("value")
	cache.Set(key, value, time.Now().Add(1*time.Second)) // Set item with very short TTL

	time.Sleep(2 * time.Second) // Ensure the item is expired

	// Attempt to access the expired item, triggering lazy cleanup
	if val, found := cache.Get(key); val != nil && found == false {
		t.Errorf("Expected item to be cleaned up lazily, but was still accessible")
	}
}

func TestNoCleanup(t *testing.T) {
	config := &common.Config{
		TickerTime:    1 * time.Minute,
		ByteSliceSize: 512,
		DefaultTTL:    1 * time.Second,
		CleanupType:   "none",
	}
	logger := newTestLogger()
	cache, _ := New("", nil)
	cache.config = *config
	cache.logger = logger

	key := "noCleanupKey"
	value := []byte("value")
	cache.Set(key, value, time.Now().Add(1*time.Second)) // Set item with very short TTL

	time.Sleep(2 * time.Second) // Ensure the item is expired

	// Access the expired item, which should not trigger cleanup
	if val, found := cache.Get(key); val != nil && found == false {
		t.Errorf("Expected no cleanup, but the item appears to have been removed")
	}
}

func TestSetAndReset(t *testing.T) {
	cache := setupTestCache()
	key := "resetKey"
	initialValue := []byte("initialValue")
	updatedValue := []byte("updatedValue")

	cache.Set(key, initialValue, time.Now().Add(cache.config.DefaultTTL))
	firstValue, _ := cache.Get(key)
	if string(firstValue) != string(initialValue) {
		t.Errorf("Expected initial value %s, got %s", initialValue, firstValue)
	}

	// Reset the value
	cache.Set(key, updatedValue, time.Now().Add(cache.config.DefaultTTL))
	secondValue, _ := cache.Get(key)
	if string(secondValue) != string(updatedValue) {
		t.Errorf("Expected updated value %s, got %s", updatedValue, secondValue)
	}
}

func TestLargeDataHandling(t *testing.T) {
	cache := setupTestCache()
	key := "largeKey"
	largeValue := make([]byte, 1024*1024) // 1 MB of data
	for i := range largeValue {           // Fill the byte slice with data
		largeValue[i] = byte(i % 256)
	}

	cache.Set(key, largeValue, time.Now().Add(cache.config.DefaultTTL))
	retrievedValue, _ := cache.Get(key)
	if !bytes.Equal(retrievedValue, largeValue) {
		t.Error("The retrieved value does not match the set large value")
	}
}

// BenchmarkSetAndGet benchmarks the Set and Get operations
func BenchmarkSetAndGet(b *testing.B) {
	cache := setupTestCache()
	key := "benchKey"
	value := []byte("benchValue")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(key, value, time.Now().Add(cache.config.DefaultTTL))
		_, _ = cache.Get(key)
	}
}

func BenchmarkSetAndReset(b *testing.B) {
	cache := setupTestCache()
	key := "benchResetKey"
	value := []byte("benchValue")
	newValue := []byte("newBenchValue")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(key, value, time.Now().Add(cache.config.DefaultTTL))
		cache.Set(key, newValue, time.Now().Add(cache.config.DefaultTTL))
	}
}

func BenchmarkLargeDataHandling(b *testing.B) {
	cache := setupTestCache()
	key := "benchLargeKey"
	largeValue := make([]byte, 1024*1024) // 1 MB
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(key, largeValue, time.Now().Add(cache.config.DefaultTTL))
		_, _ = cache.Get(key)
	}
}

// setupTestCache helps set up a cache with testing configuration
func setupTestCache() *Cache {
	config := &common.Config{
		TickerTime:    1 * time.Second,
		ByteSliceSize: 512,
		DefaultTTL:    2 * time.Second, // Reduce TTL for quicker test turnover
		CleanupType:   "active",
	}
	logger := newTestLogger()
	cache, _ := New("", nil)
	cache.config = *config
	cache.logger = logger
	return cache
}

// TestLogger implements the Logger interface for testing purposes
type TestLogger struct{}

var _ interfaces.Logger = (*TestLogger)(nil)

func newTestLogger() *TestLogger {
	return &TestLogger{}
}

func (tl *TestLogger) Printf(format string, v ...interface{}) {
	// Output can be seen in verbose test mode if needed
	testing.Verbose()
}

func (tl *TestLogger) Println(v ...interface{}) {
	// Output can be seen in verbose test mode if needed
	testing.Verbose()
}

func (tl *TestLogger) Fatal(v ...any) {
	// Output can be seen in verbose test mode if needed
	testing.Verbose()
}
