package godebugbar

import (
	"sync"
	"time"
)

// RequestInfo holds information about an HTTP request
type RequestInfo struct {
	ID            string            `json:"id"`
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	StatusCode    int               `json:"status_code"`
	Duration      time.Duration     `json:"duration"`
	DurationMs    float64           `json:"duration_ms"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	Headers       map[string]string `json:"headers"`
	QueryParams   map[string]string `json:"query_params"`
	RequestBody   string            `json:"request_body,omitempty"`
	ResponseSize  int               `json:"response_size"`
	ClientIP      string            `json:"client_ip"`
	Queries       []QueryInfo       `json:"queries"`
	Errors        []ErrorInfo       `json:"errors"`
	MemoryUsage   uint64            `json:"memory_usage"`
	CustomData    map[string]any    `json:"custom_data,omitempty"`
}

// QueryInfo holds information about a database query
type QueryInfo struct {
	ID        string        `json:"id"`
	RequestID string        `json:"request_id"`
	Query     string        `json:"query"`
	Args      []any         `json:"args,omitempty"`
	Duration  time.Duration `json:"duration"`
	DurationMs float64      `json:"duration_ms"`
	RowsAffected int64      `json:"rows_affected"`
	Error     string        `json:"error,omitempty"`
	StartTime time.Time     `json:"start_time"`
	Source    string        `json:"source,omitempty"`
}

// ErrorInfo holds information about an error
type ErrorInfo struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	Message    string    `json:"message"`
	Stack      string    `json:"stack,omitempty"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	Context    map[string]any `json:"context,omitempty"`
}

// WebSocketMessage represents a message sent over WebSocket
type WebSocketMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Message types for WebSocket communication
const (
	MessageTypeRequest     = "request"
	MessageTypeQuery       = "query"
	MessageTypeError       = "error"
	MessageTypeRequestEnd  = "request_end"
	MessageTypeHistory     = "history"
	MessageTypePing        = "ping"
	MessageTypePong        = "pong"
)

// Config holds the debug bar configuration
type Config struct {
	// Enabled determines if the debug bar is active
	Enabled bool

	// WebSocketPath is the path where the WebSocket server will listen
	WebSocketPath string

	// MaxRequests is the maximum number of requests to keep in history
	MaxRequests int

	// CaptureRequestBody determines if request bodies should be captured
	CaptureRequestBody bool

	// MaxBodySize is the maximum size of request body to capture
	MaxBodySize int

	// AllowedOrigins for WebSocket CORS
	AllowedOrigins []string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Enabled:            true,
		WebSocketPath:      "/_debugbar/ws",
		MaxRequests:        100,
		CaptureRequestBody: true,
		MaxBodySize:        64 * 1024, // 64KB
		AllowedOrigins:     []string{"*"},
	}
}

// RequestStore stores request history with thread-safe access
type RequestStore struct {
	mu       sync.RWMutex
	requests []*RequestInfo
	maxSize  int
}

// NewRequestStore creates a new request store
func NewRequestStore(maxSize int) *RequestStore {
	return &RequestStore{
		requests: make([]*RequestInfo, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Add adds a request to the store
func (s *RequestStore) Add(req *RequestInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.requests) >= s.maxSize {
		// Remove oldest request
		s.requests = s.requests[1:]
	}
	s.requests = append(s.requests, req)
}

// Get returns a request by ID
func (s *RequestStore) Get(id string) *RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, req := range s.requests {
		if req.ID == id {
			return req
		}
	}
	return nil
}

// GetAll returns all stored requests
func (s *RequestStore) GetAll() []*RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*RequestInfo, len(s.requests))
	copy(result, s.requests)
	return result
}

// GetRecent returns the most recent n requests
func (s *RequestStore) GetRecent(n int) []*RequestInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if n > len(s.requests) {
		n = len(s.requests)
	}

	result := make([]*RequestInfo, n)
	copy(result, s.requests[len(s.requests)-n:])
	return result
}

// Clear removes all stored requests
func (s *RequestStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = make([]*RequestInfo, 0, s.maxSize)
}
