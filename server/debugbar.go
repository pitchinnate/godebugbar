package godebugbar

import (
	"context"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// DebugBarContextKey is the context key for storing request info
	DebugBarContextKey ContextKey = "debugbar_request"
)

// DebugBar is the main debug bar instance
type DebugBar struct {
	config    Config
	store     *RequestStore
	wsHub     *WebSocketHub
	mu        sync.RWMutex
}

// New creates a new DebugBar instance with the given configuration
func New(config Config) *DebugBar {
	db := &DebugBar{
		config: config,
		store:  NewRequestStore(config.MaxRequests),
		wsHub:  NewWebSocketHub(),
	}

	if config.Enabled {
		go db.wsHub.Run()
	}

	return db
}

// NewWithDefaults creates a new DebugBar instance with default configuration
func NewWithDefaults() *DebugBar {
	return New(DefaultConfig())
}

// Middleware returns the Gin middleware for request tracking
func (d *DebugBar) Middleware() gin.HandlerFunc {
	return d.ginMiddleware()
}

// GormPlugin returns the GORM plugin for query tracking
func (d *DebugBar) GormPlugin() gorm.Plugin {
	return &GormDebugBarPlugin{debugBar: d}
}

// RegisterRoutes registers the WebSocket endpoint with the Gin router
func (d *DebugBar) RegisterRoutes(router *gin.Engine) {
	if !d.config.Enabled {
		return
	}

	router.GET(d.config.WebSocketPath, d.handleWebSocket)
}

// AddError adds an error to the current request context
func (d *DebugBar) AddError(c *gin.Context, err error, errType string, ctx map[string]any) {
	reqInfo := d.GetRequestInfo(c)
	if reqInfo == nil {
		return
	}

	errorInfo := ErrorInfo{
		ID:        uuid.New().String(),
		RequestID: reqInfo.ID,
		Message:   err.Error(),
		Type:      errType,
		Context:   ctx,
	}

	// Capture stack trace
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	errorInfo.Stack = string(buf[:n])

	d.mu.Lock()
	reqInfo.Errors = append(reqInfo.Errors, errorInfo)
	d.mu.Unlock()

	// Broadcast error to WebSocket clients
	d.broadcast(WebSocketMessage{
		Type:    MessageTypeError,
		Payload: errorInfo,
	})
}

// AddCustomData adds custom data to the current request
func (d *DebugBar) AddCustomData(c *gin.Context, key string, value any) {
	reqInfo := d.GetRequestInfo(c)
	if reqInfo == nil {
		return
	}

	d.mu.Lock()
	if reqInfo.CustomData == nil {
		reqInfo.CustomData = make(map[string]any)
	}
	reqInfo.CustomData[key] = value
	d.mu.Unlock()
}

// GetRequestInfo retrieves the request info from the Gin context
func (d *DebugBar) GetRequestInfo(c *gin.Context) *RequestInfo {
	if val, exists := c.Get(string(DebugBarContextKey)); exists {
		if reqInfo, ok := val.(*RequestInfo); ok {
			return reqInfo
		}
	}
	return nil
}

// GetRequestInfoFromContext retrieves the request info from a standard context
func (d *DebugBar) GetRequestInfoFromContext(ctx context.Context) *RequestInfo {
	if val := ctx.Value(DebugBarContextKey); val != nil {
		if reqInfo, ok := val.(*RequestInfo); ok {
			return reqInfo
		}
	}
	return nil
}

// GetHistory returns the request history
func (d *DebugBar) GetHistory() []*RequestInfo {
	return d.store.GetAll()
}

// GetRecentHistory returns the most recent n requests
func (d *DebugBar) GetRecentHistory(n int) []*RequestInfo {
	return d.store.GetRecent(n)
}

// ClearHistory clears all stored requests
func (d *DebugBar) ClearHistory() {
	d.store.Clear()
}

// IsEnabled returns whether the debug bar is enabled
func (d *DebugBar) IsEnabled() bool {
	return d.config.Enabled
}

// SetEnabled enables or disables the debug bar
func (d *DebugBar) SetEnabled(enabled bool) {
	d.mu.Lock()
	d.config.Enabled = enabled
	d.mu.Unlock()
}

// broadcast sends a message to all connected WebSocket clients
func (d *DebugBar) broadcast(msg WebSocketMessage) {
	if d.wsHub != nil {
		d.wsHub.Broadcast(msg)
	}
}

// addQuery adds a query to the current request
func (d *DebugBar) addQuery(ctx context.Context, query QueryInfo) {
	reqInfo := d.GetRequestInfoFromContext(ctx)
	if reqInfo == nil {
		return
	}

	query.RequestID = reqInfo.ID

	d.mu.Lock()
	reqInfo.Queries = append(reqInfo.Queries, query)
	d.mu.Unlock()

	// Broadcast query to WebSocket clients
	d.broadcast(WebSocketMessage{
		Type:    MessageTypeQuery,
		Payload: query,
	})
}

// storeRequest stores a completed request
func (d *DebugBar) storeRequest(req *RequestInfo) {
	d.store.Add(req)
}
