import { defaultConfig } from './types.js';
/**
 * WebSocket client for connecting to the Go debug bar server
 */
export class DebugBarWebSocket {
    ws = null;
    config;
    messageHandlers = new Set();
    statusHandlers = new Set();
    reconnectAttempts = 0;
    reconnectTimeout = null;
    pingInterval = null;
    _status = 'disconnected';
    constructor(config) {
        this.config = { ...defaultConfig, ...config };
        if (this.config.autoConnect) {
            this.connect();
        }
    }
    /**
     * Current connection status
     */
    get status() {
        return this._status;
    }
    /**
     * Connect to the WebSocket server
     */
    connect() {
        if (this.ws?.readyState === WebSocket.OPEN) {
            return;
        }
        this.setStatus('connecting');
        try {
            this.ws = new WebSocket(this.config.wsUrl);
            this.setupEventHandlers();
        }
        catch (error) {
            this.setStatus('error', error instanceof Error ? error.message : 'Failed to connect');
            this.scheduleReconnect();
        }
    }
    /**
     * Disconnect from the WebSocket server
     */
    disconnect() {
        this.clearReconnectTimeout();
        this.clearPingInterval();
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.setStatus('disconnected');
    }
    /**
     * Subscribe to WebSocket messages
     */
    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => this.messageHandlers.delete(handler);
    }
    /**
     * Subscribe to connection status changes
     */
    onStatusChange(handler) {
        this.statusHandlers.add(handler);
        // Immediately call with current status
        handler(this._status);
        return () => this.statusHandlers.delete(handler);
    }
    /**
     * Send a message to the server
     */
    send(message) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        }
    }
    /**
     * Send a ping message
     */
    ping() {
        this.send({ type: 'ping', payload: null });
    }
    setupEventHandlers() {
        if (!this.ws)
            return;
        this.ws.onopen = () => {
            this.reconnectAttempts = 0;
            this.setStatus('connected');
            this.startPingInterval();
        };
        this.ws.onclose = (event) => {
            this.clearPingInterval();
            if (event.wasClean) {
                this.setStatus('disconnected');
            }
            else {
                this.setStatus('error', 'Connection lost');
                this.scheduleReconnect();
            }
        };
        this.ws.onerror = () => {
            this.setStatus('error', 'WebSocket error');
        };
        this.ws.onmessage = (event) => {
            this.handleMessage(event.data);
        };
    }
    handleMessage(data) {
        // Handle multiple messages (server may batch them with newlines)
        const messages = data.split('\n').filter((m) => m.trim());
        for (const msgStr of messages) {
            try {
                const message = JSON.parse(msgStr);
                this.notifyMessageHandlers(message);
            }
            catch (error) {
                console.error('Failed to parse WebSocket message:', error);
            }
        }
    }
    notifyMessageHandlers(message) {
        for (const handler of this.messageHandlers) {
            try {
                handler(message);
            }
            catch (error) {
                console.error('Error in message handler:', error);
            }
        }
    }
    setStatus(status, error) {
        this._status = status;
        for (const handler of this.statusHandlers) {
            try {
                handler(status, error);
            }
            catch (err) {
                console.error('Error in status handler:', err);
            }
        }
    }
    scheduleReconnect() {
        if (!this.config.autoReconnect)
            return;
        if (this.config.maxReconnectAttempts > 0 &&
            this.reconnectAttempts >= this.config.maxReconnectAttempts) {
            this.setStatus('error', 'Max reconnect attempts reached');
            return;
        }
        this.clearReconnectTimeout();
        this.reconnectTimeout = setTimeout(() => {
            this.reconnectAttempts++;
            this.connect();
        }, this.config.reconnectDelay);
    }
    clearReconnectTimeout() {
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
            this.reconnectTimeout = null;
        }
    }
    startPingInterval() {
        this.clearPingInterval();
        // Send ping every 30 seconds to keep connection alive
        this.pingInterval = setInterval(() => this.ping(), 30000);
    }
    clearPingInterval() {
        if (this.pingInterval) {
            clearInterval(this.pingInterval);
            this.pingInterval = null;
        }
    }
    /**
     * Clean up resources
     */
    destroy() {
        this.disconnect();
        this.messageHandlers.clear();
        this.statusHandlers.clear();
    }
}
/**
 * Parse a history message payload into RequestInfo array
 */
export function parseHistoryPayload(payload) {
    if (!Array.isArray(payload))
        return [];
    return payload;
}
/**
 * Parse a request message payload
 */
export function parseRequestPayload(payload) {
    if (!payload || typeof payload !== 'object')
        return null;
    return payload;
}
/**
 * Parse a query message payload
 */
export function parseQueryPayload(payload) {
    if (!payload || typeof payload !== 'object')
        return null;
    return payload;
}
/**
 * Parse an error message payload
 */
export function parseErrorPayload(payload) {
    if (!payload || typeof payload !== 'object')
        return null;
    return payload;
}
