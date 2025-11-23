<script lang="ts">
	import { onMount } from 'svelte';

	let stats = $state({
		qsos: 0,
		users: 0,
		countries: 0
	});

	onMount(async () => {
		// Fetch stats from API
		try {
			const response = await fetch('/api/v1/health');
			const data = await response.json();
			console.log('API Health:', data);
		} catch (error) {
			console.error('Failed to fetch stats:', error);
		}
	});
</script>

<svelte:head>
	<title>Ham Radio Cloud - Cloud Logbook & Propagation Platform</title>
</svelte:head>

<!-- Hero Section -->
<div class="relative overflow-hidden">
	<div
		class="absolute inset-0 bg-gradient-to-br from-primary-900/20 via-slate-950 to-slate-950"
	></div>
	<div class="relative mx-auto max-w-7xl px-4 py-24 sm:px-6 lg:px-8">
		<div class="text-center">
			<h1 class="mb-6 text-6xl font-bold tracking-tight">
				<span class="text-white">Your Station.</span>
				<span class="text-primary-500">In the Cloud.</span>
			</h1>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-slate-300">
				Cloud logbook, real-time propagation insights, and SDR integration for amateur radio
				operators. Auto-sync with LoTW, track awards, and access your logs anywhere.
			</p>
			<div class="flex justify-center gap-4">
				<a href="/signup" class="btn-primary text-lg">Start Free Trial</a>
				<a href="/demo" class="btn-secondary text-lg">Watch Demo</a>
			</div>
		</div>

		<!-- Stats -->
		<div class="mt-20 grid grid-cols-1 gap-8 sm:grid-cols-3">
			<div class="card text-center">
				<div class="text-4xl font-bold text-primary-500">{stats.qsos.toLocaleString()}+</div>
				<div class="mt-2 text-slate-400">QSOs Logged</div>
			</div>
			<div class="card text-center">
				<div class="text-4xl font-bold text-primary-500">{stats.users.toLocaleString()}+</div>
				<div class="mt-2 text-slate-400">Active Operators</div>
			</div>
			<div class="card text-center">
				<div class="text-4xl font-bold text-primary-500">{stats.countries}</div>
				<div class="mt-2 text-slate-400">Countries Reached</div>
			</div>
		</div>
	</div>
</div>

<!-- Features Section -->
<div class="mx-auto max-w-7xl px-4 py-24 sm:px-6 lg:px-8">
	<h2 class="mb-12 text-center">Core Features</h2>
	<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
		<div class="card">
			<div class="mb-4 text-4xl">📝</div>
			<h3 class="mb-2 text-xl font-semibold">Cloud Logbook</h3>
			<p class="text-slate-400">
				Manual entry, ADIF import/export, auto-sync with LoTW. Your logs backed up and accessible
				from any device.
			</p>
		</div>
		<div class="card">
			<div class="mb-4 text-4xl">📡</div>
			<h3 class="mb-2 text-xl font-semibold">Propagation Engine</h3>
			<p class="text-slate-400">
				Real-time solar data, Kp index, and band condition forecasts. Get alerts when bands open
				up.
			</p>
		</div>
		<div class="card">
			<div class="mb-4 text-4xl">🌐</div>
			<h3 class="mb-2 text-xl font-semibold">Web-SDR Integration</h3>
			<p class="text-slate-400">
				Access public KiwiSDR receivers worldwide. Listen, monitor, and discover new frequencies.
			</p>
		</div>
		<div class="card">
			<div class="mb-4 text-4xl">🏆</div>
			<h3 class="mb-2 text-xl font-semibold">Award Tracking</h3>
			<p class="text-slate-400">
				Track DXCC, WAS, WAZ, and other awards. See your progress and what you need next.
			</p>
		</div>
		<div class="card">
			<div class="mb-4 text-4xl">📱</div>
			<h3 class="mb-2 text-xl font-semibold">Mobile PWA</h3>
			<p class="text-slate-400">
				Install as a mobile app. Log contacts from field days, POTA/SOTA activations, and remote
				operations.
			</p>
		</div>
		<div class="card">
			<div class="mb-4 text-4xl">🔌</div>
			<h3 class="mb-2 text-xl font-semibold">API Access</h3>
			<p class="text-slate-400">
				REST & WebSocket API for integrations. Connect your rig control software and hardware.
			</p>
		</div>
	</div>
</div>

<!-- Pricing Teaser -->
<div class="bg-slate-900/30 py-24">
	<div class="mx-auto max-w-7xl px-4 text-center sm:px-6 lg:px-8">
		<h2 class="mb-6">Simple, Transparent Pricing</h2>
		<p class="mb-12 text-xl text-slate-400">Start free. Upgrade when you need more.</p>
		<div class="grid grid-cols-1 gap-8 md:grid-cols-3">
			<div class="card">
				<h3 class="mb-4">Free</h3>
				<div class="mb-6 text-5xl font-bold">$0</div>
				<ul class="space-y-3 text-left text-slate-400">
					<li>✓ 500 QSOs limit</li>
					<li>✓ Basic propagation display</li>
					<li>✓ Public SDR links</li>
					<li>✓ ADIF import/export</li>
				</ul>
			</div>
			<div class="card border-2 border-primary-500">
				<div class="mb-2 text-xs font-semibold uppercase text-primary-500">Most Popular</div>
				<h3 class="mb-4">Operator</h3>
				<div class="mb-6 text-5xl font-bold">$7<span class="text-xl text-slate-500">/mo</span></div>
				<ul class="space-y-3 text-left text-slate-400">
					<li>✓ 20,000 QSOs</li>
					<li>✓ LoTW auto-sync</li>
					<li>✓ Propagation alerts</li>
					<li>✓ Club groups</li>
					<li>✓ Priority support</li>
				</ul>
			</div>
			<div class="card opacity-75">
				<div class="mb-2 text-xs font-semibold uppercase text-slate-500">Coming Soon</div>
				<h3 class="mb-4">Contester</h3>
				<div class="mb-6 text-5xl font-bold">$19<span class="text-xl text-slate-500">/mo</span></div>
				<ul class="space-y-3 text-left text-slate-400">
					<li>✓ Unlimited QSOs</li>
					<li>✓ Private SDR hosting</li>
					<li>✓ API access</li>
					<li>✓ Advanced analytics</li>
					<li>✓ Premium support</li>
				</ul>
			</div>
		</div>
	</div>
</div>
