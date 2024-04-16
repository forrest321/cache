> **Notice**: This package is currently under active development. It is not a finalized release, and as such, it may undergo 
> significant changes, including breaking changes, as well as changes to its name and location. Please use it with this understanding.

# Cache

The Cache package provides a high-performance, flexible caching solution designed to be adaptable for various applications and environments. It features configurable expiration strategies, making it suitable for diverse deployment scenarios.

## Features

- **Configurable Expiration**: Choose between active, lazy, or no cleanup strategies.
- **Adaptable**: Ideal for standalone applications, Docker, or Kubernetes environments.
- **Dynamic Configuration**: Manage settings through a JSON configuration file and environment variables.

## Installation

To use this cache package in your Go projects, run the following command:

```bash
go get -u github.com/forrest321/cache
```

This command will retrieve the package and ensure you have the latest version. Make sure to import it in your Go files with:
```go
import "github.com/forrest321/cache"
```

## Usage
Here’s how to import and initialize the cache in your Go application:
```go
package main

import (
    "fmt"
    "github.com/forrest321/cache"
    "log"
    "os"
)

func main() {
    logger := log.New(os.Stdout, "cache: ", log.LstdFlags)
    config, err := cache.LoadConfig("path/to/config.json") // Adjust path as needed
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }
    cacheInstance := cache.NewCache(config, logger)

    // Example of setting a value in the cache
    cacheInstance.Set([]byte("key"), []byte("value"), config.DefaultTTL)

    // Example of getting a value from the cache
    value, _ := cacheInstance.Get([]byte("key"))
    fmt.Println("Cached Value:", string(value))
}
```

## Configuration
Modify config.json to adjust cache settings. Example configuration:
```json
{
    "ticker_time": "20s",
    "byte_slice_size": 1024,
    "default_ttl": "30m",
    "cleanup_type": "active"
}
```
Adjust the config.json path according to your project structure.

## Contributing
Contributions are welcome. Please submit pull requests or create issues for any bugs and feature suggestions.