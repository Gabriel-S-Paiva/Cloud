<script lang="ts">
	import { Folder, FileText, Image, FileSpreadsheet } from '@lucide/svelte';

	const entries = [
		{ type: 'folder', name: 'Family Photos', meta: '48 items', icon: Folder },
		{ type: 'folder', name: 'Home Videos', meta: '12 items', icon: Folder },
		{ type: 'file', name: 'Tax Documents 2026.pdf', meta: '2.4 MB', icon: FileText },
		{ type: 'file', name: 'Reunion Panorama.jpg', meta: '18.1 MB', icon: Image },
		{ type: 'file', name: 'Household Budget.xlsx', meta: '340 KB', icon: FileSpreadsheet }
	];
</script>

<div class="overflow-hidden rounded-lg border border-border bg-surface-raised shadow-sm">
	<div class="flex items-center justify-between border-b border-border px-4 py-2.5">
		<span class="font-mono text-xs tracking-wide text-text-muted">/home</span>
		<span class="font-mono text-xs text-text-muted">5 items</span>
	</div>
	<ul>
		{#each entries as entry, i (entry.name)}
			<li
				class="flex items-center gap-3 px-4 py-2.5 {i !== entries.length - 1
					? 'border-b border-border'
					: ''}"
			>
				<entry.icon
					size={16}
					class={entry.type === 'folder'
						? 'shrink-0 text-accent'
						: 'shrink-0 text-accent-secondary'}
				/>
				<span class="flex-1 truncate text-sm text-text">{entry.name}</span>
				<span class="shrink-0 font-mono text-xs text-text-muted">{entry.meta}</span>
			</li>
		{/each}
	</ul>
</div>
