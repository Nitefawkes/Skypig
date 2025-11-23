<script lang="ts">
	import { onMount } from 'svelte';
	import Navigation from '$components/Navigation.svelte';
	import QSOTable from '$components/QSOTable.svelte';
	import QSOForm from '$components/QSOForm.svelte';
	import QSOStats from '$components/QSOStats.svelte';
	import ADIFManager from '$components/ADIFManager.svelte';
	import type { QSO } from '$types';

	let qsos = $state<QSO[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let showForm = $state(false);
	let editingQSO = $state<QSO | null>(null);
	let formLoading = $state(false);

	let statsComponent: QSOStats;

	onMount(async () => {
		await fetchQSOs();
	});

	async function fetchQSOs() {
		try {
			loading = true;
			const response = await fetch('/api/qsos?limit=100');
			if (!response.ok) {
				throw new Error('Failed to fetch QSOs');
			}
			const result = await response.json();
			qsos = result.data || [];
			error = null;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			loading = false;
		}
	}

	async function handleSubmit(qso: Partial<QSO>) {
		try {
			formLoading = true;

			const url = editingQSO ? `/api/qsos/${editingQSO.id}` : '/api/qsos';
			const method = editingQSO ? 'PUT' : 'POST';

			const response = await fetch(url, {
				method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(qso)
			});

			if (!response.ok) {
				const result = await response.json();
				throw new Error(result.error?.message || 'Failed to save QSO');
			}

			// Success! Refresh the list
			await fetchQSOs();
			if (statsComponent) {
				statsComponent.refresh();
			}

			// Close form
			showForm = false;
			editingQSO = null;
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Failed to save QSO');
		} finally {
			formLoading = false;
		}
	}

	async function handleDelete(id: number) {
		if (!confirm('Are you sure you want to delete this QSO?')) {
			return;
		}

		try {
			const response = await fetch(`/api/qsos/${id}`, {
				method: 'DELETE'
			});

			if (!response.ok) {
				throw new Error('Failed to delete QSO');
			}

			// Refresh list
			await fetchQSOs();
			if (statsComponent) {
				statsComponent.refresh();
			}
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Failed to delete QSO');
		}
	}

	function handleEdit(qso: QSO) {
		editingQSO = qso;
		showForm = true;
	}

	function handleCancel() {
		showForm = false;
		editingQSO = null;
	}

	function handleNewQSO() {
		editingQSO = null;
		showForm = true;
	}

	async function handleImportComplete() {
		await fetchQSOs();
		if (statsComponent) {
			statsComponent.refresh();
		}
	}
</script>

<svelte:head>
	<title>Logbook - Ham-Radio Cloud</title>
</svelte:head>

<Navigation />

<div class="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 p-4">
	<div class="mx-auto max-w-7xl">
		<!-- Header -->
		<div class="mb-6">
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-3xl font-bold text-gray-900">Logbook</h1>
					<p class="mt-1 text-sm text-gray-600">Manage your QSO contacts</p>
				</div>
			</div>
		</div>

		<!-- Stats and Tools Row -->
		<div class="mb-6 grid grid-cols-1 gap-6 lg:grid-cols-3">
			<!-- Stats -->
			<div class="lg:col-span-1">
				<QSOStats bind:this={statsComponent} />
			</div>

			<!-- ADIF Manager -->
			<div class="lg:col-span-2">
				<ADIFManager onImportComplete={handleImportComplete} />
			</div>
		</div>

		<!-- QSO List -->
		<div class="card">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-xl font-bold text-gray-900">Recent Contacts</h2>
				<button onclick={handleNewQSO} class="btn btn-primary">+ New QSO</button>
			</div>

			{#if error}
				<div class="rounded-lg bg-red-50 p-4 text-red-600">
					<strong>Error:</strong> {error}
				</div>
			{:else}
				<QSOTable {qsos} {loading} onEdit={handleEdit} onDelete={handleDelete} />
			{/if}
		</div>

		<!-- QSO Form Modal -->
		{#if showForm}
			<div
				class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4"
				onclick={(e) => e.target === e.currentTarget && handleCancel()}
			>
				<div class="max-h-[90vh] w-full max-w-4xl overflow-y-auto rounded-lg bg-white p-6 shadow-xl">
					<h2 class="mb-4 text-2xl font-bold text-gray-900">
						{editingQSO ? 'Edit QSO' : 'New QSO'}
					</h2>
					<QSOForm qso={editingQSO} onSubmit={handleSubmit} onCancel={handleCancel} loading={formLoading} />
				</div>
			</div>
		{/if}
	</div>
</div>
