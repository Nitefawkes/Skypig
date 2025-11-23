<script lang="ts">
	import { onMount } from 'svelte';

	interface Stats {
		total_qsos: number;
		qso_limit: number;
		remaining_qsos: number;
	}

	let stats = $state<Stats | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		await fetchStats();
	});

	async function fetchStats() {
		try {
			loading = true;
			const response = await fetch('/api/qsos/stats');
			if (!response.ok) {
				throw new Error('Failed to fetch stats');
			}
			const result = await response.json();
			stats = result.data;
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			loading = false;
		}
	}

	function getUsagePercent(): number {
		if (!stats || stats.qso_limit === -1) return 0;
		return (stats.total_qsos / stats.qso_limit) * 100;
	}

	function getUsageColor(): string {
		const percent = getUsagePercent();
		if (percent >= 90) return 'bg-red-500';
		if (percent >= 75) return 'bg-yellow-500';
		return 'bg-green-500';
	}

	export function refresh() {
		fetchStats();
	}
</script>

<div class="card">
	<h2 class="mb-4 text-xl font-bold text-gray-900">Logbook Statistics</h2>

	{#if loading}
		<div class="flex items-center justify-center py-8">
			<div class="text-gray-500">Loading stats...</div>
		</div>
	{:else if error}
		<div class="rounded-lg bg-red-50 p-4 text-red-600">
			<strong>Error:</strong> {error}
		</div>
	{:else if stats}
		<!-- Total QSOs -->
		<div class="mb-6">
			<div class="flex items-end justify-between">
				<div>
					<div class="text-sm font-medium text-gray-500">Total QSOs</div>
					<div class="text-4xl font-bold text-primary-600">{stats.total_qsos.toLocaleString()}</div>
				</div>
				{#if stats.qso_limit !== -1}
					<div class="text-right">
						<div class="text-sm text-gray-500">Limit</div>
						<div class="text-xl font-semibold text-gray-700">
							{stats.qso_limit.toLocaleString()}
						</div>
					</div>
				{:else}
					<div class="rounded bg-green-100 px-3 py-1 text-sm font-semibold text-green-800">
						Unlimited
					</div>
				{/if}
			</div>

			<!-- Progress Bar (if not unlimited) -->
			{#if stats.qso_limit !== -1}
				<div class="mt-4">
					<div class="h-2 w-full overflow-hidden rounded-full bg-gray-200">
						<div
							class="h-full transition-all {getUsageColor()}"
							style="width: {getUsagePercent()}%"
						></div>
					</div>
					<div class="mt-1 flex justify-between text-xs text-gray-500">
						<span>{getUsagePercent().toFixed(1)}% used</span>
						{#if stats.remaining_qsos >= 0}
							<span>{stats.remaining_qsos.toLocaleString()} remaining</span>
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<!-- Quick Stats Grid -->
		<div class="grid grid-cols-2 gap-4 border-t border-gray-200 pt-4">
			<div>
				<div class="text-xs font-medium text-gray-500">Today</div>
				<div class="text-2xl font-bold text-gray-900">-</div>
			</div>
			<div>
				<div class="text-xs font-medium text-gray-500">This Week</div>
				<div class="text-2xl font-bold text-gray-900">-</div>
			</div>
			<div>
				<div class="text-xs font-medium text-gray-500">Countries</div>
				<div class="text-2xl font-bold text-gray-900">-</div>
			</div>
			<div>
				<div class="text-xs font-medium text-gray-500">States</div>
				<div class="text-2xl font-bold text-gray-900">-</div>
			</div>
		</div>

		<!-- Warning if approaching limit -->
		{#if stats.qso_limit !== -1 && stats.remaining_qsos < 50}
			<div class="mt-4 rounded-lg bg-yellow-50 p-3 text-sm text-yellow-800">
				<strong>⚠️ Warning:</strong> You're approaching your QSO limit. Consider upgrading to the Operator
				tier for 20,000 QSOs.
			</div>
		{/if}
	{/if}
</div>
