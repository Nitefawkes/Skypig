<script lang="ts">
	import type { QSO } from '$types';
	import { BANDS, MODES } from '$types';

	interface Props {
		qso?: QSO | null;
		onSubmit: (qso: Partial<QSO>) => void;
		onCancel?: () => void;
		loading?: boolean;
	}

	let { qso = null, onSubmit, onCancel, loading = false }: Props = $props();

	let formData = $state({
		callsign: qso?.callsign || '',
		time_on: qso?.time_on || new Date().toISOString().slice(0, 16),
		time_off: qso?.time_off || '',
		band: qso?.band || '20m',
		mode: qso?.mode || 'SSB',
		freq: qso?.freq || 0,
		rst_sent: qso?.rst_sent || '59',
		rst_rcvd: qso?.rst_rcvd || '59',
		name: qso?.name || '',
		qth: qso?.qth || '',
		gridsquare: qso?.gridsquare || '',
		country: qso?.country || '',
		state: qso?.state || '',
		tx_pwr: qso?.tx_pwr || 100,
		comment: qso?.comment || ''
	});

	function handleSubmit(e: Event) {
		e.preventDefault();

		const submitData: Partial<QSO> = {
			...formData,
			callsign: formData.callsign.toUpperCase(),
			gridsquare: formData.gridsquare.toUpperCase(),
			state: formData.state.toUpperCase(),
			freq: formData.freq || undefined,
			tx_pwr: formData.tx_pwr || undefined
		};

		// Remove empty strings
		Object.keys(submitData).forEach((key) => {
			if (submitData[key as keyof typeof submitData] === '') {
				delete submitData[key as keyof typeof submitData];
			}
		});

		onSubmit(submitData);
	}

	function setCurrentTime() {
		formData.time_on = new Date().toISOString().slice(0, 16);
	}
</script>

<form onsubmit={handleSubmit} class="space-y-6">
	<!-- Essential Fields -->
	<div class="rounded-lg border border-gray-200 p-4">
		<h3 class="mb-4 font-semibold text-gray-900">Essential Information</h3>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="callsign" class="label">Callsign *</label>
				<input
					id="callsign"
					type="text"
					bind:value={formData.callsign}
					required
					class="input"
					placeholder="W1AW"
				/>
			</div>

			<div>
				<label for="time_on" class="label">
					Time On (UTC) *
					<button
						type="button"
						onclick={setCurrentTime}
						class="ml-2 text-xs text-primary-600 hover:text-primary-800"
					>
						Set Now
					</button>
				</label>
				<input
					id="time_on"
					type="datetime-local"
					bind:value={formData.time_on}
					required
					class="input"
				/>
			</div>

			<div>
				<label for="band" class="label">Band</label>
				<select id="band" bind:value={formData.band} class="input">
					<option value="">Select band</option>
					{#each BANDS as band}
						<option value={band}>{band}</option>
					{/each}
				</select>
			</div>

			<div>
				<label for="mode" class="label">Mode</label>
				<select id="mode" bind:value={formData.mode} class="input">
					<option value="">Select mode</option>
					{#each MODES as mode}
						<option value={mode}>{mode}</option>
					{/each}
				</select>
			</div>

			<div>
				<label for="freq" class="label">Frequency (MHz)</label>
				<input
					id="freq"
					type="number"
					step="0.001"
					bind:value={formData.freq}
					class="input"
					placeholder="14.250"
				/>
			</div>

			<div>
				<label for="tx_pwr" class="label">TX Power (W)</label>
				<input
					id="tx_pwr"
					type="number"
					bind:value={formData.tx_pwr}
					class="input"
					placeholder="100"
				/>
			</div>
		</div>
	</div>

	<!-- Signal Reports -->
	<div class="rounded-lg border border-gray-200 p-4">
		<h3 class="mb-4 font-semibold text-gray-900">Signal Reports</h3>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="rst_sent" class="label">RST Sent</label>
				<input
					id="rst_sent"
					type="text"
					bind:value={formData.rst_sent}
					class="input"
					placeholder="59"
				/>
			</div>

			<div>
				<label for="rst_rcvd" class="label">RST Received</label>
				<input
					id="rst_rcvd"
					type="text"
					bind:value={formData.rst_rcvd}
					class="input"
					placeholder="59"
				/>
			</div>
		</div>
	</div>

	<!-- Contact Information -->
	<div class="rounded-lg border border-gray-200 p-4">
		<h3 class="mb-4 font-semibold text-gray-900">Contact Information</h3>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<div>
				<label for="name" class="label">Name</label>
				<input id="name" type="text" bind:value={formData.name} class="input" placeholder="John" />
			</div>

			<div>
				<label for="qth" class="label">QTH</label>
				<input
					id="qth"
					type="text"
					bind:value={formData.qth}
					class="input"
					placeholder="Boston, MA"
				/>
			</div>

			<div>
				<label for="gridsquare" class="label">Grid Square</label>
				<input
					id="gridsquare"
					type="text"
					bind:value={formData.gridsquare}
					class="input"
					placeholder="FN42"
					maxlength="8"
				/>
			</div>

			<div>
				<label for="country" class="label">Country</label>
				<input
					id="country"
					type="text"
					bind:value={formData.country}
					class="input"
					placeholder="USA"
				/>
			</div>

			<div>
				<label for="state" class="label">State/Province</label>
				<input
					id="state"
					type="text"
					bind:value={formData.state}
					class="input"
					placeholder="MA"
					maxlength="2"
				/>
			</div>
		</div>
	</div>

	<!-- Comments -->
	<div class="rounded-lg border border-gray-200 p-4">
		<h3 class="mb-4 font-semibold text-gray-900">Notes</h3>

		<div>
			<label for="comment" class="label">Comment</label>
			<textarea
				id="comment"
				bind:value={formData.comment}
				rows="3"
				class="input"
				placeholder="Great contact!"
			></textarea>
		</div>
	</div>

	<!-- Form Actions -->
	<div class="flex justify-end gap-4">
		{#if onCancel}
			<button type="button" onclick={onCancel} class="btn btn-secondary" disabled={loading}>
				Cancel
			</button>
		{/if}
		<button type="submit" class="btn btn-primary" disabled={loading}>
			{loading ? 'Saving...' : qso ? 'Update QSO' : 'Add QSO'}
		</button>
	</div>
</form>
