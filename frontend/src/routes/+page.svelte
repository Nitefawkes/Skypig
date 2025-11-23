<script lang="ts">
	import { onMount } from 'svelte';
	import PropagationWidget from '$components/PropagationWidget.svelte';
	import Navigation from '$components/Navigation.svelte';

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

<Navigation />

<div class="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 p-4">
	<div class="mx-auto max-w-7xl">
		<div class="card mb-8 max-w-2xl text-center mx-auto">
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

		<div class="mt-8 flex flex-col sm:flex-row gap-4 justify-center">
			<a href="/logbook" class="btn btn-primary">
				📖 Try the Logbook
			</a>
			<a href="/login" class="btn btn-secondary">Sign In</a>
			<a href="/register" class="btn btn-secondary">Get Started</a>
		</div>

		<div class="mt-8 text-sm text-gray-500">
			<p>Phase 2: Core Features Complete</p>
			<p class="mt-1">✅ QSO Logging | ✅ ADIF Import/Export | ✅ Propagation Data</p>
		</div>
	</div>

	<!-- Propagation Widget -->
	<div id="propagation" class="mx-auto max-w-7xl mt-8">
		<PropagationWidget />
	</div>
</div>
</div>
