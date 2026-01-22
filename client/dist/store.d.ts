import type { DebugBarConfig, DebugBarState, RequestInfo, QueryInfo, ErrorInfo, ConnectionStatus, TabType } from './types.js';
/**
 * Creates a debug bar store with WebSocket connection
 */
export declare function createDebugBarStore(config: DebugBarConfig): {
    subscribe: (this: void, run: import("svelte/store").Subscriber<DebugBarState>, invalidate?: import("svelte/store").Invalidator<DebugBarState> | undefined) => import("svelte/store").Unsubscriber;
    status: import("svelte/store").Readable<ConnectionStatus>;
    requests: import("svelte/store").Readable<RequestInfo[]>;
    selectedRequest: import("svelte/store").Readable<RequestInfo | null>;
    activeTab: import("svelte/store").Readable<TabType>;
    minimized: import("svelte/store").Readable<boolean>;
    totalQueries: import("svelte/store").Readable<number>;
    totalErrors: import("svelte/store").Readable<number>;
    selectedQueries: import("svelte/store").Readable<QueryInfo[]>;
    selectedErrors: import("svelte/store").Readable<ErrorInfo[]>;
    init: () => void;
    connect: () => void;
    disconnect: () => void;
    selectRequest: (request: RequestInfo | null) => void;
    setActiveTab: (tab: TabType) => void;
    toggleMinimized: () => void;
    setMinimized: (minimized: boolean) => void;
    clearRequests: () => void;
    destroy: () => void;
};
export type DebugBarStore = ReturnType<typeof createDebugBarStore>;
