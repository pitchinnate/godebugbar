/**
 * Information about an HTTP request
 */
export interface RequestInfo {
    id: string;
    method: string;
    path: string;
    status_code: number;
    duration: number;
    duration_ms: number;
    start_time: string;
    end_time: string;
    headers: Record<string, string>;
    query_params: Record<string, string>;
    request_body?: string;
    response_size: number;
    client_ip: string;
    queries: QueryInfo[];
    errors: ErrorInfo[];
    memory_usage: number;
    custom_data?: Record<string, unknown>;
}
/**
 * Information about a database query
 */
export interface QueryInfo {
    id: string;
    request_id: string;
    query: string;
    args?: unknown[];
    duration: number;
    duration_ms: number;
    rows_affected: number;
    error?: string;
    start_time: string;
    source?: string;
}
/**
 * Information about an error
 */
export interface ErrorInfo {
    id: string;
    request_id: string;
    message: string;
    stack?: string;
    type: ErrorType;
    timestamp: string;
    context?: Record<string, unknown>;
}
/**
 * Error severity types
 */
export type ErrorType = 'exception' | 'warning' | 'notice' | 'debug';
/**
 * WebSocket message from the server
 */
export interface WebSocketMessage {
    type: MessageType;
    payload: unknown;
}
/**
 * Message types for WebSocket communication
 */
export type MessageType = 'request' | 'query' | 'error' | 'request_end' | 'history' | 'ping' | 'pong';
/**
 * Debug bar configuration options
 */
export interface DebugBarConfig {
    /** WebSocket server URL (e.g., 'ws://localhost:8080/_debugbar/ws') */
    wsUrl: string;
    /** Auto-connect on initialization */
    autoConnect?: boolean;
    /** Auto-reconnect on disconnect */
    autoReconnect?: boolean;
    /** Reconnect delay in milliseconds */
    reconnectDelay?: number;
    /** Maximum reconnect attempts (0 = infinite) */
    maxReconnectAttempts?: number;
    /** Maximum requests to keep in history */
    maxRequests?: number;
    /** Start minimized */
    startMinimized?: boolean;
    /** Position of the debug bar */
    position?: 'bottom' | 'top' | 'left' | 'right';
    /** Default height/width in pixels */
    defaultSize?: number;
}
/**
 * Default configuration values
 */
export declare const defaultConfig: Required<DebugBarConfig>;
/**
 * Connection status
 */
export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'error';
/**
 * Debug bar state
 */
export interface DebugBarState {
    /** Current connection status */
    status: ConnectionStatus;
    /** All tracked requests */
    requests: RequestInfo[];
    /** Currently selected request */
    selectedRequest: RequestInfo | null;
    /** Currently active tab */
    activeTab: TabType;
    /** Whether the bar is minimized */
    minimized: boolean;
    /** Error message if connection failed */
    errorMessage?: string;
}
/**
 * Available tabs in the debug bar
 */
export type TabType = 'requests' | 'queries' | 'errors' | 'timeline';
