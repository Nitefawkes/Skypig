<script lang="ts">
	import { onMount } from 'svelte';

	interface PropagationData {
		solar_flux: number;
		k_index: number;
		a_index: number;
		sunspot_number: number;
		updated_at: string;
	}

	interface BandCondition {
		band: string;
		day: string;
		night: string;
		reasoning: string;
	}

	interface PropagationForecast {
		current: PropagationData;
		band_conditions: BandCondition[];
		summary: string;
		last_updated: string;
	}

	let forecast = $state<PropagationForecast | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const conditionColors: Record<string, string> = {
		excellent: 'bg-green-500',
		good: 'bg-blue-500',
		fair: 'bg-yellow-500',
		poor: 'bg-red-500'
	};

	const conditionText: Record<string, string> = {
		excellent: 'Excellent',
		good: 'Good',
		fair: 'Fair',
		poor: 'Poor'
	};

	onMount(async () => {
		await fetchPropagation();
		// Refresh every 15 minutes
		const interval = setInterval(fetchPropagation, 15 * 60 * 1000);
		return () => clearInterval(interval);
	});

	async function fetchPropagation() {
		try {
			loading = true;
			const response = await fetch('/api/propagation/forecast');
			if (!response.ok) {
				throw new Error('Failed to fetch propagation data');
			}
			const result = await response.json();
			forecast = result.data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Unknown error';
			console.error('Propagation fetch error:', e);
		} finally {
			loading = false;
		}
	}

	function getSolarActivity(sfi: number): string {
		if (sfi < 70) return 'Very Low';
		if (sfi < 100) return 'Low';
		if (sfi < 150) return 'Moderate';
		if (sfi < 200) return 'High';
		return 'Very High';
	}

	function getGeomagActivity(k: number): string {
		if (k <= 2) return 'Quiet';
		if (k <= 4) return 'Unsettled';
		if (k <= 6) return 'Active';
		if (k <= 8) return 'Storm';
		return 'Severe Storm';
	}
</script>

<div class="card">
	<div class="mb-4 flex items-center justify-between">
		<h2 class="text-xl font-bold text-gray-900">Propagation Conditions</h2>
		{#if forecast}
			<span class="text-xs text-gray-500">
				Updated: {new Date(forecast.last_updated).toLocaleTimeString()}
			</span>
		{/if}
	</div>

	{#if loading && !forecast}
		<div class="flex items-center justify-center py-8">
			<div class="text-gray-500">Loading propagation data...</div>
		</div>
	{:else if error}
		<div class="rounded-lg bg-red-50 p-4 text-red-600">
			<strong>Error:</strong>
			{error}
		</div>
	{:else if forecast}
		<!-- Summary -->
		<div class="mb-6 rounded-lg bg-blue-50 p-4">
			<p class="text-sm text-blue-900">{forecast.summary}</p>
		</div>

		<!-- Solar Data -->
		<div class="mb-6 grid grid-cols-2 gap-4 md:grid-cols-4">
			<div class="rounded-lg border border-gray-200 p-3">
				<div class="text-xs font-medium text-gray-500">Solar Flux</div>
				<div class="text-2xl font-bold text-gray-900">{forecast.current.solar_flux.toFixed(0)}</div>
				<div class="text-xs text-gray-600">{getSolarActivity(forecast.current.solar_flux)}</div>
			</div>

			<div class="rounded-lg border border-gray-200 p-3">
				<div class="text-xs font-medium text-gray-500">K-Index</div>
				<div class="text-2xl font-bold text-gray-900">{forecast.current.k_index}</div>
				<div class="text-xs text-gray-600">{getGeomagActivity(forecast.current.k_index)}</div>
			</div>

			<div class="rounded-lg border border-gray-200 p-3">
				<div class="text-xs font-medium text-gray-500">A-Index</div>
				<div class="text-2xl font-bold text-gray-900">{forecast.current.a_index}</div>
			</div>

			<div class="rounded-lg border border-gray-200 p-3">
				<div class="text-xs font-medium text-gray-500">Sunspots</div>
				<div class="text-2xl font-bold text-gray-900">{forecast.current.sunspot_number}</div>
			</div>
		</div>

		<!-- Band Conditions -->
		<div>
			<h3 class="mb-3 text-sm font-semibold text-gray-700">Band Conditions</h3>
			<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
				{#each forecast.band_conditions as band}
					<div class="rounded-lg border border-gray-200 p-3">
						<div class="mb-2 text-center font-bold text-gray-900">{band.band}</div>
						<div class="space-y-1">
							<div class="flex items-center justify-between text-xs">
								<span class="text-gray-600">Day:</span>
								<span
									class="rounded px-2 py-0.5 text-white {conditionColors[band.day] ||
										'bg-gray-400'}"
								>
									{conditionText[band.day] || 'N/A'}
								</span>
							</div>
							<div class="flex items-center justify-between text-xs">
								<span class="text-gray-600">Night:</span>
								<span
									class="rounded px-2 py-0.5 text-white {conditionColors[band.night] ||
										'bg-gray-400'}"
								>
									{conditionText[band.night] || 'N/A'}
								</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Legend -->
		<div class="mt-4 flex flex-wrap gap-2 text-xs">
			<span class="text-gray-500">Conditions:</span>
			<span class="flex items-center gap-1">
				<span class="h-3 w-3 rounded bg-green-500"></span>
				Excellent
			</span>
			<span class="flex items-center gap-1">
				<span class="h-3 w-3 rounded bg-blue-500"></span>
				Good
			</span>
			<span class="flex items-center gap-1">
				<span class="h-3 w-3 rounded bg-yellow-500"></span>
				Fair
			</span>
			<span class="flex items-center gap-1">
				<span class="h-3 w-3 rounded bg-red-500"></span>
				Poor
			</span>
		</div>
	{/if}
</div>
