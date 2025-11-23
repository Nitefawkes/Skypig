<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { QSO } from '$types';

	let { qso = $bindable(null), mode = 'create' }: { qso?: QSO | null; mode?: 'create' | 'edit' } = $props();

	const dispatch = createEventDispatcher();

	let formData = $state({
		callsign: qso?.callsign || '',
		frequency: qso?.frequency || 14.074,
		band: qso?.band || '20m',
		mode: qso?.mode || 'FT8',
		rst_sent: qso?.rst_sent || '-10',
		rst_received: qso?.rst_received || '-10',
		qso_date: qso?.qso_date ? qso.qso_date.split('T')[0] : new Date().toISOString().split('T')[0],
		time_on: qso?.time_on ? qso.time_on.slice(0, 16) : new Date().toISOString().slice(0, 16),
		time_off: qso?.time_off ? qso.time_off.slice(0, 16) : '',
		grid_square: qso?.grid_square || '',
		country: qso?.country || '',
		state: qso?.state || '',
		county: qso?.county || '',
		comment: qso?.comment || '',
		tx_power: qso?.tx_power || 100
	});

	const bands = ['160m', '80m', '40m', '30m', '20m', '17m', '15m', '12m', '10m', '6m', '2m', '70cm'];
	const modes = ['SSB', 'CW', 'FT8', 'FT4', 'RTTY', 'PSK31', 'PSK63', 'JT65', 'JT9', 'AM', 'FM'];

	function handleSubmit() {
		// Convert form data to API format
		const payload = {
			...formData,
			time_on: new Date(formData.time_on).toISOString(),
			time_off: formData.time_off ? new Date(formData.time_off).toISOString() : undefined,
			callsign: formData.callsign.toUpperCase()
		};

		dispatch('submit', payload);
	}

	function handleCancel() {
		dispatch('cancel');
	}
</script>

<form onsubmit|preventDefault={handleSubmit} class="space-y-6">
	<!-- Callsign (Required) -->
	<div>
		<label for="callsign" class="mb-2 block text-sm font-medium text-slate-300">
			Callsign <span class="text-red-500">*</span>
		</label>
		<input
			id="callsign"
			type="text"
			bind:value={formData.callsign}
			placeholder="W1AW"
			required
			class="input uppercase"
			autocomplete="off"
		/>
	</div>

	<!-- Frequency and Band -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<label for="frequency" class="mb-2 block text-sm font-medium text-slate-300">
				Frequency (MHz)
			</label>
			<input
				id="frequency"
				type="number"
				step="0.001"
				bind:value={formData.frequency}
				class="input"
			/>
		</div>
		<div>
			<label for="band" class="mb-2 block text-sm font-medium text-slate-300">
				Band <span class="text-red-500">*</span>
			</label>
			<select id="band" bind:value={formData.band} required class="input">
				{#each bands as band}
					<option value={band}>{band}</option>
				{/each}
			</select>
		</div>
	</div>

	<!-- Mode and Power -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<label for="mode" class="mb-2 block text-sm font-medium text-slate-300">
				Mode <span class="text-red-500">*</span>
			</label>
			<select id="mode" bind:value={formData.mode} required class="input">
				{#each modes as mode}
					<option value={mode}>{mode}</option>
				{/each}
			</select>
		</div>
		<div>
			<label for="tx_power" class="mb-2 block text-sm font-medium text-slate-300">
				TX Power (W)
			</label>
			<input id="tx_power" type="number" bind:value={formData.tx_power} class="input" />
		</div>
	</div>

	<!-- RST Sent and Received -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<label for="rst_sent" class="mb-2 block text-sm font-medium text-slate-300">
				RST Sent
			</label>
			<input id="rst_sent" type="text" bind:value={formData.rst_sent} class="input" />
		</div>
		<div>
			<label for="rst_received" class="mb-2 block text-sm font-medium text-slate-300">
				RST Received
			</label>
			<input id="rst_received" type="text" bind:value={formData.rst_received} class="input" />
		</div>
	</div>

	<!-- Date and Time -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
		<div>
			<label for="qso_date" class="mb-2 block text-sm font-medium text-slate-300">
				QSO Date
			</label>
			<input id="qso_date" type="date" bind:value={formData.qso_date} class="input" />
		</div>
		<div>
			<label for="time_on" class="mb-2 block text-sm font-medium text-slate-300">
				Time On (UTC)
			</label>
			<input id="time_on" type="datetime-local" bind:value={formData.time_on} class="input" />
		</div>
		<div>
			<label for="time_off" class="mb-2 block text-sm font-medium text-slate-300">
				Time Off (UTC)
			</label>
			<input id="time_off" type="datetime-local" bind:value={formData.time_off} class="input" />
		</div>
	</div>

	<!-- Location -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<label for="grid_square" class="mb-2 block text-sm font-medium text-slate-300">
				Grid Square
			</label>
			<input
				id="grid_square"
				type="text"
				bind:value={formData.grid_square}
				placeholder="FN42"
				class="input uppercase"
				maxlength="6"
			/>
		</div>
		<div>
			<label for="country" class="mb-2 block text-sm font-medium text-slate-300">
				Country
			</label>
			<input id="country" type="text" bind:value={formData.country} class="input" />
		</div>
	</div>

	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div>
			<label for="state" class="mb-2 block text-sm font-medium text-slate-300">
				State/Province
			</label>
			<input id="state" type="text" bind:value={formData.state} class="input" />
		</div>
		<div>
			<label for="county" class="mb-2 block text-sm font-medium text-slate-300">
				County
			</label>
			<input id="county" type="text" bind:value={formData.county} class="input" />
		</div>
	</div>

	<!-- Comment -->
	<div>
		<label for="comment" class="mb-2 block text-sm font-medium text-slate-300">
			Comment
		</label>
		<textarea
			id="comment"
			bind:value={formData.comment}
			rows="3"
			class="input"
			placeholder="Additional notes..."
		></textarea>
	</div>

	<!-- Buttons -->
	<div class="flex justify-end gap-4">
		<button type="button" onclick={handleCancel} class="btn-secondary">
			Cancel
		</button>
		<button type="submit" class="btn-primary">
			{mode === 'create' ? 'Log QSO' : 'Update QSO'}
		</button>
	</div>
</form>
