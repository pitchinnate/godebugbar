<script>import { onMount, onDestroy } from "svelte";
import { createDebugBarStore } from "./store.js";
import RequestList from "./RequestList.svelte";
import RequestDetail from "./RequestDetail.svelte";
import QueryList from "./QueryList.svelte";
import ErrorList from "./ErrorList.svelte";
export let config;
let store;
let containerHeight = 300;
let isResizing = false;
$: status = store?.status;
$: requests = store?.requests;
$: selectedRequest = store?.selectedRequest;
$: activeTab = store?.activeTab;
$: minimized = store?.minimized;
$: totalQueries = store?.totalQueries;
$: totalErrors = store?.totalErrors;
onMount(() => {
  store = createDebugBarStore(config);
});
onDestroy(() => {
  store?.destroy();
});
function handleTabClick(tab) {
  store?.setActiveTab(tab);
}
function handleRequestSelect(request) {
  store?.selectRequest(request);
}
function handleToggleMinimize() {
  store?.toggleMinimized();
}
function handleClear() {
  store?.clearRequests();
}
function startResize(e) {
  isResizing = true;
  const startY = e.clientY;
  const startHeight = containerHeight;
  function onMouseMove(e2) {
    if (!isResizing) return;
    const delta = startY - e2.clientY;
    containerHeight = Math.max(100, Math.min(600, startHeight + delta));
  }
  function onMouseUp() {
    isResizing = false;
    window.removeEventListener("mousemove", onMouseMove);
    window.removeEventListener("mouseup", onMouseUp);
  }
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseup", onMouseUp);
}
function getStatusColor(status2) {
  switch (status2) {
    case "connected":
      return "#22c55e";
    case "connecting":
      return "#eab308";
    case "error":
      return "#ef4444";
    default:
      return "#6b7280";
  }
}
function getMethodColor(method) {
  switch (method) {
    case "GET":
      return "#22c55e";
    case "POST":
      return "#3b82f6";
    case "PUT":
      return "#f59e0b";
    case "PATCH":
      return "#8b5cf6";
    case "DELETE":
      return "#ef4444";
    default:
      return "#6b7280";
  }
}
</script>

{#if store}
	<div class="debugbar" class:minimized={$minimized}>
		<!-- Resize handle -->
		{#if !$minimized}
			<div class="resize-handle" on:mousedown={startResize} role="separator" aria-orientation="horizontal"></div>
		{/if}

		<!-- Header bar -->
		<div class="header">
			<div class="header-left">
				<span class="logo">Debug Bar</span>
				<span class="status-indicator" style="background-color: {getStatusColor($status)}"></span>
				<span class="status-text">{$status}</span>
			</div>

			<div class="header-center">
				{#if !$minimized}
					<button
						class="tab"
						class:active={$activeTab === 'requests'}
						on:click={() => handleTabClick('requests')}
					>
						Requests ({$requests?.length || 0})
					</button>
					<button
						class="tab"
						class:active={$activeTab === 'queries'}
						on:click={() => handleTabClick('queries')}
					>
						Queries ({$totalQueries || 0})
					</button>
					<button
						class="tab"
						class:active={$activeTab === 'errors'}
						on:click={() => handleTabClick('errors')}
						class:has-errors={($totalErrors || 0) > 0}
					>
						Errors ({$totalErrors || 0})
					</button>
				{/if}
			</div>

			<div class="header-right">
				{#if !$minimized}
					<button class="action-btn" on:click={handleClear} title="Clear requests">
						<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" />
						</svg>
					</button>
				{/if}
				<button class="action-btn" on:click={handleToggleMinimize} title={$minimized ? 'Expand' : 'Minimize'}>
					{#if $minimized}
						<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M18 15l-6-6-6 6" />
						</svg>
					{:else}
						<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M6 9l6 6 6-6" />
						</svg>
					{/if}
				</button>
			</div>
		</div>

		<!-- Content area -->
		{#if !$minimized}
			<div class="content" style="height: {containerHeight}px">
				<div class="panel-left">
					{#if $activeTab === 'requests'}
						<RequestList
							requests={$requests || []}
							selectedId={$selectedRequest?.id}
							onSelect={handleRequestSelect}
							{getMethodColor}
						/>
					{:else if $activeTab === 'queries'}
						<QueryList requests={$requests || []} />
					{:else if $activeTab === 'errors'}
						<ErrorList requests={$requests || []} />
					{/if}
				</div>

				{#if $selectedRequest && $activeTab === 'requests'}
					<div class="panel-right">
						<RequestDetail request={$selectedRequest} {getMethodColor} />
					</div>
				{/if}
			</div>
		{/if}
	</div>
{/if}

<style>
	.debugbar {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		background: #1e1e2e;
		color: #cdd6f4;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
		font-size: 12px;
		z-index: 99999;
		box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.3);
	}

	.debugbar.minimized {
		height: auto;
	}

	.resize-handle {
		height: 4px;
		background: #313244;
		cursor: ns-resize;
		transition: background 0.2s;
	}

	.resize-handle:hover {
		background: #45475a;
	}

	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 6px 12px;
		background: #181825;
		border-bottom: 1px solid #313244;
	}

	.header-left,
	.header-right {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.header-center {
		display: flex;
		gap: 4px;
	}

	.logo {
		font-weight: 600;
		color: #cba6f7;
	}

	.status-indicator {
		width: 8px;
		height: 8px;
		border-radius: 50%;
	}

	.status-text {
		color: #6c7086;
		font-size: 11px;
	}

	.tab {
		padding: 4px 12px;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #6c7086;
		cursor: pointer;
		font-size: 12px;
		transition: all 0.2s;
	}

	.tab:hover {
		background: #313244;
		color: #cdd6f4;
	}

	.tab.active {
		background: #45475a;
		color: #cdd6f4;
	}

	.tab.has-errors {
		color: #f38ba8;
	}

	.action-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #6c7086;
		cursor: pointer;
		transition: all 0.2s;
	}

	.action-btn:hover {
		background: #313244;
		color: #cdd6f4;
	}

	.content {
		display: flex;
		overflow: hidden;
	}

	.panel-left {
		flex: 1;
		overflow: auto;
		border-right: 1px solid #313244;
	}

	.panel-right {
		flex: 1;
		overflow: auto;
	}
</style>
