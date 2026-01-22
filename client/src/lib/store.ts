import { writable, derived, get } from 'svelte/store';
import type {
	DebugBarConfig,
	DebugBarState,
	RequestInfo,
	QueryInfo,
	ErrorInfo,
	ConnectionStatus,
	TabType
} from './types.js';
import { defaultConfig } from './types.js';
import {
	DebugBarWebSocket,
	parseHistoryPayload,
	parseRequestPayload,
	parseQueryPayload,
	parseErrorPayload
} from './websocket.js';

/**
 * Creates a debug bar store with WebSocket connection
 */
export function createDebugBarStore(config: DebugBarConfig) {
	const fullConfig = { ...defaultConfig, ...config };

	// Internal state
	const state = writable<DebugBarState>({
		status: 'disconnected',
		requests: [],
		selectedRequest: null,
		activeTab: 'requests',
		minimized: fullConfig.startMinimized
	});

	let wsClient: DebugBarWebSocket | null = null;

	// Initialize WebSocket connection
	function init(): void {
		if (wsClient) return;

		wsClient = new DebugBarWebSocket(fullConfig);

		// Handle status changes
		wsClient.onStatusChange((status, error) => {
			state.update((s) => ({
				...s,
				status,
				errorMessage: error
			}));
		});

		// Handle messages
		wsClient.onMessage((message) => {
			switch (message.type) {
				case 'history': {
					const requests = parseHistoryPayload(message.payload);
					state.update((s) => ({
						...s,
						requests: requests.slice(-fullConfig.maxRequests)
					}));
					break;
				}

				case 'request': {
					const request = parseRequestPayload(message.payload);
					if (request) {
						state.update((s) => {
							// Add new request (it's still in progress)
							const exists = s.requests.find((r) => r.id === request.id);
							if (exists) return s;

							const requests = [...s.requests, request].slice(-fullConfig.maxRequests);
							return { ...s, requests };
						});
					}
					break;
				}

				case 'request_end': {
					const request = parseRequestPayload(message.payload);
					if (request) {
						state.update((s) => {
							// Update existing request with final data
							const requests = s.requests.map((r) => (r.id === request.id ? request : r));
							const selectedRequest =
								s.selectedRequest?.id === request.id ? request : s.selectedRequest;
							return { ...s, requests, selectedRequest };
						});
					}
					break;
				}

				case 'query': {
					const query = parseQueryPayload(message.payload);
					if (query) {
						state.update((s) => {
							const requests = s.requests.map((r) => {
								if (r.id === query.request_id) {
									return {
										...r,
										queries: [...(r.queries || []), query]
									};
								}
								return r;
							});

							const selectedRequest =
								s.selectedRequest?.id === query.request_id
									? {
											...s.selectedRequest,
											queries: [...(s.selectedRequest.queries || []), query]
										}
									: s.selectedRequest;

							return { ...s, requests, selectedRequest };
						});
					}
					break;
				}

				case 'error': {
					const error = parseErrorPayload(message.payload);
					if (error) {
						state.update((s) => {
							const requests = s.requests.map((r) => {
								if (r.id === error.request_id) {
									return {
										...r,
										errors: [...(r.errors || []), error]
									};
								}
								return r;
							});

							const selectedRequest =
								s.selectedRequest?.id === error.request_id
									? {
											...s.selectedRequest,
											errors: [...(s.selectedRequest.errors || []), error]
										}
									: s.selectedRequest;

							return { ...s, requests, selectedRequest };
						});
					}
					break;
				}
			}
		});
	}

	// Actions
	function connect(): void {
		wsClient?.connect();
	}

	function disconnect(): void {
		wsClient?.disconnect();
	}

	function selectRequest(request: RequestInfo | null): void {
		state.update((s) => ({ ...s, selectedRequest: request }));
	}

	function setActiveTab(tab: TabType): void {
		state.update((s) => ({ ...s, activeTab: tab }));
	}

	function toggleMinimized(): void {
		state.update((s) => ({ ...s, minimized: !s.minimized }));
	}

	function setMinimized(minimized: boolean): void {
		state.update((s) => ({ ...s, minimized }));
	}

	function clearRequests(): void {
		state.update((s) => ({ ...s, requests: [], selectedRequest: null }));
	}

	function destroy(): void {
		wsClient?.destroy();
		wsClient = null;
	}

	// Derived stores
	const status = derived(state, ($state) => $state.status);
	const requests = derived(state, ($state) => $state.requests);
	const selectedRequest = derived(state, ($state) => $state.selectedRequest);
	const activeTab = derived(state, ($state) => $state.activeTab);
	const minimized = derived(state, ($state) => $state.minimized);

	const totalQueries = derived(state, ($state) =>
		$state.requests.reduce((sum, r) => sum + (r.queries?.length || 0), 0)
	);

	const totalErrors = derived(state, ($state) =>
		$state.requests.reduce((sum, r) => sum + (r.errors?.length || 0), 0)
	);

	const selectedQueries = derived(state, ($state) => $state.selectedRequest?.queries || []);

	const selectedErrors = derived(state, ($state) => $state.selectedRequest?.errors || []);

	// Auto-initialize if autoConnect is true
	if (fullConfig.autoConnect) {
		init();
	}

	return {
		// State
		subscribe: state.subscribe,
		status,
		requests,
		selectedRequest,
		activeTab,
		minimized,
		totalQueries,
		totalErrors,
		selectedQueries,
		selectedErrors,

		// Actions
		init,
		connect,
		disconnect,
		selectRequest,
		setActiveTab,
		toggleMinimized,
		setMinimized,
		clearRequests,
		destroy
	};
}

export type DebugBarStore = ReturnType<typeof createDebugBarStore>;
