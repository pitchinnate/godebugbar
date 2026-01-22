import type { DebugBarConfig, WebSocketMessage, RequestInfo, QueryInfo, ErrorInfo, ConnectionStatus } from './types.js';
export type MessageHandler = (message: WebSocketMessage) => void;
export type StatusHandler = (status: ConnectionStatus, error?: string) => void;
/**
 * WebSocket client for connecting to the Go debug bar server
 */
export declare class DebugBarWebSocket {
    private ws;
    private config;
    private messageHandlers;
    private statusHandlers;
    private reconnectAttempts;
    private reconnectTimeout;
    private pingInterval;
    private _status;
    constructor(config: DebugBarConfig);
    /**
     * Current connection status
     */
    get status(): ConnectionStatus;
    /**
     * Connect to the WebSocket server
     */
    connect(): void;
    /**
     * Disconnect from the WebSocket server
     */
    disconnect(): void;
    /**
     * Subscribe to WebSocket messages
     */
    onMessage(handler: MessageHandler): () => void;
    /**
     * Subscribe to connection status changes
     */
    onStatusChange(handler: StatusHandler): () => void;
    /**
     * Send a message to the server
     */
    send(message: WebSocketMessage): void;
    /**
     * Send a ping message
     */
    ping(): void;
    private setupEventHandlers;
    private handleMessage;
    private notifyMessageHandlers;
    private setStatus;
    private scheduleReconnect;
    private clearReconnectTimeout;
    private startPingInterval;
    private clearPingInterval;
    /**
     * Clean up resources
     */
    destroy(): void;
}
/**
 * Parse a history message payload into RequestInfo array
 */
export declare function parseHistoryPayload(payload: unknown): RequestInfo[];
/**
 * Parse a request message payload
 */
export declare function parseRequestPayload(payload: unknown): RequestInfo | null;
/**
 * Parse a query message payload
 */
export declare function parseQueryPayload(payload: unknown): QueryInfo | null;
/**
 * Parse an error message payload
 */
export declare function parseErrorPayload(payload: unknown): ErrorInfo | null;
