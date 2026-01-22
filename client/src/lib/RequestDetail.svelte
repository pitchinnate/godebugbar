<script lang="ts">
	import type { RequestInfo } from './types.js';

	export let request: RequestInfo;
	export let getMethodColor: (method: string) => string;

	let activeSection: 'overview' | 'headers' | 'queries' | 'errors' | 'body' = 'overview';

	function formatDuration(ms: number): string {
		if (ms < 1) return '<1ms';
		if (ms < 1000) return `${Math.round(ms)}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatBytes(bytes: number): string {
		if (bytes < 1024) return `${bytes} B`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
	}

	function formatTimestamp(timestamp: string): string {
		return new Date(timestamp).toLocaleString();
	}

	function getStatusColor(status: number): string {
		if (status >= 500) return '#ef4444';
		if (status >= 400) return '#f59e0b';
		if (status >= 300) return '#3b82f6';
		if (status >= 200) return '#22c55e';
		return '#6b7280';
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

	$: hasQueries = request.queries?.length > 0;
	$: hasErrors = request.errors?.length > 0;
	$: hasBody = request.request_body && request.request_body.length > 0;
	$: totalQueryTime = request.queries?.reduce((sum, q) => sum + q.duration_ms, 0) || 0;
</script>

<div class="request-detail">
	<div class="detail-header">
		<span class="method" style="color: {getMethodColor(request.method)}">{request.method}</span>
		<span class="path">{request.path}</span>
		<span class="status" style="color: {getStatusColor(request.status_code)}">
			{request.status_code}
		</span>
	</div>

	<div class="section-tabs">
		<button class:active={activeSection === 'overview'} on:click={() => (activeSection = 'overview')}>
			Overview
		</button>
		<button class:active={activeSection === 'headers'} on:click={() => (activeSection = 'headers')}>
			Headers
		</button>
		<button class:active={activeSection === 'queries'} on:click={() => (activeSection = 'queries')} disabled={!hasQueries}>
			Queries ({request.queries?.length || 0})
		</button>
		<button class:active={activeSection === 'errors'} on:click={() => (activeSection = 'errors')} disabled={!hasErrors} class:has-errors={hasErrors}>
			Errors ({request.errors?.length || 0})
		</button>
		{#if hasBody}
			<button class:active={activeSection === 'body'} on:click={() => (activeSection = 'body')}>
				Body
			</button>
		{/if}
	</div>

	<div class="section-content">
		{#if activeSection === 'overview'}
			<div class="overview-grid">
				<div class="stat">
					<span class="stat-label">Duration</span>
					<span class="stat-value">{formatDuration(request.duration_ms)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Response Size</span>
					<span class="stat-value">{formatBytes(request.response_size)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Memory</span>
					<span class="stat-value">{formatBytes(request.memory_usage)}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Client IP</span>
					<span class="stat-value">{request.client_ip}</span>
				</div>
				<div class="stat">
					<span class="stat-label">Queries</span>
					<span class="stat-value">{request.queries?.length || 0} ({formatDuration(totalQueryTime)})</span>
				</div>
				<div class="stat">
					<span class="stat-label">Start Time</span>
					<span class="stat-value">{formatTimestamp(request.start_time)}</span>
				</div>
			</div>

			{#if Object.keys(request.query_params || {}).length > 0}
				<div class="subsection">
					<h4>Query Parameters</h4>
					<div class="key-value-list">
						{#each Object.entries(request.query_params) as [key, value]}
							<div class="key-value">
								<span class="key">{key}</span>
								<span class="value">{value}</span>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			{#if request.custom_data && Object.keys(request.custom_data).length > 0}
				<div class="subsection">
					<h4>Custom Data</h4>
					<pre class="code-block">{JSON.stringify(request.custom_data, null, 2)}</pre>
				</div>
			{/if}
		{:else if activeSection === 'headers'}
			<div class="key-value-list">
				{#each Object.entries(request.headers || {}) as [key, value]}
					<div class="key-value">
						<span class="key">{key}</span>
						<span class="value">{value}</span>
					</div>
				{/each}
			</div>
		{:else if activeSection === 'queries'}
			<div class="query-list">
				{#each request.queries || [] as query, i}
					<div class="query-item">
						<div class="query-header">
							<span class="query-index">#{i + 1}</span>
							<span class="query-duration">{formatDuration(query.duration_ms)}</span>
							{#if query.rows_affected > 0}
								<span class="query-rows">{query.rows_affected} rows</span>
							{/if}
							{#if query.error}
								<span class="query-error">Error</span>
							{/if}
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
							<div class="query-error-msg">{query.error}</div>
						{/if}
					</div>
				{/each}
			</div>
		{:else if activeSection === 'errors'}
			<div class="error-list">
				{#each request.errors || [] as error}
					<div class="error-item">
						<div class="error-header">
							<span class="error-type" style="color: {getErrorTypeColor(error.type)}">{error.type}</span>
							<span class="error-time">{formatTimestamp(error.timestamp)}</span>
						</div>
						<div class="error-message">{error.message}</div>
						{#if error.context}
							<pre class="error-context">{JSON.stringify(error.context, null, 2)}</pre>
						{/if}
						{#if error.stack}
							<details class="error-stack">
								<summary>Stack Trace</summary>
								<pre>{error.stack}</pre>
							</details>
						{/if}
					</div>
				{/each}
			</div>
		{:else if activeSection === 'body'}
			<pre class="code-block">{request.request_body}</pre>
		{/if}
	</div>
</div>

<style>
	.request-detail {
		height: 100%;
		display: flex;
		flex-direction: column;
	}

	.detail-header {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px;
		background: #181825;
		border-bottom: 1px solid #313244;
	}

	.method {
		font-weight: 600;
	}

	.path {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.status {
		font-weight: 600;
	}

	.section-tabs {
		display: flex;
		gap: 2px;
		padding: 8px 12px;
		background: #181825;
		border-bottom: 1px solid #313244;
	}

	.section-tabs button {
		padding: 4px 10px;
		background: transparent;
		border: none;
		border-radius: 4px;
		color: #6c7086;
		cursor: pointer;
		font-size: 11px;
		transition: all 0.15s;
	}

	.section-tabs button:hover:not(:disabled) {
		background: #313244;
		color: #cdd6f4;
	}

	.section-tabs button.active {
		background: #45475a;
		color: #cdd6f4;
	}

	.section-tabs button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.section-tabs button.has-errors {
		color: #f38ba8;
	}

	.section-content {
		flex: 1;
		overflow: auto;
		padding: 12px;
	}

	.overview-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
		gap: 12px;
	}

	.stat {
		background: #181825;
		padding: 10px;
		border-radius: 4px;
	}

	.stat-label {
		display: block;
		color: #6c7086;
		font-size: 10px;
		text-transform: uppercase;
		margin-bottom: 4px;
	}

	.stat-value {
		font-weight: 500;
	}

	.subsection {
		margin-top: 16px;
	}

	.subsection h4 {
		margin: 0 0 8px 0;
		color: #6c7086;
		font-size: 11px;
		text-transform: uppercase;
	}

	.key-value-list {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.key-value {
		display: flex;
		gap: 12px;
		padding: 6px 8px;
		background: #181825;
		border-radius: 4px;
	}

	.key {
		color: #89b4fa;
		min-width: 120px;
		flex-shrink: 0;
	}

	.value {
		color: #a6e3a1;
		word-break: break-all;
	}

	.code-block {
		background: #181825;
		padding: 12px;
		border-radius: 4px;
		overflow: auto;
		margin: 0;
		font-size: 11px;
		white-space: pre-wrap;
		word-break: break-all;
	}

	.query-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.query-item {
		background: #181825;
		border-radius: 4px;
		padding: 10px;
	}

	.query-header {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-bottom: 8px;
		font-size: 11px;
	}

	.query-index {
		color: #6c7086;
	}

	.query-duration {
		color: #f9e2af;
	}

	.query-rows {
		color: #89b4fa;
	}

	.query-error {
		color: #f38ba8;
		font-weight: 600;
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

	.query-error-msg {
		margin-top: 6px;
		color: #f38ba8;
		font-size: 11px;
	}

	.error-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
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
		margin-bottom: 6px;
		font-size: 11px;
	}

	.error-type {
		font-weight: 600;
		text-transform: uppercase;
	}

	.error-time {
		color: #6c7086;
	}

	.error-message {
		margin-bottom: 8px;
	}

	.error-context {
		background: #11111b;
		padding: 8px;
		border-radius: 4px;
		margin: 8px 0;
		font-size: 11px;
	}

	.error-stack {
		margin-top: 8px;
	}

	.error-stack summary {
		cursor: pointer;
		color: #6c7086;
		font-size: 11px;
	}

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
</style>
