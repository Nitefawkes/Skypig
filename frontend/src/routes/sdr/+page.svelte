<script lang="ts">
	import { onMount } from 'svelte';

	interface SDR {
		id: string;
		name: string;
		callsign?: string;
		url: string;
		type: string;
		location?: string;
		grid_square?: string;
		country?: string;
		bands: string[];
		modes: string[];
		antenna_info?: string;
		frequency_min?: number;
		frequency_max?: number;
		users_max?: number;
		status: string;
		last_seen?: string;
	}

	interface SDRStats {
		total: number;
		online: number;
		by_type: {
			kiwisdr: number;
			websdr: number;
			openwebrx: number;
		};
	}

	let sdrs = $state<SDR[]>([]);
	let stats = $state<SDRStats | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let searchQuery = $state('');
	let filterType = $state('');
	let filterCountry = $state('');
	let filterStatus = $state('online');
	let total = $state(0);
	let limit = $state(50);
	let offset = $state(0);
	let currentPage = $state(1);

	onMount(async () => {
		await Promise.all([loadSDRs(), loadStats()]);
	});

	async function loadSDRs() {
		loading = true;
		error = null;
		try {
			const params = new URLSearchParams({
				limit: limit.toString(),
				offset: offset.toString()
			});

			if (filterType) params.append('type', filterType);
			if (filterCountry) params.append('country', filterCountry);
			if (filterStatus) params.append('status', filterStatus);
			if (searchQuery) params.append('search', searchQuery);

			const response = await fetch(`/api/v1/sdr?${params.toString()}`);
			if (!response.ok) {
				throw new Error('Failed to fetch SDR directory');
			}

			const data = await response.json();
			sdrs = data.sdrs || [];
			total = data.total || 0;
		} catch (err: any) {
			error = err.message || 'Failed to load SDR directory';
		} finally {
			loading = false;
		}
	}

	async function loadStats() {
		try {
			const response = await fetch('/api/v1/sdr/stats');
			if (response.ok) {
				stats = await response.json();
			}
		} catch (err) {
			console.error('Failed to load stats:', err);
		}
	}

	async function handleRefresh() {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/v1/sdr/refresh', {
				method: 'POST'
			});

			if (!response.ok) {
				throw new Error('Failed to refresh directory');
			}

			// Reload after refresh
			await loadSDRs();
		} catch (err: any) {
			error = err.message || 'Failed to refresh directory';
			loading = false;
		}
	}

	async function handleSearch() {
		offset = 0;
		currentPage = 1;
		await loadSDRs();
	}

	function handleClearFilters() {
		searchQuery = '';
		filterType = '';
		filterCountry = '';
		filterStatus = 'online';
		offset = 0;
		currentPage = 1;
		loadSDRs();
	}

	function handlePageChange(newPage: number) {
		currentPage = newPage;
		offset = (newPage - 1) * limit;
		loadSDRs();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	function getTypeColor(type: string) {
		switch (type) {
			case 'kiwisdr':
				return 'bg-blue-600';
			case 'websdr':
				return 'bg-green-600';
			case 'openwebrx':
				return 'bg-purple-600';
			default:
				return 'bg-slate-600';
		}
	}

	function getStatusColor(status: string) {
		switch (status) {
			case 'online':
				return 'text-green-500';
			case 'offline':
				return 'text-red-500';
			default:
				return 'text-slate-500';
		}
	}

	$effect(() => {
		// Reactive search
		if (searchQuery !== undefined || filterType !== undefined || filterCountry !== undefined) {
			const timeout = setTimeout(() => {
				handleSearch();
			}, 500);
			return () => clearTimeout(timeout);
		}
	});

	const totalPages = $derived(Math.ceil(total / limit));
</script>

<svelte:head>
	<title>WebSDR Directory - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<!-- Header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">WebSDR Directory</h1>
			<p class="mt-2 text-slate-400">
				Browse and connect to software-defined radio receivers worldwide
			</p>
		</div>
		<button class="btn-primary" onclick={handleRefresh} disabled={loading}>
			{loading ? 'Refreshing...' : 'Refresh Directory'}
		</button>
	</div>

	<!-- Stats Cards -->
	{#if stats}
		<div class="mb-8 grid grid-cols-2 gap-4 md:grid-cols-4">
			<div class="card">
				<div class="text-sm text-slate-400">Total SDRs</div>
				<div class="mt-2 text-3xl font-bold text-primary-500">{stats.total}</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">Online Now</div>
				<div class="mt-2 text-3xl font-bold text-green-500">{stats.online}</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">KiwiSDR</div>
				<div class="mt-2 text-3xl font-bold text-blue-500">{stats.by_type.kiwisdr}</div>
			</div>
			<div class="card">
				<div class="text-sm text-slate-400">WebSDR</div>
				<div class="mt-2 text-3xl font-bold text-green-500">{stats.by_type.websdr}</div>
			</div>
		</div>
	{/if}

	<!-- Search and Filters -->
	<div class="card mb-6">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-4">
			<div class="md:col-span-2">
				<label class="mb-2 block text-sm font-medium text-slate-300">Search</label>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search by name, location, or callsign..."
					class="input w-full"
				/>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Type</label>
				<select bind:value={filterType} class="input w-full">
					<option value="">All Types</option>
					<option value="kiwisdr">KiwiSDR</option>
					<option value="websdr">WebSDR</option>
					<option value="openwebrx">OpenWebRX</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Status</label>
				<select bind:value={filterStatus} class="input w-full">
					<option value="">All Status</option>
					<option value="online">Online</option>
					<option value="offline">Offline</option>
				</select>
			</div>
		</div>
		{#if searchQuery || filterType || filterCountry}
			<div class="mt-4">
				<button onclick={handleClearFilters} class="btn-secondary text-sm">
					Clear Filters
				</button>
			</div>
		{/if}
	</div>

	<!-- Error Message -->
	{#if error}
		<div class="mb-6 rounded-lg border border-red-700 bg-red-900/20 p-4 text-red-400">
			{error}
		</div>
	{/if}

	<!-- SDR List -->
	{#if loading && sdrs.length === 0}
		<div class="card py-12 text-center">
			<div class="text-slate-400">Loading SDR directory...</div>
		</div>
	{:else if sdrs.length === 0}
		<div class="card py-12 text-center">
			<div class="text-slate-400">No SDRs found matching your criteria</div>
		</div>
	{:else}
		<div class="mb-6 grid grid-cols-1 gap-4 lg:grid-cols-2">
			{#each sdrs as sdr}
				<div class="card hover:border-primary-700 transition-colors">
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<div class="flex items-center gap-3">
								<h3 class="text-xl font-bold">{sdr.name}</h3>
								<span class="rounded px-2 py-1 text-xs font-semibold text-white {getTypeColor(sdr.type)}">
									{sdr.type}
								</span>
								<span class="text-sm font-semibold {getStatusColor(sdr.status)}">
									● {sdr.status}
								</span>
							</div>

							{#if sdr.callsign}
								<div class="mt-1 text-sm text-primary-400 font-mono">{sdr.callsign}</div>
							{/if}

							{#if sdr.location}
								<div class="mt-2 flex items-center gap-2 text-sm text-slate-400">
									<svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"></path>
									</svg>
									{sdr.location}
									{#if sdr.grid_square}
										<span class="font-mono text-xs">({sdr.grid_square})</span>
									{/if}
								</div>
							{/if}

							{#if sdr.bands && sdr.bands.length > 0}
								<div class="mt-3 flex flex-wrap gap-1">
									{#each sdr.bands.slice(0, 8) as band}
										<span class="rounded bg-slate-700 px-2 py-1 text-xs font-medium">{band}</span>
									{/each}
									{#if sdr.bands.length > 8}
										<span class="rounded bg-slate-700 px-2 py-1 text-xs font-medium">+{sdr.bands.length - 8} more</span>
									{/if}
								</div>
							{/if}

							{#if sdr.frequency_min && sdr.frequency_max}
								<div class="mt-2 text-xs text-slate-500">
									{sdr.frequency_min.toFixed(3)} - {sdr.frequency_max.toFixed(3)} MHz
								</div>
							{/if}

							{#if sdr.antenna_info}
								<div class="mt-2 text-xs text-slate-500">
									📡 {sdr.antenna_info}
								</div>
							{/if}
						</div>
					</div>

					<div class="mt-4 flex items-center justify-between border-t border-slate-700 pt-4">
						<div class="text-xs text-slate-500">
							{#if sdr.users_max}
								Max users: {sdr.users_max}
							{/if}
						</div>
						<a
							href={sdr.url}
							target="_blank"
							rel="noopener noreferrer"
							class="btn-primary text-sm"
						>
							Open SDR →
						</a>
					</div>
				</div>
			{/each}
		</div>

		<!-- Pagination -->
		{#if totalPages > 1}
			<div class="card flex items-center justify-between">
				<div class="text-sm text-slate-400">
					Showing {offset + 1}-{Math.min(offset + limit, total)} of {total} SDRs
				</div>
				<div class="flex gap-2">
					<button
						onclick={() => handlePageChange(currentPage - 1)}
						disabled={currentPage === 1}
						class="btn-secondary disabled:opacity-50"
					>
						← Previous
					</button>
					<div class="flex items-center gap-1">
						{#each Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
							let page;
							if (totalPages <= 5) {
								page = i + 1;
							} else if (currentPage <= 3) {
								page = i + 1;
							} else if (currentPage >= totalPages - 2) {
								page = totalPages - 4 + i;
							} else {
								page = currentPage - 2 + i;
							}
							return page;
						}) as page}
							<button
								onclick={() => handlePageChange(page)}
								class="btn-secondary {currentPage === page ? 'bg-primary-600' : ''}"
							>
								{page}
							</button>
						{/each}
					</div>
					<button
						onclick={() => handlePageChange(currentPage + 1)}
						disabled={currentPage === totalPages}
						class="btn-secondary disabled:opacity-50"
					>
						Next →
					</button>
				</div>
			</div>
		{/if}
	{/if}

	<!-- Info Panel -->
	<div class="card mt-8 border-primary-700 bg-primary-900/10">
		<h3 class="mb-4 text-xl font-semibold">About WebSDR</h3>
		<div class="grid grid-cols-1 gap-4 text-sm text-slate-400 md:grid-cols-2">
			<div class="space-y-3">
				<p>
					<strong class="text-slate-300">What is WebSDR?</strong> Software-Defined Radio (SDR)
					receivers accessible through your web browser. Listen to HF, VHF, and UHF bands from
					around the world without owning radio equipment.
				</p>
				<p>
					<strong class="text-slate-300">KiwiSDR:</strong> Popular SDR platform covering 0-30 MHz
					(HF bands). Great for shortwave listening, amateur radio monitoring, and propagation
					research.
				</p>
			</div>
			<div class="space-y-3">
				<p>
					<strong class="text-slate-300">How to use:</strong> Click "Open SDR" to launch the
					receiver in a new tab. Most SDRs support multiple simultaneous users. Be courteous and
					share the spectrum!
				</p>
				<p class="text-xs italic">
					Directory automatically refreshes every 6 hours from public WebSDR networks. Status and
					availability may vary.
				</p>
			</div>
		</div>
	</div>
</div>
