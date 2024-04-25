>Notice: This package is currently under active development. It is not a finalized release, and as such, it may undergo
significant changes, including breaking changes, as well as changes to its name and location. Please use it with this understanding.

# Cache

The Cache package provides a high-performance, flexible caching solution designed to be adaptable for various applications and environments. It features configurable expiration strategies, making it suitable for diverse deployment scenarios.

## Features

Configurable Expiration: Choose between active, lazy, or no cleanup strategies.

Adaptable: Ideal for standalone applications, Docker, or Kubernetes environments.

Dynamic Configuration: Manage settings through a JSON configuration file and environment variables.

## Installation

To use this cache package in your Go projects, run the following command:

```shell
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

You can modify settings through a JSON configuration file as well as environment variables, giving you flexibility according to your use case. Here's what each configuration option means:

`ticker_time` : This is the interval time for checking and clearing expired keys. The default value is "20s" (20 seconds). Accepts a string of number followed by "s" for seconds, "m" for minutes, or "h" for hours.

`byte_slice_size` : The size of byte slices that the cache will batch together for clearing up. The default value is 1024.

`default_ttl` : This is the default time to live for each cache entry if no TTL is provided during the Set operation. The default value is "30m" (30 minutes). Accepts a string of number followed by "s" for seconds, "m" for minutes, or "h" for hours.

`cleanup_type` : This defines the strategy for clearing expired keys. There are three options - "active", "lazy", and "none". The default is "active".

Here's an example of a `config.json` file:

```json
{
  "ticker_time": "20s",
  "byte_slice_size": 1024,
  "default_ttl": "30m",
  "cleanup_type": "active"
}
```
Ensure that you update the path to `config.json` in the code according to your project's file structure.
