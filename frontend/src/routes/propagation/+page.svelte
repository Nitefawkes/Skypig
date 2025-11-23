<script lang="ts">
	import { onMount } from 'svelte';

	let solarData = $state({
		solar_flux: 0,
		sunspot_number: 0,
		a_index: 0,
		k_index: 0,
		xray_flux: 'N/A',
		timestamp: '',
		source: ''
	});

	let bandConditions = $state<any[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let lastUpdate = $state<string>('');

	onMount(async () => {
		await loadPropagationData();
		// Refresh every 15 minutes
		setInterval(loadPropagationData, 15 * 60 * 1000);
	});

	async function loadPropagationData() {
		loading = true;
		error = null;
		try {
			// Fetch current solar data
			const currentResponse = await fetch('/api/v1/propagation/current');
			if (!currentResponse.ok) {
				throw new Error('Failed to fetch propagation data');
			}
			const currentData = await currentResponse.json();

			solarData = {
				solar_flux: currentData.solar_flux || 0,
				sunspot_number: currentData.sunspot_number || 0,
				a_index: currentData.a_index || 0,
				k_index: currentData.k_index || 0,
				xray_flux: currentData.xray_flux || 'N/A',
				timestamp: currentData.timestamp,
				source: currentData.source
			};

			lastUpdate = new Date(currentData.timestamp).toLocaleString();

			// Fetch band conditions
			const bandsResponse = await fetch('/api/v1/propagation/bands');
			if (!bandsResponse.ok) {
				throw new Error('Failed to fetch band conditions');
			}
			const bandsData = await bandsResponse.json();
			bandConditions = bandsData.conditions || [];
		} catch (err: any) {
			error = err.message || 'Failed to load propagation data';
		} finally {
			loading = false;
		}
	}

	async function handleRefresh() {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/v1/propagation/refresh', {
				method: 'POST'
			});

			if (!response.ok) {
				throw new Error('Failed to refresh data');
			}

			// Reload after refresh
			await loadPropagationData();
		} catch (err: any) {
			error = err.message || 'Failed to refresh data';
			loading = false;
		}
	}

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

	function getConditionBgColor(condition: string) {
		switch (condition) {
			case 'good':
				return 'bg-gradient-to-r from-green-600 to-green-400';
			case 'fair':
				return 'bg-gradient-to-r from-yellow-600 to-yellow-400';
			case 'poor':
				return 'bg-gradient-to-r from-red-600 to-red-400';
			default:
				return 'bg-gradient-to-r from-slate-600 to-slate-400';
		}
	}

	function getSolarFluxStatus(flux: number) {
		if (flux > 150) return { text: 'Excellent', color: 'text-green-500' };
		if (flux > 100) return { text: 'Good', color: 'text-yellow-500' };
		if (flux > 70) return { text: 'Fair', color: 'text-orange-500' };
		return { text: 'Poor', color: 'text-red-500' };
	}

	function getKIndexStatus(k: number) {
		if (k <= 2) return { text: 'Quiet', color: 'text-green-500' };
		if (k <= 3) return { text: 'Unsettled', color: 'text-yellow-500' };
		if (k <= 4) return { text: 'Active', color: 'text-orange-500' };
		return { text: 'Stormy', color: 'text-red-500' };
	}
</script>

