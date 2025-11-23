<script lang="ts">
	import { onMount } from 'svelte';

	let apiHealth = $state<{ status: string; service: string; version: string } | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			const response = await fetch('/api/health');
			if (!response.ok) {
				throw new Error('Failed to fetch API health');
			}
			apiHealth = await response.json();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			loading = false;
		}
	});
</script>

<div class="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-primary-50 to-secondary-50">
	<div class="card max-w-2xl text-center">
		<h1 class="mb-4 text-5xl font-bold text-primary-700">
			Ham-Radio Cloud
		</h1>
		<p class="mb-8 text-xl text-gray-600">
			Modern cloud platform for amateur radio operators
		</p>

		<div class="mb-8 grid grid-cols-1 gap-4 md:grid-cols-3">
			<div class="card border-l-4 border-primary-500">
				<div class="text-3xl font-bold text-primary-600">500+</div>
				<div class="text-sm text-gray-600">Free QSOs</div>
			</div>
			<div class="card border-l-4 border-secondary-500">
				<div class="text-3xl font-bold text-secondary-600">Real-time</div>
				<div class="text-sm text-gray-600">Propagation Data</div>
			</div>
			<div class="card border-l-4 border-primary-500">
				<div class="text-3xl font-bold text-primary-600">LoTW</div>
				<div class="text-sm text-gray-600">Auto-Sync</div>
			</div>
		</div>

		{#if loading}
			<div class="text-gray-500">Checking API connection...</div>
		{:else if error}
			<div class="rounded-lg bg-red-50 p-4 text-red-600">
				<strong>API Error:</strong> {error}
				<div class="mt-2 text-sm">Make sure the backend is running on port 8080</div>
			</div>
		{:else if apiHealth}
			<div class="rounded-lg bg-green-50 p-4 text-green-700">
				<strong>✓ API Connected</strong>
				<div class="mt-1 text-sm text-green-600">
					{apiHealth.service} v{apiHealth.version} - Status: {apiHealth.status}
				</div>
			</div>
		{/if}

		<div class="mt-8 flex gap-4 justify-center">
			<a href="/login" class="btn btn-primary">Sign In</a>
			<a href="/register" class="btn btn-secondary">Get Started</a>
		</div>

		<div class="mt-8 text-sm text-gray-500">
			<p>Phase 1: Project Foundation & Core Infrastructure</p>
			<p class="mt-1">MVP Target: 60-90 days</p>
		</div>
	</div>
</div>
