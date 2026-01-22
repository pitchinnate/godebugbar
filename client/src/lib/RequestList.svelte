<script lang="ts">
	import type { RequestInfo } from './types.js';

	export let requests: RequestInfo[];
	export let selectedId: string | undefined;
	export let onSelect: (request: RequestInfo) => void;
	export let getMethodColor: (method: string) => string;

	function formatDuration(ms: number): string {
		if (ms < 1) return '<1ms';
		if (ms < 1000) return `${Math.round(ms)}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatTime(timestamp: string): string {
		const date = new Date(timestamp);
		return date.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	}

	function getStatusColor(status: number): string {
		if (status >= 500) return '#ef4444';
		if (status >= 400) return '#f59e0b';
		if (status >= 300) return '#3b82f6';
		if (status >= 200) return '#22c55e';
		return '#6b7280';
	}

	// Sort requests by start time, newest first
	$: sortedRequests = [...requests].reverse();
</script>

<div class="request-list">
	{#each sortedRequests as request (request.id)}
		<button
			class="request-item"
			class:selected={request.id === selectedId}
			class:has-errors={request.errors?.length > 0}
			on:click={() => onSelect(request)}
		>
			<div class="request-main">
				<span class="method" style="color: {getMethodColor(request.method)}">{request.method}</span>
				<span class="path" title={request.path}>{request.path}</span>
			</div>
			<div class="request-meta">
				<span class="status" style="color: {getStatusColor(request.status_code)}">
					{request.status_code || '...'}
				</span>
				<span class="duration">{formatDuration(request.duration_ms)}</span>
				<span class="queries" title="Database queries">
					<svg width="10" height="10" viewBox="0 0 24 24" fill="currentColor">
						<path d="M12 3C7.58 3 4 4.79 4 7v10c0 2.21 3.58 4 8 4s8-1.79 8-4V7c0-2.21-3.58-4-8-4z"/>
					</svg>
					{request.queries?.length || 0}
				</span>
				<span class="time">{formatTime(request.start_time)}</span>
			</div>
		</button>
	{/each}

	{#if requests.length === 0}
		<div class="empty">No requests captured yet</div>
	{/if}
</div>

<style>
	.request-list {
		display: flex;
		flex-direction: column;
	}

	.request-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 8px 12px;
		background: transparent;
		border: none;
		border-bottom: 1px solid #313244;
		cursor: pointer;
		text-align: left;
		width: 100%;
		transition: background 0.15s;
	}

	.request-item:hover {
		background: #313244;
	}

	.request-item.selected {
		background: #45475a;
	}

	.request-item.has-errors {
		border-left: 2px solid #f38ba8;
	}

	.request-main {
		display: flex;
		align-items: center;
		gap: 8px;
		min-width: 0;
		flex: 1;
	}

	.method {
		font-weight: 600;
		font-size: 11px;
		flex-shrink: 0;
	}

	.path {
		color: #cdd6f4;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.request-meta {
		display: flex;
		align-items: center;
		gap: 12px;
		flex-shrink: 0;
		color: #6c7086;
		font-size: 11px;
	}

	.status {
		font-weight: 600;
	}

	.duration {
		min-width: 45px;
		text-align: right;
	}

	.queries {
		display: flex;
		align-items: center;
		gap: 3px;
	}

	.time {
		min-width: 60px;
		text-align: right;
	}

	.empty {
		padding: 24px;
		text-align: center;
		color: #6c7086;
	}
</style>