<svelte:head>
	<title>Propagation - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<!-- Header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Propagation Conditions</h1>
			<p class="mt-2 text-slate-400">
				Real-time solar data from {solarData.source || 'NOAA SWPC'}
			</p>
			{#if lastUpdate}
				<p class="mt-1 text-sm text-slate-500">Last updated: {lastUpdate}</p>
			{/if}
		</div>
		<button class="btn-primary" onclick={handleRefresh} disabled={loading}>
			{loading ? 'Refreshing...' : 'Refresh Data'}
		</button>
	</div>

	<!-- Error Message -->
	{#if error}
		<div class="mb-6 rounded-lg border border-red-700 bg-red-900/20 p-4 text-red-400">
			{error}
		</div>
	{/if}

	<!-- Solar Data -->
	{#if loading && !solarData.solar_flux}
		<div class="card mb-8 py-12 text-center">
			<div class="text-slate-400">Loading propagation data...</div>
		</div>
	{:else}
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-3 lg:grid-cols-5">
			<div class="card">
				<div class="text-sm text-slate-400">Solar Flux</div>
				<div class="mt-2 text-3xl font-bold text-primary-500">
					{solarData.solar_flux.toFixed(0)}
				</div>
				<div class="mt-1 text-xs text-slate-500">SFU</div>
				<div class="mt-2 text-sm font-semibold {getSolarFluxStatus(solarData.solar_flux).color}">
					{getSolarFluxStatus(solarData.solar_flux).text}
				</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">Sunspot Number</div>
				<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.sunspot_number}</div>
				<div class="mt-1 text-xs text-slate-500">SSN</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">A-Index</div>
				<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.a_index}</div>
				<div class="mt-1 text-xs text-slate-500">Geomagnetic Activity</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">K-Index</div>
				<div class="mt-2 text-3xl font-bold text-primary-500">{solarData.k_index}</div>
				<div class="mt-2 text-sm font-semibold {getKIndexStatus(solarData.k_index).color}">
					{getKIndexStatus(solarData.k_index).text}
				</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">X-Ray Flux</div>
				<div class="mt-2 text-2xl font-bold text-primary-500">{solarData.xray_flux}</div>
				<div class="mt-1 text-xs text-slate-500">Solar Flares</div>
			</div>
		</div>

		<!-- Band Conditions -->
		{#if bandConditions.length > 0}
			<div class="card mb-8">
				<div class="mb-6 flex items-center justify-between">
					<h2 class="text-2xl font-bold">Current Band Conditions</h2>
					<div class="flex gap-4 text-sm">
						<span class="text-slate-400">
							<span class="text-green-500">●</span> Good
						</span>
						<span class="text-slate-400">
							<span class="text-yellow-500">●</span> Fair
						</span>
						<span class="text-slate-400">
							<span class="text-red-500">●</span> Poor
						</span>
					</div>
				</div>
				<div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
					{#each bandConditions as band}
						<div class="rounded-lg border border-slate-700 bg-slate-800/50 p-4">
							<div class="flex items-center justify-between">
								<div class="text-xl font-bold">{band.band}</div>
								<div class="text-xs font-semibold uppercase {getConditionColor(band.condition)}">
									{band.condition}
								</div>
							</div>
							<div class="mt-3 h-2 overflow-hidden rounded-full bg-slate-700">
								<div
									class="h-full {getConditionBgColor(band.condition)}"
									style="width: {(band.score / 10) * 100}%"
								></div>
							</div>
							<div class="mt-2 flex items-center justify-between text-xs text-slate-500">
								<span>Score: {band.score.toFixed(1)}/10</span>
								<span class="capitalize">{band.day_night}</span>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Info Panel -->
		<div class="card border-primary-700 bg-primary-900/10">
			<h3 class="mb-4 text-xl font-semibold">Understanding Propagation Data</h3>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div class="space-y-3 text-sm text-slate-400">
					<p>
						<strong class="text-slate-300">Solar Flux (SFI):</strong> Measured in Solar Flux Units.
						Higher values (>150) indicate better HF propagation, especially on higher bands. Values
						below 70 mean limited DX opportunities.
					</p>
					<p>
						<strong class="text-slate-300">Sunspot Number (SSN):</strong> More sunspots correlate with
						better HF conditions. Solar maximum can have 200+ sunspots.
					</p>
					<p>
						<strong class="text-slate-300">K-Index:</strong> Measures geomagnetic activity on a scale
						of 0-9. Values of 0-2 are ideal. Above 4 indicates disturbed conditions that degrade propagation.
					</p>
				</div>
				<div class="space-y-3 text-sm text-slate-400">
					<p>
						<strong class="text-slate-300">A-Index:</strong> Average geomagnetic activity over 24
						hours. Values below 10 are excellent for DX. Above 20 means poor conditions.
					</p>
					<p>
						<strong class="text-slate-300">X-Ray Flux:</strong> Indicates solar flare activity. Class
						M and X flares can cause radio blackouts but may enhance propagation later.
					</p>
					<p class="text-xs italic">
						Data refreshes every 15 minutes from NOAA Space Weather Prediction Center. Band
						conditions use rule-based analysis considering solar flux, K-index, and time of day.
					</p>
				</div>
			</div>
		</div>
	{/if}
</div>
