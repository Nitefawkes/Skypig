<script lang="ts">
	interface ImportResult {
		total_records: number;
		imported_records: number;
		failed_records: number;
		skipped_records: number;
		errors?: string[];
	}

	interface Props {
		onImportComplete?: () => void;
	}

	let { onImportComplete }: Props = $props();

	let importLoading = $state(false);
	let importResult = $state<ImportResult | null>(null);
	let importError = $state<string | null>(null);
	let exportLoading = $state(false);

	async function handleImport(e: Event) {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		try {
			importLoading = true;
			importError = null;
			importResult = null;

			const text = await file.text();

			const response = await fetch('/api/qsos/import', {
				method: 'POST',
				headers: {
					'Content-Type': 'text/plain'
				},
				body: text
			});

			const result = await response.json();

			if (response.ok || response.status === 206) {
				importResult = result.data;
				if (onImportComplete) {
					onImportComplete();
				}
			} else {
				importError = result.error?.message || 'Import failed';
			}
		} catch (e) {
			importError = e instanceof Error ? e.message : 'Unknown error';
		} finally {
			importLoading = false;
			// Reset file input
			target.value = '';
		}
	}

	async function handleExport() {
		try {
			exportLoading = true;

			const response = await fetch('/api/qsos/export');
			if (!response.ok) {
				throw new Error('Export failed');
			}

			const blob = await response.blob();
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `hamradio_cloud_export_${new Date().toISOString().slice(0, 10)}.adi`;
			document.body.appendChild(a);
			a.click();
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (e) {
			alert(e instanceof Error ? e.message : 'Export failed');
		} finally {
			exportLoading = false;
		}
	}
</script>

<div class="card">
	<h2 class="mb-4 text-xl font-bold text-gray-900">ADIF Import/Export</h2>

	<div class="space-y-6">
		<!-- Export Section -->
		<div class="rounded-lg border border-gray-200 p-4">
			<h3 class="mb-2 font-semibold text-gray-900">Export Logbook</h3>
			<p class="mb-4 text-sm text-gray-600">
				Download your entire logbook in ADIF format for backup or use with other logging software.
			</p>
			<button
				onclick={handleExport}
				disabled={exportLoading}
				class="btn btn-primary flex items-center gap-2"
			>
				{#if exportLoading}
					<span>Exporting...</span>
				{:else}
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
						/>
					</svg>
					<span>Export to ADIF</span>
				{/if}
			</button>
		</div>

		<!-- Import Section -->
		<div class="rounded-lg border border-gray-200 p-4">
			<h3 class="mb-2 font-semibold text-gray-900">Import from ADIF</h3>
			<p class="mb-4 text-sm text-gray-600">
				Import QSOs from an ADIF file exported from another logging program.
			</p>

			<label class="btn btn-secondary inline-flex cursor-pointer items-center gap-2">
				{#if importLoading}
					<span>Importing...</span>
				{:else}
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
						/>
					</svg>
					<span>Choose ADIF File</span>
				{/if}
				<input
					type="file"
					accept=".adi,.adif"
					onchange={handleImport}
					disabled={importLoading}
					class="hidden"
				/>
			</label>

			<!-- Import Result -->
			{#if importResult}
				<div class="mt-4 rounded-lg bg-blue-50 p-4">
					<div class="mb-2 font-semibold text-blue-900">Import Complete</div>
					<div class="space-y-1 text-sm text-blue-800">
						<div>✓ Total Records: {importResult.total_records}</div>
						<div class="text-green-700">✓ Imported: {importResult.imported_records}</div>
						{#if importResult.failed_records > 0}
							<div class="text-red-700">✗ Failed: {importResult.failed_records}</div>
						{/if}
						{#if importResult.skipped_records > 0}
							<div class="text-yellow-700">⚠ Skipped: {importResult.skipped_records}</div>
						{/if}
					</div>

					{#if importResult.errors && importResult.errors.length > 0}
						<details class="mt-3">
							<summary class="cursor-pointer text-sm font-medium text-blue-900">
								View Errors ({importResult.errors.length})
							</summary>
							<div class="mt-2 max-h-40 overflow-y-auto rounded bg-white p-2 text-xs">
								{#each importResult.errors as error}
									<div class="text-red-600">{error}</div>
								{/each}
							</div>
						</details>
					{/if}
				</div>
			{/if}

			<!-- Import Error -->
			{#if importError}
				<div class="mt-4 rounded-lg bg-red-50 p-4 text-red-600">
					<strong>Import Failed:</strong> {importError}
				</div>
			{/if}
		</div>

		<!-- Help Text -->
		<div class="text-xs text-gray-500">
			<strong>Supported Format:</strong> ADIF 3.1.4 (.adi, .adif files)
		</div>
	</div>
</div>
