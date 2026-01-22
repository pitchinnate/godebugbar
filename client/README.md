# Go Debug Bar - Client

A Svelte component for displaying debug information from Go web applications in real-time.

## Installation

```bash
pnpm add go-debug-bar
# or
npm install go-debug-bar
```

## Quick Start

```svelte
<script>
  import { DebugBar } from 'go-debug-bar';
</script>

<DebugBar config={{ wsUrl: 'ws://localhost:8080/_debugbar/ws' }} />
```

## Configuration

```typescript
interface DebugBarConfig {
  // WebSocket server URL (required)
  wsUrl: string;

  // Auto-connect on initialization (default: true)
  autoConnect?: boolean;

  // Auto-reconnect on disconnect (default: true)
  autoReconnect?: boolean;

  // Reconnect delay in milliseconds (default: 3000)
  reconnectDelay?: number;

  // Maximum reconnect attempts, 0 = infinite (default: 0)
  maxReconnectAttempts?: number;

  // Maximum requests to keep in history (default: 100)
  maxRequests?: number;

  // Start minimized (default: false)
  startMinimized?: boolean;

  // Position of the debug bar (default: 'bottom')
  position?: 'bottom' | 'top' | 'left' | 'right';

  // Default height/width in pixels (default: 300)
  defaultSize?: number;
}
```

## Full Example

```svelte
<script>
  import { DebugBar } from 'go-debug-bar';

  const config = {
    wsUrl: 'ws://localhost:8080/_debugbar/ws',
    autoConnect: true,
    autoReconnect: true,
    reconnectDelay: 3000,
    maxRequests: 100,
    startMinimized: false
  };
</script>

<DebugBar {config} />
```

## Using the Store Directly

For more control, you can use the store directly:

```svelte
<script>
  import { createDebugBarStore } from 'go-debug-bar';
  import { onDestroy } from 'svelte';

  const store = createDebugBarStore({
    wsUrl: 'ws://localhost:8080/_debugbar/ws'
  });

  // Access reactive state
  $: status = $store.status;
  $: requests = $store.requests;

  // Actions
  function clearAll() {
    store.clearRequests();
  }

  onDestroy(() => {
    store.destroy();
  });
</script>

<p>Connection: {$status}</p>
<p>Requests: {$requests.length}</p>
<button on:click={clearAll}>Clear</button>
```

## Using the WebSocket Client Directly

For non-Svelte applications or custom implementations:

```typescript
import { DebugBarWebSocket } from 'go-debug-bar';

const ws = new DebugBarWebSocket({
  wsUrl: 'ws://localhost:8080/_debugbar/ws',
  autoConnect: true
});

// Listen for messages
ws.onMessage((message) => {
  console.log('Received:', message.type, message.payload);
});

// Listen for status changes
ws.onStatusChange((status, error) => {
  console.log('Status:', status, error);
});

// Manual control
ws.connect();
ws.disconnect();
ws.destroy();
```

## Features

- **Real-time Updates** - Requests, queries, and errors stream in via WebSocket
- **Request Tracking** - View all HTTP requests with method, path, status, duration
- **Query Tracking** - See all database queries with SQL, args, duration, and source
- **Error Tracking** - View errors with stack traces and context
- **Resizable Panel** - Drag to resize the debug bar
- **Minimizable** - Click to minimize/expand
- **Dark Theme** - Easy on the eyes

## Tabs

| Tab | Description |
|-----|-------------|
| Requests | List of all HTTP requests with details |
| Queries | All database queries across all requests |
| Errors | All errors/warnings across all requests |

## Keyboard Shortcuts

- Click the chevron icon to minimize/expand
- Click trash icon to clear all requests

## TypeScript Support

Full TypeScript support with exported types:

```typescript
import type {
  RequestInfo,
  QueryInfo,
  ErrorInfo,
  DebugBarConfig,
  ConnectionStatus
} from 'go-debug-bar';
```

## Requirements

- Svelte 4.x or 5.x
- Go backend with `godebugbar` package

## License

MIT
