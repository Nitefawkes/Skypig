<script lang="ts">
	import type { QSO } from '$types';

	interface Props {
		qsos: QSO[];
		loading?: boolean;
		onEdit?: (qso: QSO) => void;
		onDelete?: (id: number) => void;
	}

	let { qsos, loading = false, onEdit, onDelete }: Props = $props();

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString();
	}

	function formatTime(dateStr: string): string {
		return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}
</script>

<div class="overflow-x-auto">
	{#if loading}
		<div class="flex items-center justify-center py-8">
			<div class="text-gray-500">Loading QSOs...</div>
		</div>
	{:else if qsos.length === 0}
		<div class="rounded-lg bg-gray-50 p-8 text-center">
			<div class="text-gray-500">No QSOs found</div>
			<div class="mt-2 text-sm text-gray-400">Add your first contact to get started!</div>
		</div>
	{:else}
		<table class="min-w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Date/Time
					</th>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Callsign
					</th>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Band
					</th>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Mode
					</th>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						RST
					</th>
					<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Name/QTH
					</th>
					<th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
						Actions
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200 bg-white">
				{#each qsos as qso}
					<tr class="hover:bg-gray-50">
						<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
							<div>{formatDate(qso.time_on)}</div>
							<div class="text-xs text-gray-500">{formatTime(qso.time_on)} UTC</div>
						</td>
						<td class="whitespace-nowrap px-4 py-3">
							<div class="text-sm font-bold text-primary-600">{qso.callsign}</div>
							{#if qso.gridsquare}
								<div class="text-xs text-gray-500">{qso.gridsquare}</div>
							{/if}
						</td>
						<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
							{qso.band || '-'}
						</td>
						<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
							{qso.mode || '-'}
						</td>
						<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500">
							<div class="text-xs">S: {qso.rst_sent || '-'}</div>
							<div class="text-xs">R: {qso.rst_rcvd || '-'}</div>
						</td>
						<td class="px-4 py-3 text-sm text-gray-900">
							{#if qso.name}
								<div class="font-medium">{qso.name}</div>
							{/if}
							{#if qso.qth}
								<div class="text-xs text-gray-500">{qso.qth}</div>
							{/if}
						</td>
						<td class="whitespace-nowrap px-4 py-3 text-right text-sm">
							<div class="flex justify-end gap-2">
								{#if onEdit}
									<button
										onclick={() => onEdit(qso)}
										class="text-primary-600 hover:text-primary-900"
									>
										Edit
									</button>
								{/if}
								{#if onDelete}
									<button
										onclick={() => onDelete(qso.id)}
										class="text-red-600 hover:text-red-900"
									>
										Delete
									</button>
								{/if}
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>

		<div class="mt-4 text-sm text-gray-500">
			Showing {qsos.length} QSO{qsos.length !== 1 ? 's' : ''}
		</div>
	{/if}
</div>
