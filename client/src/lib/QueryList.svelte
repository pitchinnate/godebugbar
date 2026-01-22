<script lang="ts">
	import type { RequestInfo, QueryInfo } from './types.js';

	export let requests: RequestInfo[];

	interface QueryWithRequest extends QueryInfo {
		request_path: string;
		request_method: string;
	}

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

	// Flatten all queries from all requests, newest first
	$: allQueries = requests
		.flatMap((req) =>
			(req.queries || []).map((q) => ({
				...q,
				request_path: req.path,
				request_method: req.method
			}))
		)
		.reverse() as QueryWithRequest[];

	$: totalTime = allQueries.reduce((sum, q) => sum + q.duration_ms, 0);
</script>

<div class="query-list-container">
	<div class="summary-bar">
		<span>{allQueries.length} queries</span>
		<span class="separator">|</span>
		<span>Total: {formatDuration(totalTime)}</span>
	</div>

	<div class="query-items">
		{#each allQueries as query (query.id)}
			<div class="query-item" class:has-error={!!query.error}>
				<div class="query-header">
					<span class="query-request">{query.request_method} {query.request_path}</span>
					<span class="query-meta">
						<span class="query-duration">{formatDuration(query.duration_ms)}</span>
						{#if query.rows_affected > 0}
							<span class="query-rows">{query.rows_affected} rows</span>
						{/if}
						<span class="query-time">{formatTime(query.start_time)}</span>
					</span>
				</div>
				<pre class="query-sql">{query.query}</pre>
				{#if query.args && query.args.length > 0}
					<div class="query-args">
						<span class="args-label">Args:</span>
						<code>{JSON.stringify(query.args)}</code>
					</div>
				{/if}
				{#if query.source}
					<div class="query-source">{query.source}</div>
				{/if}
				{#if query.error}
					<div class="query-error">{query.error}</div>
				{/if}
			</div>
		{/each}

		{#if allQueries.length === 0}
			<div class="empty">No queries captured yet</div>
		{/if}
	</div>
</div>

<style>
	.query-list-container {
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.summary-bar {
		padding: 8px 12px;
		background: #181825;
		border-bottom: 1px solid #313244;
		color: #6c7086;
		font-size: 11px;
	}

	.separator {
		margin: 0 8px;
	}

	.query-items {
		flex: 1;
		overflow: auto;
		padding: 12px;
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.query-item {
		background: #181825;
		border-radius: 4px;
		padding: 10px;
	}

	.query-item.has-error {
		border-left: 3px solid #f38ba8;
	}

	.query-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 8px;
		font-size: 11px;
	}

	.query-request {
		color: #89b4fa;
	}

	.query-meta {
		display: flex;
		gap: 12px;
		color: #6c7086;
	}

	.query-duration {
		color: #f9e2af;
	}

	.query-rows {
		color: #89b4fa;
	}

	.query-sql {
		background: #11111b;
		padding: 8px;
		border-radius: 4px;
		margin: 0;
		font-size: 11px;
		overflow-x: auto;
		white-space: pre-wrap;
		word-break: break-all;
		color: #a6e3a1;
	}

	.query-args {
		margin-top: 6px;
		font-size: 11px;
		color: #6c7086;
	}

	.query-args code {
		color: #fab387;
	}

	.query-source {
		margin-top: 6px;
		font-size: 10px;
		color: #6c7086;
	}

	.query-error {
		margin-top: 6px;
		color: #f38ba8;
		font-size: 11px;
	}

	.empty {
		padding: 24px;
		text-align: center;
		color: #6c7086;
	}
</style>
