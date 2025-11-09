# Atomic

A collection of thread-safe, generic data structures and utilities for Go. Built with modern Go generics and designed for concurrent applications.

## Overview

Atomic provides high-performance, concurrent-safe data structures with rich utility methods. All data structures are designed to be thread-safe by default, using efficient read-write locking mechanisms.

## Packages

### 📦 [Collection](./collection)

A powerful, generic, thread-safe map-like data structure with 50+ utility methods inspired by JavaScript's Map and Array APIs.

**Features:**
- Thread-safe concurrent access with RWMutex
- Generic support for any comparable key type and any value type
- Functional programming methods (Map, Filter, Reduce, FlatMap)
- Array-like access with positive/negative indexing
- Set operations (Union, Intersection, Difference)
- Sorting, reversing, and randomization
- Query methods (Find, Partition, Some, Every)

[→ View Collection Documentation](./collection/README.md)

## Installation

```bash
# Install specific package
go get github.com/kolosys/atomic/collection

# Or add to your go.mod
require github.com/kolosys/atomic/collection latest
```

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/kolosys/atomic/collection"
)

func main() {
    // Create a thread-safe collection
    users := collection.New[string, User]()

    // Safely add items from multiple goroutines
    users.Set("alice", User{Name: "Alice", Age: 30})
    users.Set("bob", User{Name: "Bob", Age: 25})

    // Filter and transform
    adults := users.Filter(func(user User, id string, c *collection.Collection[string, User]) bool {
        return user.Age >= 18
    })

    // Map to a different type
    names := collection.MapCollection(users, func(user User, id string, c *collection.Collection[string, User]) string {
        return user.Name
    })

    fmt.Printf("Found %d adults\n", adults.Size())
    fmt.Printf("Names: %v\n", names)
}
```

## Design Philosophy

1. **Thread-Safe by Default**: All data structures are designed for concurrent use without requiring external synchronization
2. **Generic First**: Built with Go 1.18+ generics for type safety and flexibility
3. **Functional Patterns**: Rich APIs inspired by functional programming languages
4. **Performance**: Efficient implementations using read-write locks for concurrent reads
5. **Ergonomic APIs**: Chainable methods and intuitive interfaces

## Requirements

- Go 1.24 or later

## Performance

All atomic data structures are optimized for concurrent access:

- **Concurrent Reads**: Multiple goroutines can read simultaneously using `RLock()`
- **Safe Writes**: Write operations are protected with exclusive locks
- **Minimal Lock Contention**: Fine-grained locking strategies where applicable
- **Zero Dependencies**: Pure Go implementations with no external dependencies

## Project Structure

```
atomic/
├── collection/          # Thread-safe generic map with utility methods
│   ├── collection.go
│   ├── collection_mappers.go
│   ├── collection_test.go
│   └── README.md
├── go.work             # Go workspace configuration
├── LICENSE             # MIT License
└── README.md           # This file
```

## Roadmap

Future packages planned for the atomic library:

- **Set**: Thread-safe generic set with mathematical set operations
- **List**: Thread-safe generic slice with array-like operations
- **Queue**: Thread-safe FIFO queue with blocking/non-blocking operations
- **Stack**: Thread-safe LIFO stack
- **Cache**: Thread-safe LRU/LFU cache implementations
- **Counter**: Thread-safe counter with atomic operations
- **Registry**: Thread-safe service/dependency registry

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test -v ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development

```bash
# Clone the repository
git clone https://github.com/kolosys/atomic.git
cd atomic

# Run all tests
go test -v ./...

# Run tests with race detector
go test -race -v ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Check test coverage
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://pkg.go.dev/github.com/kolosys/atomic)
- 🐛 [Issue Tracker](https://github.com/kolosys/atomic/issues)
- 💬 [Discussions](https://github.com/kolosys/atomic/discussions)

## Acknowledgments

Inspired by:
- Discord.js Collection API
- JavaScript Array and Map methods
- Java's Concurrent Collections
- Rust's standard collection types

---

Made with ❤️ by [Kolosys](https://github.com/kolosys)
