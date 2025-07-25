---
id: memory
title: Memory
---


![Release](https://img.shields.io/github/v/tag/gofiber/storage?filter=memory*)
[![Discord](https://img.shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](https://gofiber.io/discord)
![Test](https://img.shields.io/github/actions/workflow/status/gofiber/storage/test-memory.yml?label=Tests)

An in-memory storage driver.

### Table of Contents
- [Signatures](#signatures)
- [Installation](#installation)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)


### Signatures
```go
func New(config ...Config) Storage
func (s *Storage) Get(key string) ([]byte, error)
func (s *Storage) GetWithContext(ctx context.Context, key string) ([]byte, error)
func (s *Storage) Set(key string, val []byte, exp time.Duration) error
func (s *Storage) SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error
func (s *Storage) Delete(key string) error
func (s *Storage) DeleteWithContext(ctx context.Context, key string) error
func (s *Storage) Reset() error
func (s *Storage) ResetWithContext(ctx context.Context) error
func (s *Storage) Close() error
func (s *Storage) Conn() map[string]entry
func (s *Storage) Keys() ([][]byte, error)
```

**Note:** The context methods are dummy methods and don't have any functionality, as memory storage does not support context cancellation. They are provided for compliance with the Fiber storage interface.

### Installation
Memory is tested on the 2 last [Go versions](https://golang.org/dl/) with support for modules. So make sure to initialize one first if you didn't do that yet:
```bash
go mod init github.com/<user>/<repo>
```
And then install the memory implementation:
```bash
go get github.com/gofiber/storage/memory/v2
```

### Examples
Import the storage package.
```go
import "github.com/gofiber/storage/memory/v2"
```

You can use the following possibilities to create a storage:
```go
// Initialize default config
store := memory.New()

// Initialize custom config
store := memory.New(memory.Config{
	GCInterval: 10 * time.Second,
})
```

### Config
```go
type Config struct {
	// Time before deleting expired keys
	//
	// Default is 10 * time.Second
	GCInterval time.Duration
}
```

### Default Config
```go
var ConfigDefault = Config{
	GCInterval: 10 * time.Second,
}
```
