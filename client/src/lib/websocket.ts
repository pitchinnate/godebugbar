import type {
	DebugBarConfig,
	WebSocketMessage,
	RequestInfo,
	QueryInfo,
	ErrorInfo,
	ConnectionStatus
} from './types.js';
import { defaultConfig } from './types.js';

export type MessageHandler = (message: WebSocketMessage) => void;
export type StatusHandler = (status: ConnectionStatus, error?: string) => void;

/**
 * WebSocket client for connecting to the Go debug bar server
 */
export class DebugBarWebSocket {
	private ws: WebSocket | null = null;
	private config: Required<DebugBarConfig>;
	private messageHandlers: Set<MessageHandler> = new Set();
	private statusHandlers: Set<StatusHandler> = new Set();
	private reconnectAttempts = 0;
	private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	private pingInterval: ReturnType<typeof setInterval> | null = null;
	private _status: ConnectionStatus = 'disconnected';

	constructor(config: DebugBarConfig) {
		this.config = { ...defaultConfig, ...config };

		if (this.config.autoConnect) {
			this.connect();
		}
	}

	/**
	 * Current connection status
	 */
	get status(): ConnectionStatus {
		return this._status;
	}

	/**
	 * Connect to the WebSocket server
	 */
	connect(): void {
		if (this.ws?.readyState === WebSocket.OPEN) {
			return;
		}

		this.setStatus('connecting');

		try {
			this.ws = new WebSocket(this.config.wsUrl);
			this.setupEventHandlers();
		} catch (error) {
			this.setStatus('error', error instanceof Error ? error.message : 'Failed to connect');
			this.scheduleReconnect();
		}
	}

	/**
	 * Disconnect from the WebSocket server
	 */
	disconnect(): void {
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
	onMessage(handler: MessageHandler): () => void {
		this.messageHandlers.add(handler);
		return () => this.messageHandlers.delete(handler);
	}

	/**
	 * Subscribe to connection status changes
	 */
	onStatusChange(handler: StatusHandler): () => void {
		this.statusHandlers.add(handler);
		// Immediately call with current status
		handler(this._status);
		return () => this.statusHandlers.delete(handler);
	}

	/**
	 * Send a message to the server
	 */
	send(message: WebSocketMessage): void {
		if (this.ws?.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(message));
		}
	}

	/**
	 * Send a ping message
	 */
	ping(): void {
		this.send({ type: 'ping', payload: null });
	}

	private setupEventHandlers(): void {
		if (!this.ws) return;

		this.ws.onopen = () => {
			this.reconnectAttempts = 0;
			this.setStatus('connected');
			this.startPingInterval();
		};

		this.ws.onclose = (event) => {
			this.clearPingInterval();

			if (event.wasClean) {
				this.setStatus('disconnected');
			} else {
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

	private handleMessage(data: string): void {
		// Handle multiple messages (server may batch them with newlines)
		const messages = data.split('\n').filter((m) => m.trim());

		for (const msgStr of messages) {
			try {
				const message = JSON.parse(msgStr) as WebSocketMessage;
				this.notifyMessageHandlers(message);
			} catch (error) {
				console.error('Failed to parse WebSocket message:', error);
			}
		}
	}

	private notifyMessageHandlers(message: WebSocketMessage): void {
		for (const handler of this.messageHandlers) {
			try {
				handler(message);
			} catch (error) {
				console.error('Error in message handler:', error);
			}
		}
	}

	private setStatus(status: ConnectionStatus, error?: string): void {
		this._status = status;
		for (const handler of this.statusHandlers) {
			try {
				handler(status, error);
			} catch (err) {
				console.error('Error in status handler:', err);
			}
		}
	}

	private scheduleReconnect(): void {
		if (!this.config.autoReconnect) return;

		if (
			this.config.maxReconnectAttempts > 0 &&
			this.reconnectAttempts >= this.config.maxReconnectAttempts
		) {
			this.setStatus('error', 'Max reconnect attempts reached');
			return;
		}

		this.clearReconnectTimeout();

		this.reconnectTimeout = setTimeout(() => {
			this.reconnectAttempts++;
			this.connect();
		}, this.config.reconnectDelay);
	}

	private clearReconnectTimeout(): void {
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}
	}

	private startPingInterval(): void {
		this.clearPingInterval();
		// Send ping every 30 seconds to keep connection alive
		this.pingInterval = setInterval(() => this.ping(), 30000);
	}

	private clearPingInterval(): void {
		if (this.pingInterval) {
			clearInterval(this.pingInterval);
			this.pingInterval = null;
		}
	}

	/**
	 * Clean up resources
	 */
	destroy(): void {
		this.disconnect();
		this.messageHandlers.clear();
		this.statusHandlers.clear();
	}
}

/**
 * Parse a history message payload into RequestInfo array
 */
export function parseHistoryPayload(payload: unknown): RequestInfo[] {
	if (!Array.isArray(payload)) return [];
	return payload as RequestInfo[];
}

/**
 * Parse a request message payload
 */
export function parseRequestPayload(payload: unknown): RequestInfo | null {
	if (!payload || typeof payload !== 'object') return null;
	return payload as RequestInfo;
}

/**
 * Parse a query message payload
 */
export function parseQueryPayload(payload: unknown): QueryInfo | null {
	if (!payload || typeof payload !== 'object') return null;
	return payload as QueryInfo;
}

/**
 * Parse an error message payload
 */
export function parseErrorPayload(payload: unknown): ErrorInfo | null {
	if (!payload || typeof payload !== 'object') return null;
	return payload as ErrorInfo;
}
