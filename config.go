package cache

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	EnvTickerTime    = "CACHE_TICKER_TIME"
	EnvByteSliceSize = "CACHE_BYTE_SLICE_SIZE"
	EnvDefaultTTL    = "CACHE_DEFAULT_TTL"
	EnvCleanupType   = "CACHE_CLEANUP_TYPE"
	EnvStoreType     = "CACHE_STORE_TYPE"

	DefaultByteSliceSize = 1024
	DefaultTickerTime    = 20 * time.Second
	DefaultTTL           = 30 * time.Minute
	DefaultCleanupType   = Active
	DefaultStoreType     = InMemoryStoreType
)

type CleanupType string

const (
	Active CleanupType = "active"
	Lazy   CleanupType = "lazy"
	None   CleanupType = "none"
)

type StoreType string

const (
	InMemoryStoreType StoreType = "inmemory"
	// Here is where will add store types later
	// DiskStoreType   StoreType = "disk"
	// RemoteStoreType StoreType = "remote"
)

type Config struct {
	TickerTime    time.Duration `json:"ticker_time"`
	ByteSliceSize int           `json:"byte_slice_size"`
	DefaultTTL    time.Duration `json:"default_ttl"`
	CleanupType   CleanupType   `json:"cleanup_type"`
	StoreType     StoreType     `json:"store_type"`
}

func DefaultConfig() *Config {
	return &Config{
		TickerTime:    DefaultTickerTime,
		ByteSliceSize: DefaultByteSliceSize,
		DefaultTTL:    DefaultTTL,
		CleanupType:   DefaultCleanupType,
		StoreType:     DefaultStoreType,
	}
}

func LoadConfig(filePath string) (*Config, error) {
	config := DefaultConfig()

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening config file: %v", err)
		} else {
			defer func() {
				_ = file.Close()
			}()
			decoder := json.NewDecoder(file)
			if err = decoder.Decode(config); err != nil {
				log.Printf("Error decoding config file: %v", err)
			}
		}
	}

	// Override config with environment variables if they exist
	if envVar := os.Getenv(EnvTickerTime); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.TickerTime = duration
		}
	}
	if envVar := os.Getenv(EnvByteSliceSize); envVar != "" {
		if size, err := strconv.Atoi(envVar); err == nil {
			config.ByteSliceSize = size
		}
	}
	if envVar := os.Getenv(EnvDefaultTTL); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.DefaultTTL = duration
		}
	}
	if envVar := os.Getenv(EnvCleanupType); envVar != "" {
		envCleanType := CleanupType(envVar)

		switch envCleanType {
		case Active, Lazy, None:
			config.CleanupType = envCleanType
		default:
			log.Printf("Invalid cleanup type provided: %v", envVar)
			config.CleanupType = Lazy
		}
	}
	if envVar := os.Getenv(EnvStoreType); envVar != "" {
		envStoreType := StoreType(envVar)

		//When adding more store types, create more cases here
		switch envStoreType {
		case InMemoryStoreType:
			config.StoreType = envStoreType
		default:
			log.Printf("Invalid store type provided: %v", envVar)
			config.StoreType = InMemoryStoreType
		}
	}

	log.Printf("Loaded Config: %+v", config)

	return config, nil
}
