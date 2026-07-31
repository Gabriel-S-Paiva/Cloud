<script lang="ts">
	import { sharedContents } from '$lib/stores/sharedContents.svelte';
	import Button from '$lib/components/UI/Button/Button.svelte';
	import Modal from '$lib/components/UI/Modal/Modal.svelte';
	import { Folder as FolderIcon, FileText } from '@lucide/svelte';
	import type { SharedFolder, SharedFile } from '$lib/types';

	let {
		item,
		direction
	}: { item: SharedFolder | SharedFile; direction: 'incoming' | 'outgoing' } = $props();

	const isFile = (i: SharedFolder | SharedFile): i is SharedFile => 'size' in i;
	const otherUser = direction === 'incoming' ? item.ownedByUsername : item.sharedWith;

	let confirmOpen = $state(false);
</script>

<div class="flex items-center gap-3 px-3 py-2.5 border border-border rounded-lg bg-surface-raised">
	{#if isFile(item)}
		<FileText size={18} class="text-accent-secondary shrink-0" />
	{:else}
		<FolderIcon size={18} class="text-accent shrink-0" />
	{/if}

	<span class="flex-1 truncate text-sm text-text">{item.displayName}</span>

	<span class="text-xs text-text-muted shrink-0 hidden sm:inline">
		{direction === 'incoming' ? 'from' : 'with'} {otherUser}
	</span>

	{#if direction === 'outgoing'}
		<select
			value={item.permissions}
			onchange={(e) =>
				sharedContents.updatePermission(item.shareId, e.currentTarget.value as 'Edit' | 'View')}
			class="h-8 rounded-md border border-border bg-surface px-2 text-xs text-text shrink-0"
		>
			<option value="View">View</option>
			<option value="Edit">Edit</option>
		</select>
	{:else}
		<span
			class="text-xs px-2 py-1 rounded-full border shrink-0"
			style={item.permissions === 'Edit'
				? 'color: var(--accent-secondary); border-color: var(--accent-secondary);'
				: 'color: var(--text-muted); border-color: var(--border);'}
		>
			{item.permissions}
		</span>
	{/if}

	<Button variant="secondary" onclick={() => (confirmOpen = true)}>Unshare</Button>
</div>

<Modal open={confirmOpen} onclose={() => (confirmOpen = false)}>
	<div class="flex flex-col gap-3">
		<h3 class="font-display text-lg text-text">Remove share?</h3>
		<p class="text-sm text-text">
			{#if direction === 'incoming'}
				You'll lose access to <strong>{item.displayName}</strong>.
			{:else}
				<strong>{item.sharedWith}</strong> will lose access to <strong>{item.displayName}</strong>.
			{/if}
		</p>
		<div class="flex gap-2">
			<Button variant="secondary" onclick={() => (confirmOpen = false)}>Cancel</Button>
			<Button
				variant="danger"
				onclick={() => {
					sharedContents.unshare(item.shareId, direction);
					confirmOpen = false;
				}}
			>
				Confirm
			</Button>
		</div>
	</div>
</Modal>