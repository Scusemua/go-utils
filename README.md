# go-utils
A collection of Go utility modules and helper functions designed to streamline common tasks in Go applications. This repository contains a range of modules for caching, configuration, logging, concurrency, and more, allowing you to efficiently handle repetitive tasks in your Go projects.

## Overview of Modules

### Cache
Lightweight, in-memory caching utility that provides a simple API for storing and retrieving values. Useful for caching expensive computations or external resource calls.

### Config
Configuration loader that supports parsing from multiple formats, including JSON, YAML, and environment variables. Helps centralize and manage your application's settings.

### Logger
Logging utility that offers structured logging with support for various logging levels. Easily configurable to suit different verbosity and output formats.

### MapReduce
Provides a functional approach to parallelizing work across multiple goroutines with map-reduce style processing. Ideal for tasks that can be broken down into smaller, independent units.

### Net
Network utilities, including helpers for checking port availability, resolving IPs, and more. Simplifies common network-related tasks in Go.

### Promise
A lightweight promise implementation for Go, enabling async-style programming. Useful for handling asynchronous operations cleanly and effectively.

### Sync
Concurrency utilities for synchronizing goroutines, including enhanced mutexes and wait groups. Designed to make concurrent programming in Go easier and safer.

## Installation
``` bash
go get github.com/Scusemua/go-utils
```

Each module can be imported individually, so you can use only what you need:
``` go
import (
    "github.com/Scusemua/go-utils/cache"
    "github.com/Scusemua/go-utils/config"
    "github.com/Scusemua/go-utils/logger"
    // Add other modules as needed
)
```

## Usage Examples
Here are some quick examples to get started with each module.

### Cache Example
``` go
import "github.com/Scusemua/go-utils/cache"

...
```

### Config Example
``` go
import "github.com/Scusemua/go-utils/config"

...
```

### Logger Example
``` go
import "github.com/Scusemua/go-utils/logger"

...
```

### MapReduce Example
``` go
import "github.com/Scusemua/go-utils/mapreduce"

...
```

### Net Example
``` go
import "github.com/Scusemua/go-utils/net"

...
```

### Promise Example
``` go
import "github.com/Scusemua/go-utils/promise"

p := promise.NewPromise()

go func() {
    res := doSomeWork()
    p.Resolve(res)
}
```

### Sync Example
``` go
import "github.com/Scusemua/go-utils/sync"

waitGroup := sync.WaitGroup()
wg.Add(1)

go func() {
    doSomething()
    wg.Done()
}

wg.Wait()
```

## Contributing
Contributions are welcome! Please feel free to open issues or submit pull requests.

### Fork the repository
Create a new branch (`git checkout -b feature/your-feature`)
Commit your changes (`git commit -am 'Add new feature'`)
Push to the branch (`git push origin feature/your-feature`)
Open a pull request

## License
This project is licensed under the MIT License. See the LICENSE file for details.
