package cache

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	CacheTickerTime    = "CACHE_TICKER_TIME"
	CacheByteSliceSize = "CACHE_BYTE_SLICE_SIZE"
	CacheDefaultTTL    = "CACHE_DEFAULT_TTL"
	CacheCleanupType   = "CACHE_CLEANUP_TYPE"

	DefaultByteSliceSize = 1024
	DefaultTickerTime    = 20 * time.Second
	DefaultTTL           = 30 * time.Minute
	DefaultCleanupType   = "active"
)

type CleanupType string

const (
	Active CleanupType = "active"
	Lazy   CleanupType = "lazy"
	None   CleanupType = "none"
)

type Config struct {
	TickerTime    time.Duration `json:"ticker_time"`
	ByteSliceSize int           `json:"byte_slice_size"`
	DefaultTTL    time.Duration `json:"default_ttl"`
	CleanupType   CleanupType   `json:"cleanup_type"`
}

func DefaultConfig() *Config {
	return &Config{
		TickerTime:    DefaultTickerTime,
		ByteSliceSize: DefaultByteSliceSize,
		DefaultTTL:    DefaultTTL,
		CleanupType:   Active,
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
	if envVar := os.Getenv(CacheTickerTime); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.TickerTime = duration
		}
	}
	if envVar := os.Getenv(CacheByteSliceSize); envVar != "" {
		if size, err := strconv.Atoi(envVar); err == nil {
			config.ByteSliceSize = size
		}
	}
	if envVar := os.Getenv(CacheDefaultTTL); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.DefaultTTL = duration
		}
	}
	if envVar := os.Getenv(CacheCleanupType); envVar != "" {
		config.CleanupType = CleanupType(envVar)
	}

	log.Printf("Loaded Config: %+v", config)

	return config, nil
}
