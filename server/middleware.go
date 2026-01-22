package godebugbar

import (
	"bytes"
	"context"
	"io"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseWriter wraps gin.ResponseWriter to capture response size
type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (w *responseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *responseWriter) WriteString(s string) (int, error) {
	n, err := w.ResponseWriter.WriteString(s)
	w.size += n
	return n, err
}

// ginMiddleware creates the Gin middleware for request tracking
func (d *DebugBar) ginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !d.config.Enabled {
			c.Next()
			return
		}

		// Skip the debug bar WebSocket endpoint
		if c.Request.URL.Path == d.config.WebSocketPath {
			c.Next()
			return
		}

		startTime := time.Now()

		// Create request info
		reqInfo := &RequestInfo{
			ID:          uuid.New().String(),
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			StartTime:   startTime,
			Headers:     make(map[string]string),
			QueryParams: make(map[string]string),
			Queries:     make([]QueryInfo, 0),
			Errors:      make([]ErrorInfo, 0),
			ClientIP:    c.ClientIP(),
		}

		// Capture headers
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				reqInfo.Headers[key] = values[0]
			}
		}

		// Capture query parameters
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				reqInfo.QueryParams[key] = values[0]
			}
		}

		// Capture request body if enabled
		if d.config.CaptureRequestBody && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, int64(d.config.MaxBodySize)))
			if err == nil {
				reqInfo.RequestBody = string(bodyBytes)
				// Restore the body for the actual handler
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Store request info in context
		c.Set(string(DebugBarContextKey), reqInfo)

		// Also store in standard context for GORM integration
		ctx := context.WithValue(c.Request.Context(), DebugBarContextKey, reqInfo)
		c.Request = c.Request.WithContext(ctx)

		// Wrap response writer to capture size
		rw := &responseWriter{ResponseWriter: c.Writer, size: 0}
		c.Writer = rw

		// Broadcast request start
		d.broadcast(WebSocketMessage{
			Type:    MessageTypeRequest,
			Payload: reqInfo,
		})

		// Process request
		c.Next()

		// Calculate final metrics
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		reqInfo.EndTime = endTime
		reqInfo.Duration = duration
		reqInfo.DurationMs = float64(duration.Nanoseconds()) / 1e6
		reqInfo.StatusCode = c.Writer.Status()
		reqInfo.ResponseSize = rw.size

		// Capture memory usage
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		reqInfo.MemoryUsage = memStats.Alloc

		// Store completed request
		d.storeRequest(reqInfo)

		// Broadcast request completion
		d.broadcast(WebSocketMessage{
			Type:    MessageTypeRequestEnd,
			Payload: reqInfo,
		})
	}
}
