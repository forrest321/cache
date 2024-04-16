package cache

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TickerTime    time.Duration `json:"ticker_time"`
	ByteSliceSize int           `json:"byte_slice_size"`
	DefaultTTL    time.Duration `json:"default_ttl"`
	CleanupType   string        `json:"cleanup_type"` // "active", "lazy", or "none"
}

func DefaultConfig() *Config {
	return &Config{
		TickerTime:    20 * time.Second,
		ByteSliceSize: 1024,
		DefaultTTL:    30 * time.Minute,
		CleanupType:   "active",
	}
}

func LoadConfig(filePath string) (*Config, error) {
	config := DefaultConfig()

	if filePath != "" {
		file, err := os.Open(filePath)
		if err == nil {
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			decoder := json.NewDecoder(file)
			_ = decoder.Decode(config) // No need to handle error as defaults are already set
		}
	}

	// Override config with environment variables if they exist
	if envVar := os.Getenv("CACHE_TICKER_TIME"); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.TickerTime = duration
		}
	}
	if envVar := os.Getenv("CACHE_BYTE_SLICE_SIZE"); envVar != "" {
		if size, err := strconv.Atoi(envVar); err == nil {
			config.ByteSliceSize = size
		}
	}
	if envVar := os.Getenv("CACHE_DEFAULT_TTL"); envVar != "" {
		if duration, err := time.ParseDuration(envVar); err == nil {
			config.DefaultTTL = duration
		}
	}
	if envVar := os.Getenv("CACHE_CLEANUP_TYPE"); envVar != "" {
		config.CleanupType = envVar
	}

	return config, nil
}
