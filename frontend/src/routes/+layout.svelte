<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';

	let { children } = $props();
</script>

<div class="min-h-screen bg-slate-950">
	<!-- Navigation -->
	<nav class="border-b border-slate-800 bg-slate-900/50 backdrop-blur-lg">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<a href="/" class="flex items-center space-x-2">
						<div class="text-2xl">📻</div>
						<span class="text-xl font-bold text-primary-500">Ham Radio Cloud</span>
					</a>
					<div class="ml-10 flex space-x-4">
						<a
							href="/logbook"
							class="rounded-md px-3 py-2 text-sm font-medium transition-colors {$page.url
								.pathname === '/logbook'
								? 'bg-slate-800 text-white'
								: 'text-slate-300 hover:bg-slate-800 hover:text-white'}"
						>
							Logbook
						</a>
						<a
							href="/propagation"
							class="rounded-md px-3 py-2 text-sm font-medium transition-colors {$page.url
								.pathname === '/propagation'
								? 'bg-slate-800 text-white'
								: 'text-slate-300 hover:bg-slate-800 hover:text-white'}"
						>
							Propagation
						</a>
						<a
							href="/sdr"
							class="rounded-md px-3 py-2 text-sm font-medium transition-colors {$page.url
								.pathname === '/sdr'
								? 'bg-slate-800 text-white'
								: 'text-slate-300 hover:bg-slate-800 hover:text-white'}"
						>
							SDR
						</a>
					</div>
				</div>
				<div class="flex items-center space-x-4">
					<a href="/login" class="btn-secondary">Sign In</a>
				</div>
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	<main>
		{@render children()}
	</main>

	<!-- Footer -->
	<footer class="mt-20 border-t border-slate-800 bg-slate-900/30 py-8">
		<div class="mx-auto max-w-7xl px-4 text-center text-sm text-slate-500">
			<p>Ham Radio Cloud v1.0 - Built for the amateur radio community</p>
			<p class="mt-2">73 de W1AW</p>
		</div>
	</footer>
</div>
