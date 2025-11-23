<script lang="ts">
	let sdrs = $state([
		{
			id: '1',
			name: 'WebSDR - University of Twente',
			location: 'Netherlands',
			frequency: '7.100 MHz',
			status: 'online'
		},
		{
			id: '2',
			name: 'KiwiSDR - W1AW',
			location: 'Connecticut, USA',
			frequency: '14.200 MHz',
			status: 'online'
		},
		{
			id: '3',
			name: 'WebSDR - K3FEF',
			location: 'Pennsylvania, USA',
			frequency: '3.850 MHz',
			status: 'online'
		}
	]);
</script>

<svelte:head>
	<title>Web-SDR Directory - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<div class="mb-8">
		<h1 class="text-3xl font-bold">Web-SDR Directory</h1>
		<p class="mt-2 text-slate-400">
			Access public software-defined radio receivers around the world
		</p>
	</div>

	<!-- Info Banner -->
	<div class="card mb-8 border-primary-700 bg-primary-900/10">
		<h3 class="mb-2 text-lg font-semibold">About Web-SDR</h3>
		<p class="text-slate-400">
			Web-SDRs allow you to listen to HF, VHF, and UHF frequencies remotely through your browser.
			These are read-only receivers operated by amateur radio operators and institutions worldwide.
		</p>
		<p class="mt-2 text-sm text-slate-500">
			Note: Private SDR hosting and transmit capabilities will be available in future tiers.
		</p>
	</div>

	<!-- Search and Filters -->
	<div class="card mb-6">
		<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Search</label>
				<input type="text" placeholder="Search by name or location..." class="input" />
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Band</label>
				<select class="input">
					<option value="">All Bands</option>
					<option value="hf">HF (3-30 MHz)</option>
					<option value="vhf">VHF (30-300 MHz)</option>
					<option value="uhf">UHF (300-3000 MHz)</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Status</label>
				<select class="input">
					<option value="">All</option>
					<option value="online">Online</option>
					<option value="offline">Offline</option>
				</select>
			</div>
		</div>
	</div>

	<!-- SDR Grid -->
	<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
		{#each sdrs as sdr}
			<div class="card hover:border-primary-700 transition-colors">
				<div class="mb-4 flex items-start justify-between">
					<div>
						<h3 class="text-lg font-semibold">{sdr.name}</h3>
						<p class="mt-1 text-sm text-slate-400">📍 {sdr.location}</p>
					</div>
					<span class="rounded-full bg-green-500/20 px-2 py-1 text-xs text-green-400">
						{sdr.status}
					</span>
				</div>
				<div class="mb-4 font-mono text-sm text-primary-400">{sdr.frequency}</div>
				<a href="#" class="btn-primary w-full text-center"> Connect to SDR </a>
			</div>
		{/each}
	</div>

	<!-- Coming Soon -->
	<div class="card mt-8 border-slate-700 bg-slate-800/30">
		<h3 class="mb-4 text-xl font-semibold">🚀 Coming Soon</h3>
		<ul class="space-y-2 text-slate-400">
			<li>✓ Host your own private SDR receiver</li>
			<li>✓ Docker-based SDR containers for easy deployment</li>
			<li>✓ Remote transmit capability (with strict controls)</li>
			<li>✓ DSP toolbox: filters, noise reduction, decoders</li>
		</ul>
	</div>
</div>
