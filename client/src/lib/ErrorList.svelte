<script lang="ts">
	import type { RequestInfo, ErrorInfo } from './types.js';

	export let requests: RequestInfo[];

	interface ErrorWithRequest extends ErrorInfo {
		request_path: string;
		request_method: string;
	}

	function formatTimestamp(timestamp: string): string {
		return new Date(timestamp).toLocaleString();
	}

	function getErrorTypeColor(type: string): string {
		switch (type) {
			case 'exception':
				return '#ef4444';
			case 'warning':
				return '#f59e0b';
			case 'notice':
				return '#3b82f6';
			case 'debug':
				return '#6b7280';
			default:
				return '#6b7280';
		}
	}

	// Flatten all errors from all requests, newest first
	$: allErrors = requests
		.flatMap((req) =>
			(req.errors || []).map((e) => ({
				...e,
				request_path: req.path,
				request_method: req.method
			}))
		)
		.reverse() as ErrorWithRequest[];

	$: errorCounts = allErrors.reduce(
		(acc, e) => {
			acc[e.type] = (acc[e.type] || 0) + 1;
			return acc;
		},
		{} as Record<string, number>
	);
</script>

<div class="error-list-container">
	<div class="summary-bar">
		<span>{allErrors.length} total</span>
		{#each Object.entries(errorCounts) as [type, count]}
			<span class="error-count" style="color: {getErrorTypeColor(type)}">{type}: {count}</span>
		{/each}
	</div>

	<div class="error-items">
		{#each allErrors as error (error.id)}
			<div class="error-item">
				<div class="error-header">
					<div class="error-left">
						<span class="error-type" style="color: {getErrorTypeColor(error.type)}">{error.type}</span>
						<span class="error-request">{error.request_method} {error.request_path}</span>
					</div>
					<span class="error-time">{formatTimestamp(error.timestamp)}</span>
				</div>
				<div class="error-message">{error.message}</div>
				{#if error.context && Object.keys(error.context).length > 0}
					<details class="error-context">
						<summary>Context</summary>
						<pre>{JSON.stringify(error.context, null, 2)}</pre>
					</details>
				{/if}
				{#if error.stack}
					<details class="error-stack">
						<summary>Stack Trace</summary>
						<pre>{error.stack}</pre>
					</details>
				{/if}
			</div>
		{/each}

		{#if allErrors.length === 0}
			<div class="empty">No errors captured</div>
		{/if}
	</div>
</div>

<style>
	.error-list-container {
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.summary-bar {
		display: flex;
		gap: 16px;
		padding: 8px 12px;
		background: #181825;
		border-bottom: 1px solid #313244;
		color: #6c7086;
		font-size: 11px;
	}

	.error-count {
		font-weight: 500;
	}

	.error-items {
		flex: 1;
		overflow: auto;
		padding: 12px;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.error-item {
		background: #181825;
		border-radius: 4px;
		padding: 10px;
		border-left: 3px solid #f38ba8;
	}

	.error-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 8px;
		font-size: 11px;
	}

	.error-left {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.error-type {
		font-weight: 600;
		text-transform: uppercase;
	}

	.error-request {
		color: #89b4fa;
	}

	.error-time {
		color: #6c7086;
	}

	.error-message {
		margin-bottom: 8px;
		word-break: break-word;
	}

	.error-context,
	.error-stack {
		margin-top: 8px;
	}

	.error-context summary,
	.error-stack summary {
		cursor: pointer;
		color: #6c7086;
		font-size: 11px;
		padding: 4px 0;
	}

	.error-context pre,
	.error-stack pre {
		background: #11111b;
		padding: 8px;
		border-radius: 4px;
		margin-top: 6px;
		font-size: 10px;
		overflow-x: auto;
		white-space: pre-wrap;
		word-break: break-all;
	}

	.empty {
		padding: 24px;
		text-align: center;
		color: #6c7086;
	}
</style>
