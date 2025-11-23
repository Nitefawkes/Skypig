<script lang="ts">
	import { onMount } from 'svelte';
	import type { QSO } from '$types';
	import api from '$lib/utils/api';

	let qsos = $state<QSO[]>([]);
	let loading = $state(true);
	let showAddModal = $state(false);

	onMount(async () => {
		try {
			// TODO: Get actual QSOs from API with auth token
			qsos = [];
			loading = false;
		} catch (error) {
			console.error('Failed to load QSOs:', error);
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Logbook - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Logbook</h1>
			<p class="mt-2 text-slate-400">View and manage your QSO log</p>
		</div>
		<div class="flex gap-4">
			<button class="btn-secondary">Import ADIF</button>
			<button class="btn-secondary">Export ADIF</button>
			<button class="btn-primary" onclick={() => (showAddModal = true)}>Log QSO</button>
		</div>
	</div>

	<!-- Filters -->
	<div class="card mb-6">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-4">
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Callsign</label>
				<input type="text" placeholder="Search callsign..." class="input" />
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Band</label>
				<select class="input">
					<option value="">All Bands</option>
					<option value="160m">160m</option>
					<option value="80m">80m</option>
					<option value="40m">40m</option>
					<option value="20m">20m</option>
					<option value="15m">15m</option>
					<option value="10m">10m</option>
					<option value="6m">6m</option>
					<option value="2m">2m</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Mode</label>
				<select class="input">
					<option value="">All Modes</option>
					<option value="SSB">SSB</option>
					<option value="CW">CW</option>
					<option value="FT8">FT8</option>
					<option value="FT4">FT4</option>
					<option value="RTTY">RTTY</option>
					<option value="PSK31">PSK31</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Date Range</label>
				<input type="date" class="input" />
			</div>
		</div>
	</div>

	<!-- QSO List -->
	{#if loading}
		<div class="card text-center py-12">
			<div class="text-slate-400">Loading QSOs...</div>
		</div>
	{:else if qsos.length === 0}
		<div class="card text-center py-12">
			<div class="text-6xl mb-4">📝</div>
			<h3 class="text-xl font-semibold mb-2">No QSOs Yet</h3>
			<p class="text-slate-400 mb-6">Start logging your contacts or import from ADIF</p>
			<button class="btn-primary" onclick={() => (showAddModal = true)}>Log Your First QSO</button>
		</div>
	{:else}
		<div class="card overflow-x-auto">
			<table class="w-full">
				<thead class="border-b border-slate-700">
					<tr class="text-left text-sm text-slate-400">
						<th class="pb-3">Date/Time</th>
						<th class="pb-3">Callsign</th>
						<th class="pb-3">Frequency</th>
						<th class="pb-3">Band</th>
						<th class="pb-3">Mode</th>
						<th class="pb-3">RST</th>
						<th class="pb-3">Grid</th>
						<th class="pb-3">LoTW</th>
						<th class="pb-3"></th>
					</tr>
				</thead>
				<tbody>
					{#each qsos as qso}
						<tr class="border-b border-slate-800/50 hover:bg-slate-800/30">
							<td class="py-3 font-mono text-sm">
								{new Date(qso.time_on).toLocaleString()}
							</td>
							<td class="py-3 font-bold text-primary-400">{qso.callsign}</td>
							<td class="py-3 font-mono text-sm">{qso.frequency.toFixed(3)}</td>
							<td class="py-3">{qso.band}</td>
							<td class="py-3">{qso.mode}</td>
							<td class="py-3 font-mono text-sm">{qso.rst_sent}/{qso.rst_received}</td>
							<td class="py-3">{qso.grid_square || '-'}</td>
							<td class="py-3">
								{#if qso.lotw_confirmed}
									<span class="text-green-500">✓</span>
								{:else if qso.lotw_sent}
									<span class="text-yellow-500">⋯</span>
								{:else}
									<span class="text-slate-600">-</span>
								{/if}
							</td>
							<td class="py-3">
								<button class="text-slate-400 hover:text-white">Edit</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
