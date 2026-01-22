# Go Debug Bar

A debug bar package for Go web applications, similar to Laravel's debug bar. Supports Gin and GORM with real-time WebSocket streaming.

## Features

- **Request Tracking** - HTTP method, path, headers, query parameters, request body, response status, duration, memory usage
- **Database Query Tracking** - All GORM operations with query text, parameters, duration, rows affected, and source location
- **Error Logging** - Multiple severity levels (exception, warning, notice, debug) with full stack traces
- **Real-time WebSocket** - Stream debug data to connected clients instantly
- **Request History** - Configurable storage of recent requests

## Installation

```bash
go get github.com/user/godebugbar
```

## Quick Start

```go
package main

import (
    "github.com/gin-gonic/gin"
    godebugbar "github.com/user/godebugbar"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // Initialize the debug bar with default settings
    debugBar := godebugbar.NewWithDefaults()

    // Initialize GORM
    db, _ := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

    // Register the GORM plugin for query tracking
    db.Use(debugBar.GormPlugin())

    // Initialize Gin
    r := gin.Default()

    // Add debug bar middleware
    r.Use(debugBar.Middleware())

    // Register WebSocket endpoint
    debugBar.RegisterRoutes(r)

    // Your routes here
    r.GET("/", func(c *gin.Context) {
        // Queries are automatically tracked when using context
        var users []User
        db.WithContext(c.Request.Context()).Find(&users)
        c.JSON(200, users)
    })

    r.Run(":8080")
}
```

## Configuration

```go
debugBar := godebugbar.New(godebugbar.Config{
    // Enable or disable the debug bar
    Enabled: true,

    // WebSocket endpoint path
    WebSocketPath: "/_debugbar/ws",

    // Maximum number of requests to keep in history
    MaxRequests: 100,

    // Capture request bodies
    CaptureRequestBody: true,

    // Maximum request body size to capture (bytes)
    MaxBodySize: 64 * 1024, // 64KB

    // Allowed origins for WebSocket CORS
    AllowedOrigins: []string{"*"},
})
```

## Usage

### Request Tracking

Request tracking is automatic once the middleware is added:

```go
r.Use(debugBar.Middleware())
```

Each request captures:
- Request ID, method, path, client IP
- Headers and query parameters
- Request body (if enabled)
- Response status code and size
- Duration and memory usage
- Associated database queries
- Any logged errors

### Database Query Tracking

Register the GORM plugin to track all database operations:

```go
db.Use(debugBar.GormPlugin())
```

**Important:** Always use `WithContext` to associate queries with the current request:

```go
r.GET("/users", func(c *gin.Context) {
    var users []User
    db.WithContext(c.Request.Context()).Find(&users)
    c.JSON(200, users)
})
```

Tracked query information:
- SQL query text
- Query parameters
- Execution duration
- Rows affected
- Errors (if any)
- Source file and line number

### Error Logging

Log errors at different severity levels:

```go
// Log an error with stack trace
debugBar.LogError(c, err)

// Log with additional context
debugBar.LogErrorWithContext(c, err, map[string]any{
    "user_id": 123,
    "action":  "create_order",
})

// Log warnings
debugBar.LogWarning(c, "Cache miss for key: user_123")

// Log notices
debugBar.LogNotice(c, "User not found, creating new record")

// Log debug messages
debugBar.LogDebug(c, "Processing started")
debugBar.LogDebugWithContext(c, "Step completed", map[string]any{
    "step": 1,
    "data": someData,
})
```

### Custom Data

Add arbitrary data to the current request:

```go
debugBar.AddCustomData(c, "user_id", 123)
debugBar.AddCustomData(c, "permissions", []string{"read", "write"})
```

### Recovery Middleware

Capture panics in the debug bar:

```go
r.Use(debugBar.RecoveryMiddleware())
```

### Gin Error Handler

Automatically log Gin errors:

```go
r.Use(debugBar.ErrorHandler())
```

## WebSocket API

Connect to the WebSocket endpoint to receive real-time debug data:

```
ws://localhost:8080/_debugbar/ws
```

### Message Types

| Type | Description |
|------|-------------|
| `history` | Sent on connection with all stored requests |
| `request` | Sent when a new request starts |
| `request_end` | Sent when a request completes |
| `query` | Sent for each database query |
| `error` | Sent when an error is logged |
| `ping` / `pong` | Keep-alive messages |

### Message Format

```json
{
    "type": "request_end",
    "payload": {
        "id": "uuid",
        "method": "GET",
        "path": "/users",
        "status_code": 200,
        "duration_ms": 12.5,
        "queries": [...],
        "errors": [...]
    }
}
```

## API Reference

### DebugBar Methods

| Method | Description |
|--------|-------------|
| `New(config Config)` | Create with custom configuration |
| `NewWithDefaults()` | Create with default configuration |
| `Middleware()` | Returns Gin middleware |
| `GormPlugin()` | Returns GORM plugin |
| `RegisterRoutes(r *gin.Engine)` | Register WebSocket endpoint |
| `LogError(c, err)` | Log an error |
| `LogErrorWithContext(c, err, ctx)` | Log error with context |
| `LogWarning(c, message)` | Log a warning |
| `LogNotice(c, message)` | Log a notice |
| `LogDebug(c, message)` | Log a debug message |
| `AddCustomData(c, key, value)` | Add custom data to request |
| `GetRequestInfo(c)` | Get current request info |
| `GetHistory()` | Get all stored requests |
| `GetRecentHistory(n)` | Get last n requests |
| `ClearHistory()` | Clear stored requests |
| `IsEnabled()` | Check if enabled |
| `SetEnabled(bool)` | Enable or disable |

## Production Usage

Disable the debug bar in production:

```go
debugBar := godebugbar.New(godebugbar.Config{
    Enabled: os.Getenv("ENV") != "production",
    // ... other config
})
```

Or disable at runtime:

```go
debugBar.SetEnabled(false)
```

## License

MIT
