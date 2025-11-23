<script lang="ts">
	import { onMount } from 'svelte';
	import type { QSO } from '$types';
	import api from '$lib/utils/api';
	import QSOForm from '$lib/components/QSOForm.svelte';

	let qsos = $state<QSO[]>([]);
	let loading = $state(true);
	let showAddModal = $state(false);
	let showEditModal = $state(false);
	let editingQSO = $state<QSO | null>(null);
	let error = $state<string | null>(null);
	let success = $state<string | null>(null);

	// Filters
	let filters = $state({
		callsign: '',
		band: '',
		mode: '',
		start_date: '',
		end_date: ''
	});

	onMount(async () => {
		await loadQSOs();
	});

	async function loadQSOs() {
		loading = true;
		error = null;
		try {
			const response = await fetch(`/api/v1/qso?${new URLSearchParams(filters as any)}`);
			const data = await response.json();
			qsos = data.data || [];
		} catch (err: any) {
			error = err.message || 'Failed to load QSOs';
			qsos = [];
		} finally {
			loading = false;
		}
	}

	async function handleCreateQSO(event: CustomEvent) {
		try {
			const response = await fetch('/api/v1/qso', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(event.detail)
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to create QSO');
			}

			success = 'QSO logged successfully!';
			showAddModal = false;
			await loadQSOs();

			setTimeout(() => (success = null), 3000);
		} catch (err: any) {
			error = err.message;
			setTimeout(() => (error = null), 5000);
		}
	}

	async function handleUpdateQSO(event: CustomEvent) {
		if (!editingQSO) return;

		try {
			const response = await fetch(`/api/v1/qso/${editingQSO.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(event.detail)
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.message || 'Failed to update QSO');
			}

			success = 'QSO updated successfully!';
			showEditModal = false;
			editingQSO = null;
			await loadQSOs();

			setTimeout(() => (success = null), 3000);
		} catch (err: any) {
			error = err.message;
			setTimeout(() => (error = null), 5000);
		}
	}

	async function handleDeleteQSO(id: string) {
		if (!confirm('Are you sure you want to delete this QSO?')) return;

		try {
			const response = await fetch(`/api/v1/qso/${id}`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error('Failed to delete QSO');
			}

			success = 'QSO deleted successfully!';
			await loadQSOs();
			setTimeout(() => (success = null), 3000);
		} catch (err: any) {
			error = err.message;
			setTimeout(() => (error = null), 5000);
		}
	}

	function handleEditClick(qso: QSO) {
		editingQSO = qso;
		showEditModal = true;
	}

	async function handleImportADIF() {
		const input = document.createElement('input');
		input.type = 'file';
		input.accept = '.adi,.adif';

		input.onchange = async (e: any) => {
			const file = e.target.files[0];
			if (!file) return;

			const formData = new FormData();
			formData.append('file', file);

			try {
				const response = await fetch('/api/v1/qso/import/adif', {
					method: 'POST',
					body: formData
				});

				const result = await response.json();

				if (response.ok) {
					success = `Imported ${result.imported} QSOs (${result.skipped} skipped)`;
					await loadQSOs();
					setTimeout(() => (success = null), 5000);
				} else {
					throw new Error(result.message || 'Import failed');
				}
			} catch (err: any) {
				error = err.message;
				setTimeout(() => (error = null), 5000);
			}
		};

		input.click();
	}

	async function handleExportADIF() {
		try {
			const queryParams = new URLSearchParams(filters as any);
			const response = await fetch(`/api/v1/qso/export/adif?${queryParams}`);

			if (!response.ok) {
				throw new Error('Export failed');
			}

			const blob = await response.blob();
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `logbook_${new Date().toISOString().split('T')[0]}.adi`;
			a.click();
			window.URL.revokeObjectURL(url);

			success = 'Logbook exported successfully!';
			setTimeout(() => (success = null), 3000);
		} catch (err: any) {
			error = err.message;
			setTimeout(() => (error = null), 5000);
		}
	}

	async function handleFilterChange() {
		await loadQSOs();
	}

	function clearFilters() {
		filters = {
			callsign: '',
			band: '',
			mode: '',
			start_date: '',
			end_date: ''
		};
		loadQSOs();
	}
</script>

<svelte:head>
	<title>Logbook - Ham Radio Cloud</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<!-- Header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">Logbook</h1>
			<p class="mt-2 text-slate-400">View and manage your QSO log</p>
		</div>
		<div class="flex gap-4">
			<button class="btn-secondary" onclick={handleImportADIF}>Import ADIF</button>
			<button class="btn-secondary" onclick={handleExportADIF}>Export ADIF</button>
			<button class="btn-primary" onclick={() => (showAddModal = true)}>Log QSO</button>
		</div>
	</div>

	<!-- Success/Error Messages -->
	{#if success}
		<div class="mb-6 rounded-lg border border-green-700 bg-green-900/20 p-4 text-green-400">
			{success}
		</div>
	{/if}

	{#if error}
		<div class="mb-6 rounded-lg border border-red-700 bg-red-900/20 p-4 text-red-400">
			{error}
		</div>
	{/if}

	<!-- Filters -->
	<div class="card mb-6">
		<div class="mb-4 flex items-center justify-between">
			<h3 class="text-lg font-semibold">Filters</h3>
			<button class="text-sm text-primary-500 hover:underline" onclick={clearFilters}>
				Clear All
			</button>
		</div>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-5">
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Callsign</label>
				<input
					type="text"
					placeholder="Search..."
					bind:value={filters.callsign}
					onchange={handleFilterChange}
					class="input"
				/>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Band</label>
				<select bind:value={filters.band} onchange={handleFilterChange} class="input">
					<option value="">All Bands</option>
					<option value="20m">20m</option>
					<option value="40m">40m</option>
					<option value="80m">80m</option>
					<option value="15m">15m</option>
					<option value="10m">10m</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Mode</label>
				<select bind:value={filters.mode} onchange={handleFilterChange} class="input">
					<option value="">All Modes</option>
					<option value="FT8">FT8</option>
					<option value="SSB">SSB</option>
					<option value="CW">CW</option>
					<option value="FT4">FT4</option>
				</select>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">Start Date</label>
				<input
					type="date"
					bind:value={filters.start_date}
					onchange={handleFilterChange}
					class="input"
				/>
			</div>
			<div>
				<label class="mb-2 block text-sm font-medium text-slate-300">End Date</label>
				<input
					type="date"
					bind:value={filters.end_date}
					onchange={handleFilterChange}
					class="input"
				/>
			</div>
		</div>
	</div>

	<!-- QSO List -->
	{#if loading}
		<div class="card py-12 text-center">
			<div class="text-slate-400">Loading QSOs...</div>
		</div>
	{:else if qsos.length === 0}
		<div class="card py-12 text-center">
			<div class="mb-4 text-6xl">📝</div>
			<h3 class="mb-2 text-xl font-semibold">No QSOs Yet</h3>
			<p class="mb-6 text-slate-400">Start logging your contacts or import from ADIF</p>
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
								<div class="flex gap-2">
									<button
										class="text-sm text-slate-400 hover:text-white"
										onclick={() => handleEditClick(qso)}
									>
										Edit
									</button>
									<button
										class="text-sm text-red-400 hover:text-red-300"
										onclick={() => handleDeleteQSO(qso.id)}
									>
										Delete
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			<div class="mt-4 text-sm text-slate-500">
				Showing {qsos.length} QSO{qsos.length !== 1 ? 's' : ''}
			</div>
		</div>
	{/if}
</div>

<!-- Add QSO Modal -->
{#if showAddModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="card max-h-[90vh] w-full max-w-2xl overflow-y-auto">
			<h2 class="mb-6 text-2xl font-bold">Log New QSO</h2>
			<QSOForm on:submit={handleCreateQSO} on:cancel={() => (showAddModal = false)} />
		</div>
	</div>
{/if}

<!-- Edit QSO Modal -->
{#if showEditModal && editingQSO}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="card max-h-[90vh] w-full max-w-2xl overflow-y-auto">
			<h2 class="mb-6 text-2xl font-bold">Edit QSO</h2>
			<QSOForm
				qso={editingQSO}
				mode="edit"
				on:submit={handleUpdateQSO}
				on:cancel={() => {
					showEditModal = false;
					editingQSO = null;
				}}
			/>
		</div>
	</div>
{/if}
