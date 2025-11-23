<script lang="ts">
	import { onMount } from 'svelte';

	let solarData = $state({
		solar_flux: 0,
		sunspot_number: 0,
		a_index: 0,
		k_index: 0,
		xray_flux: 'N/A'
	});

	let bandConditions = $state([
		{ band: '160m', condition: 'poor', score: 2 },
		{ band: '80m', condition: 'fair', score: 5 },
		{ band: '40m', condition: 'good', score: 8 },
		{ band: '20m', condition: 'good', score: 9 },
		{ band: '15m', condition: 'fair', score: 6 },
		{ band: '10m', condition: 'poor', score: 3 },
		{ band: '6m', condition: 'poor', score: 2 }
	]);

	function getConditionColor(condition: string) {
		switch (condition) {
			case 'good':
				return 'text-green-500';
			case 'fair':
				return 'text-yellow-500';
			case 'poor':
				return 'text-red-500';
			default:
				return 'text-slate-500';
		}
	}
</script>

<svelte:head>
	<title>Propagation - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<div class="mb-8">
		<h1 class="text-3xl font-bold">Propagation Conditions</h1>
		<p class="mt-2 text-slate-400">Real-time solar data and band conditions</p>
	</div>

	<!-- Solar Data -->
	<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-3 lg:grid-cols-5">
		<div class="card">
			<div class="text-sm text-slate-400">Solar Flux</div>
			<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.solar_flux}</div>
			<div class="mt-1 text-xs text-slate-500">SFU</div>
		</div>
		<div class="card">
			<div class="text-sm text-slate-400">Sunspot Number</div>
			<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.sunspot_number}</div>
		</div>
		<div class="card">
			<div class="text-sm text-slate-400">A-Index</div>
			<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.a_index}</div>
		</div>
		<div class="card">
			<div class="text-sm text-slate-400">K-Index</div>
			<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.k_index}</div>
		</div>
		<div class="card">
			<div class="text-sm text-slate-400">X-Ray Flux</div>
			<div class="mt-2 text-2xl font-bold text-primary-500">{solarData.xray_flux}</div>
		</div>
	</div>

	<!-- Band Conditions -->
	<div class="card">
		<h2 class="mb-6 text-2xl font-bold">Current Band Conditions</h2>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
			{#each bandConditions as band}
				<div class="rounded-lg border border-slate-700 bg-slate-800/50 p-4">
					<div class="flex items-center justify-between">
						<div class="text-xl font-bold">{band.band}</div>
						<div class="text-sm font-semibold uppercase {getConditionColor(band.condition)}">
							{band.condition}
						</div>
					</div>
					<div class="mt-3 h-2 overflow-hidden rounded-full bg-slate-700">
						<div
							class="h-full bg-gradient-to-r from-primary-600 to-primary-400"
							style="width: {(band.score / 10) * 100}%"
						></div>
					</div>
					<div class="mt-2 text-xs text-slate-500">Score: {band.score}/10</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Info Panel -->
	<div class="card mt-8 bg-primary-900/10 border-primary-700">
		<h3 class="mb-4 text-xl font-semibold">About Propagation Data</h3>
		<div class="space-y-3 text-slate-400">
			<p>
				<strong class="text-slate-300">Solar Flux:</strong> Higher values (>150) generally indicate
				better HF propagation conditions.
			</p>
			<p>
				<strong class="text-slate-300">K-Index:</strong> Lower is better. Values above 4 indicate
				disturbed conditions.
			</p>
			<p>
				<strong class="text-slate-300">A-Index:</strong> Lower is better. Values below 10 are
				ideal for propagation.
			</p>
		</div>
	</div>
</div>
