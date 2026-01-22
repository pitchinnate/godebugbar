package godebugbar

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Error types
const (
	ErrorTypeException = "exception"
	ErrorTypeWarning   = "warning"
	ErrorTypeNotice    = "notice"
	ErrorTypeDebug     = "debug"
)

// LogError logs an error to the debug bar
func (d *DebugBar) LogError(c *gin.Context, err error) {
	d.logError(c, err, ErrorTypeException, nil, 2)
}

// LogErrorWithContext logs an error with additional context
func (d *DebugBar) LogErrorWithContext(c *gin.Context, err error, ctx map[string]any) {
	d.logError(c, err, ErrorTypeException, ctx, 2)
}

// LogWarning logs a warning to the debug bar
func (d *DebugBar) LogWarning(c *gin.Context, message string) {
	d.logError(c, fmt.Errorf("%s", message), ErrorTypeWarning, nil, 2)
}

// LogWarningWithContext logs a warning with additional context
func (d *DebugBar) LogWarningWithContext(c *gin.Context, message string, ctx map[string]any) {
	d.logError(c, fmt.Errorf("%s", message), ErrorTypeWarning, ctx, 2)
}

// LogNotice logs a notice to the debug bar
func (d *DebugBar) LogNotice(c *gin.Context, message string) {
	d.logError(c, fmt.Errorf("%s", message), ErrorTypeNotice, nil, 2)
}

// LogDebug logs a debug message to the debug bar
func (d *DebugBar) LogDebug(c *gin.Context, message string) {
	d.logError(c, fmt.Errorf("%s", message), ErrorTypeDebug, nil, 2)
}

// LogDebugWithContext logs a debug message with additional context
func (d *DebugBar) LogDebugWithContext(c *gin.Context, message string, ctx map[string]any) {
	d.logError(c, fmt.Errorf("%s", message), ErrorTypeDebug, ctx, 2)
}

// logError is the internal method for logging errors
func (d *DebugBar) logError(c *gin.Context, err error, errType string, ctx map[string]any, skip int) {
	if !d.config.Enabled {
		return
	}

	reqInfo := d.GetRequestInfo(c)
	if reqInfo == nil {
		return
	}

	errorInfo := ErrorInfo{
		ID:        uuid.New().String(),
		RequestID: reqInfo.ID,
		Message:   err.Error(),
		Type:      errType,
		Timestamp: time.Now(),
		Context:   ctx,
	}

	// Capture stack trace
	errorInfo.Stack = captureStackTrace(skip + 1)

	d.mu.Lock()
	reqInfo.Errors = append(reqInfo.Errors, errorInfo)
	d.mu.Unlock()

	// Broadcast error to WebSocket clients
	d.broadcast(WebSocketMessage{
		Type:    MessageTypeError,
		Payload: errorInfo,
	})
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) string {
	const maxStackSize = 4096
	buf := make([]byte, maxStackSize)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// RecoveryMiddleware returns a Gin middleware that recovers from panics
// and logs them to the debug bar
func (d *DebugBar) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic to debug bar
				var panicErr error
				switch e := err.(type) {
				case error:
					panicErr = e
				default:
					panicErr = fmt.Errorf("%v", e)
				}

				d.logError(c, panicErr, ErrorTypeException, map[string]any{
					"panic": true,
				}, 4)

				// Re-panic to let Gin's default recovery handle it
				// or handle it here if you want custom behavior
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// ErrorHandler returns a function that can be used with Gin's error handling
func (d *DebugBar) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors after the request is processed
		for _, ginErr := range c.Errors {
			d.logError(c, ginErr.Err, ErrorTypeException, map[string]any{
				"gin_error_type": ginErr.Type,
				"gin_error_meta": ginErr.Meta,
			}, 2)
		}
	}
}
